package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kqns91/blog-searcher/usecase"
)

type Handler interface {
	IndexDocument(ctx context.Context, c *gin.Context) error
	Search(ctx context.Context, c *gin.Context) error
}

type httpHandler struct {
	uc usecase.Usecase
}

func New(uc usecase.Usecase) Handler {
	return &httpHandler{
		uc: uc,
	}
}

func SetRouteFunc(handler Handler) func(*gin.Engine) *gin.Engine {
	return func(engine *gin.Engine) *gin.Engine {
		routes := map[string]struct {
			fn     func(c *gin.Context)
			method string
		}{
			"/documents/index":  {fn: fn(handler.IndexDocument), method: http.MethodGet},
			"/documents/search": {fn: fn(handler.Search), method: http.MethodGet},
		}

		api := engine.Group("/api")
		for path, route := range routes {
			api.Handle(route.method, path, route.fn)
		}

		return engine
	}
}

func fn(f func(ctx context.Context, c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if err := f(ctx, c); err != nil {
			handleError(ctx, c, err)
		}
	}
}

func handleError(ctx context.Context, c *gin.Context, err error) {
	log.Printf("failure: %v", err.Error())

	c.JSON(
		http.StatusInternalServerError,
		map[string]any{
			"error_message": err.Error(),
		},
	)
}

func (h *httpHandler) IndexDocument(ctx context.Context, c *gin.Context) error {
	err := h.uc.IndexDocument(ctx)
	if err != nil {
		return fmt.Errorf("failed to IndexDocument: %w", err)
	}

	c.Status(http.StatusOK)

	return nil
}

func (h *httpHandler) Search(ctx context.Context, c *gin.Context) error {
	res, err := h.uc.Search(ctx, c.Query("query"))
	if err != nil {
		return fmt.Errorf("failed to IndexDocument: %w", err)
	}

	c.JSON(http.StatusOK, res)

	return nil
}
