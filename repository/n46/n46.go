package n46

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kqns91/blog-searcher/model"
)

type N46 interface {
	GetBlogs(ctx context.Context) ([]*model.Blog, error)
}

type n46 struct {
	baseURL       string
	defaultClient *http.Client
}

func New(baseURL string) N46 {
	return &n46{
		baseURL:       baseURL,
		defaultClient: http.DefaultClient,
	}
}

const defaultLimit = 100

func (n *n46) GetBlogs(ctx context.Context) ([]*model.Blog, error) {
	if n.baseURL == "" {
		return nil, errors.New("baseURL is empty")
	}

	result := []*model.Blog{}
	st := 0
	rw := defaultLimit

	query := url.Values{}

	for {
		query.Set("st", strconv.Itoa(st))
		query.Set("rw", strconv.Itoa(rw))

		urlStr := fmt.Sprintf("%s/blog?%s", n.baseURL, query.Encode())
		res := &model.Blogs{}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := n.defaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		defer resp.Body.Close()

		bytes, err := removeRes(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to remove res: %w", err)
		}

		err = json.Unmarshal(bytes, &res)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}

		result = append(result, res.Data...)

		total, err := strconv.Atoi(res.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to convert int: %w", err)
		}

		st += rw

		if total <= len(result) {
			break
		}
	}

	return result, nil
}

func removeRes(body io.Reader) ([]byte, error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read all: %w", err)
	}

	rawStr := string(bytes)
	convertStr := strings.TrimRight(strings.TrimLeft(rawStr, "res("), ");")

	return []byte(convertStr), nil
}
