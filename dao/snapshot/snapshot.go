package snapshot

import (
	"context"
	"github.com/gzltommy/go-graphql-client"
	"golang.org/x/oauth2"
)

const (
	SnapshotUrl   = "https://hub.snapshot.org/graphql"
	SnapshotVPUrl = "https://score.snapshot.org/api/scores"
)

func GetClient() *graphql.Client {
	client := graphql.NewClient(SnapshotUrl, nil)
	return client
}

func GetAuthClient(accessToken string) *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return graphql.NewClient(SnapshotUrl, httpClient)
}
