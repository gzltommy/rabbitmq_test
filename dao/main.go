package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gzltommy/go-graphql-client"
	"io/ioutil"
	"net/http"
	"rbtmq/dao/model"
	"rbtmq/dao/snapshot"
	"strconv"
)

var space = "cocoslabs.eth"

//var space = "zorromiaotommy.eth"

func main() {
	p := queryProposals("0x8344ae34789a923a5dd9c038521e1be465aff7a49ea6eaab0c3a8dbc7439d30b")
	list := statisticProposalVotes(p)

	address := make([]string, 0, len(list))
	for _, v := range list {
		address = append(address, v.Voter)
	}

	buf, _ := json.Marshal(address)

	fmt.Println(len(address), "\n", string(buf))
}

func queryProposals(id string) *model.SnapshotProposal {
	client := snapshot.GetClient()

	var query struct {
		Proposal *model.SnapshotProposal `graphql:"proposal(id:$id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(id), // limit
	}
	jsonByte, err := client.QueryRawByGet(context.Background(), &query, variables)
	if err != nil || jsonByte == nil || len(*jsonByte) == 0 {
		fmt.Errorf("queryProposals error:", err, jsonByte == nil)
		return nil
	}

	err = json.Unmarshal(*jsonByte, &query)
	if err != nil {
		fmt.Errorf("Unmarshal error:", err)
		return nil
	}
	return query.Proposal
}

func statisticProposalVotes(prop *model.SnapshotProposal) []Vote {
	// 1、解析提案详情数据
	client := snapshot.GetClient()
	var page, limit = 1, 400

	list := make([]Vote, 0, 400)
	voteList, _ := handleProposalVotes(client, page, limit, space, prop)
	list = append(list, voteList...)
	for len(voteList) == limit {
		page++
		voteList, _ = handleProposalVotes(client, page, limit, space, prop)
		list = append(list, voteList...)
	}

	return list

}

// 投票数据
type Vote struct {
	Voter   string `json:"voter"`
	Choice  int    `json:"choice"` // 第几个选项：1 或 2 或 3 或 ...
	Created int64  `json:"created"`

	VP float64 `json:"vp"` // vp 值（没有结束的提案查不到该值）
}

// 处理投票（从 snapshot 上查询出投票信息 + 保存到本地）
func handleProposalVotes(client *graphql.Client, page, limit int, space string, p *model.SnapshotProposal) ([]Vote, *model.SnapshotProposal) {
	// 从 snapshot 上获取用户的投票信息 + snapshot 上的提案信息
	var query struct {
		Votes    []Vote                 `graphql:"votes(first:$first,skip:$skip,orderBy:\"created\",orderDirection:asc,where: {proposal: $proposal,created_gt:$created_gt})"`
		Proposal model.SnapshotProposal `graphql:"proposal(id:$id)"`
	}
	variables := map[string]interface{}{
		"first":      graphql.Int(limit),              // limit
		"skip":       graphql.Int(limit * (page - 1)), // offset=limit*(page-1)
		"proposal":   graphql.String(p.ID),
		"created_gt": graphql.Int(0),
		"id":         graphql.String(p.ID),
	}
	jsonByte, err := client.QueryRawByGet(context.Background(), &query, variables)
	if err != nil || jsonByte == nil || len(*jsonByte) == 0 {
		fmt.Errorf("queryProposals error:", err, jsonByte == nil)
		return nil, nil
	}
	err = json.Unmarshal(*jsonByte, &query)
	if err != nil {
		fmt.Errorf("Unmarshal error:", err)
		return nil, nil
	}

	// 没有投票数据（1：还没开始投票；2：已经开始投票但是没人投）
	if len(query.Votes) == 0 {
		return nil, &query.Proposal
	}

	//chainId, err := strconv.Atoi(p.Proposal.Network)
	//if err != nil {
	//	fmt.Errorf("strconv.Atoi() error:", p.Proposal.Network, err)
	//	return 0, nil
	//}

	// 如果该提案的状态不是“closed”的，需要通过另外的接口获取 VP 值
	queryVPVoters := make([]string, 0, len(query.Votes))
	for _, v := range query.Votes {
		// 1、不知为何，查出的数据 choice 有可能是 0（snapshot bug）
		// 2、已经取到 VP 值的不需要再次查询了
		if v.Choice == 0 || v.Choice > len(p.Choices) || v.VP > 0 {
			continue
		}
		queryVPVoters = append(queryVPVoters, v.Voter)
	}

	// 批量查询
	if len(queryVPVoters) > 0 {
		snapshotNum, _ := strconv.Atoi(p.Snapshot)
		vps := queryVoteVP(&reqVotesVP{
			Params: ReqVPParam{
				Space:      space,
				Strategies: p.Strategies,
				Network:    p.Network,
				Snapshot:   snapshotNum,
				Addresses:  queryVPVoters,
			},
		})
		if vps == nil {
			fmt.Errorf("vps is nil")
			return nil, nil
		}

		// 设置 vp 值
		for i, v := range query.Votes {
			// 1、不知为何，查出的数据 choice 有可能是 0（snapshot bug）
			// 2、已经取到 VP 值的不需要再次查询了
			if v.Choice == 0 || v.Choice > len(p.Choices) || v.VP > 0 {
				continue
			}

			// 该用户的 vp 值（每种投票策略都对应了一个 投票 vp ）
			vp := float64(0)
			for _, vv := range vps.Result.Scores {
				if value, ok := vv[v.Voter]; ok {
					vp += value
				}
			}
			query.Votes[i].VP = vp
		}
	}

	// 处理记录投票数据：（用户有可能改选投票）
	for _, v := range query.Votes {
		// 1、不知为何，查出的数据 choice 有可能是 0（snapshot bug）
		if v.Choice == 0 || v.Choice > len(p.Choices) {
			continue
		}

		// 处理用户改选投票逻辑
		puv := model.ProposalVote{
			ProposalID: p.ID,
			VoterID:    v.Voter, //地址
			Choice:     v.Choice,
			VP:         v.VP,
			Created:    v.Created,
		}

		// 更新统计数据
		p.Scores[puv.Choice-1] += puv.VP
		p.Scores_Total += puv.VP
		p.Votes += 1
	}
	return query.Votes, &query.Proposal
}

type ReqVPParam struct {
	Space      string           `json:"space"`
	Strategies []model.Strategy `json:"strategies"`
	Network    string           `json:"network"`
	Snapshot   int              `json:"snapshot"`
	Addresses  []string         `json:"addresses"`
}

type reqVotesVP struct {
	Params ReqVPParam `json:"params"`
}

type respVotesVP struct {
	Result struct {
		State  string               `json:"state"`
		Cache  bool                 `json:"cache"`
		Scores []map[string]float64 `json:"scores"`
	} `json:"result"`
}

// 从获取 VP 值
func queryVoteVP(req *reqVotesVP) *respVotesVP {
	bs, _ := json.Marshal(req)
	resp, err := http.Post(snapshot.SnapshotVPUrl, "application/json;charset=utf-8", bytes.NewBuffer(bs))
	if err != nil {
		fmt.Errorf("http.Pos err:", err)
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	res := respVotesVP{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Errorf("Unmarshal err:", err)
		return nil
	}
	return &res
}
