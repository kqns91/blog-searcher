package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kqns91/blog-searcher/handler"
	"github.com/kqns91/blog-searcher/repository"
	"github.com/kqns91/blog-searcher/repository/n46"
	"github.com/kqns91/blog-searcher/repository/search"
	"github.com/kqns91/blog-searcher/usecase"
	"github.com/opensearch-project/opensearch-go/v2"
)

var exitCode = 0

func main() {
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	opensearchClient, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{os.Getenv("OPEN_SEARCH_ADDRESS")},
		Username:  os.Getenv("USER_NAME"),
		Password:  os.Getenv("PASSWORD"),
	})
	if err != nil {
		log.Printf("failed to create opensearch client: %v", err.Error())

		exitCode = 1

		return
	}

	n46 := n46.New(os.Getenv("N46_BASEURL"))
	search := search.New(opensearchClient)
	repo := repository.New(n46, search)
	uc := usecase.New(repo)
	h := handler.New(uc)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine = (handler.SetRouteFunc(h))(engine)

	log.Printf("listening and serving on port %s", os.Getenv("PORT"))

	if engine.Run(":" + os.Getenv("PORT")); err != nil {
		log.Printf("failed to serve: %v", err)

		exitCode = 1

		return
	}
}

