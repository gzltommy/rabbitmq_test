package model

type ProposalVote struct {
	Id         int64  `gorm:"primary_key" json:"id"`
	ProposalID string `json:"-"` // 提案的 ID

	UserID  int64   `json:"-"`        // 用户 ID (我们平台的用户才有)
	VoterID string  `json:"voter_id"` // 用户投票 ID (地址)
	Choice  int     `json:"choice"`   // 第几个选项：1 或 2 或 3 或 ...
	VP      float64 `json:"vp"`       // vp 值
	Created int64   `json:"created"`  // 投票的时间戳
}

type Strategy struct {
	Name    string      `json:"name"`
	Network string      `json:"network"`
	Params  interface{} `json:"params"` // 注意：这个任意参数的特殊解法
}

type SnapshotProposal struct {
	ID         string     `json:"id"`      // 提案的 ID，如：“0xa3af2279022db004bc5aae5df8adba5113fd5f6b12dc40c7c4c2c7ebfc4e9c7a”
	Ipfs       string     `json:"ipfs"`    // ipfs
	Author     string     `json:"author"`  // 创建人钱包地址，如："0x4bDA26282Cd8D7E5B5253e339d9E7906B039b2e6"
	Created    int        `json:"created"` // 创建时间
	Type       string     `json:"type"`    // 提案投票类型（"single-choice"）
	Network    string     `json:"network"`
	Strategies []Strategy `json:"strategies"`

	Title    string   `json:"title"`    // 主题
	Body     string   `json:"body"`     // 说明 & 选项（body）
	Choices  []string `json:"choices"`  // 说明 & 选项（选项）
	Start    int      `json:"start"`    // 开始时间
	End      int      `json:"end"`      // 结束时间
	Snapshot string   `json:"snapshot"` // blockNumber
	State    string   `json:"state"`    // 状态（"closed"、"active"）
	Link     string   `json:"link"`     // Snapshot 链接

	// 以下数据要 “closed” 状态的提案才能返回数据
	Scores       []float64 `json:"scores"`       // 选项的所得 VP 值
	Scores_Total float64   `json:"scores_total"` // 总的投票 VP 值（注意：这里不能修改名字）
	Votes        int       `json:"votes"`        // 总的投票人数

	//Scores_State graphql.String
	//Scores_By_Strategy interface{}
	//Scores_Updated     graphql.Int
}

type Proposal struct {
	Id             int64  `gorm:"primary_key" json:"id"`
	ProposalID     string `json:"-"`          // 提案的 ID，如：“0xa3af2279022db004bc5aae5df8adba5113fd5f6b12dc40c7c4c2c7ebfc4e9c7a”
	ProposalAuthor string `json:"-"`          // 创建人钱包地址，如："0x4bDA26282Cd8D7E5B5253e339d9E7906B039b2e6"
	ProposalData   string `json:"-"`          // 提案原始数据“SnapshotProposal”
	CreateTime     int64  `json:"-"`          // 创建时间
	StateCode      uint8  `json:"state_code"` // 提案状态（"pending","closed"、"active"）我们自己扩展 3 个状态：reviewing + published + deleted

	// 其他数据
	PublishUserName string `json:"publish_user_name"` // 地址（提案颁布者名字）

	// 提案统计数据
	StatisticIsEnd int8  `json:"-"` // 是否以结束了统计
	StatisticTime  int64 `json:"-"` // 上次统计时间

	Proposal SnapshotProposal `gorm:"-"` // 提案详情
}
