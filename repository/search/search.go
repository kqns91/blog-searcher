package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/kqns91/blog-searcher/model"
	"github.com/kqns91/blog-searcher/model/response"
	opensearchv2 "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type OpenSearch interface {
	Search(ctx context.Context, query string) (*response.SearchResponse, error)
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

func (o *osearch) Search(ctx context.Context, query string) (*response.SearchResponse, error) {
	q := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"multi_match": map[string]any{
							"query": query,
							"fields": []string{
								"title",
								"name",
								"text",
							},
						},
					},
				},
			},
		},
		"_source": []string{
			"title",
			"name",
			"date",
			"img",
			"link",
		},
		"highlight": map[string]any{
			"fields": map[string]any{
				"text": map[string]any{
					"number_of_fragments": 1,
					"fragment_size":       200,
				},
			},
		},
	}

	body, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	req := opensearchapi.SearchRequest{
		Index: []string{"blogs"},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, o.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}

	defer res.Body.Close()

	var v response.SearchResponse

	err = json.NewDecoder(res.Body).Decode(&v)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return &v, nil
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
