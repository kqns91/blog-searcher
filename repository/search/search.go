package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/kqns91/blog-searcher/model"
	opensearchv2 "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type OpenSearch interface {
	Search(ctx context.Context) (*model.Blog, error)
	IndexDocument(ctx context.Context, doc *model.Blog) error
}

type osearch struct {
	client *opensearchv2.Client
}

func New(client *opensearchv2.Client) OpenSearch {
	return &osearch{
		client: client,
	}
}

func (o *osearch) Search(ctx context.Context) (*model.Blog, error) {
	return nil, nil
}

func (o *osearch) IndexDocument(ctx context.Context, doc *model.Blog) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index: "blogs",
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, o.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}

	defer res.Body.Close()

	return nil
}
