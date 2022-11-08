package repository

import (
	"github.com/kqns91/blog-searcher/repository/n46"
	"github.com/kqns91/blog-searcher/repository/search"
)

type Repository interface {
	N46() n46.N46
	OpenSearch() search.OpenSearch
}

type repository struct {
	n46     n46.N46
	osearch search.OpenSearch
}

func New(n46 n46.N46, osearch search.OpenSearch) Repository {
	return &repository{
		n46:     n46,
		osearch: osearch,
	}
}

func (r *repository) N46() n46.N46 {
	return r.n46
}

func (r *repository) OpenSearch() search.OpenSearch {
	return r.osearch
}
