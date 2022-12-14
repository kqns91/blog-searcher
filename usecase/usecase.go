package usecase

import (
	"context"
	"fmt"

	"github.com/kqns91/blog-searcher/model/response"
	"github.com/kqns91/blog-searcher/repository"
)

type Usecase interface {
	Search(ctx context.Context, query string) (*response.SearchResponse, error)
	IndexDocument(ctx context.Context) error
}

type ucase struct {
	repo repository.Repository
}

func New(repo repository.Repository) Usecase {
	return &ucase{
		repo: repo,
	}
}

func (u *ucase) Search(ctx context.Context, query string) (*response.SearchResponse, error) {
	res, err := u.repo.OpenSearch().Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}

	return res, nil
}

func (u *ucase) IndexDocument(ctx context.Context) error {
	blogs, err := u.repo.N46().GetBlogs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get blogs: %w", err)
	}

	for _, b := range blogs {
		if err := u.repo.OpenSearch().IndexDocument(ctx, b); err != nil {
			return fmt.Errorf("failed to index document: %w", err)
		}
	}

	return nil
}
