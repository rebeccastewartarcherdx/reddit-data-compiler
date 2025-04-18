package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"redditDataCompiler/controllers"
	"redditDataCompiler/handlers"
	"redditDataCompiler/poller"
	"time"
)

func main() {
	// start fetching data from reddit
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	poll := poller.NewPoller(httpClient)
	pollChannel := poll.Poll()
	stats := controllers.NewStatsProcessor()

	go func() {
		for r := range pollChannel {
			if r.Err == nil {
				go stats.CalculateTopUserAndPost(r.Data)
				fmt.Println()
				fmt.Println()
				fmt.Println()
			}
		}
	}()

	// set up http server
	handler := handlers.NewStats(stats)
	e := echo.New()
	e.GET("/most_upvoted_post", handler.GetPostWithMostUpvotes)
	e.GET("/user_with_most_posts", handler.GetUserWithMostPosts)
	if err := e.Start("localhost:8005"); err != nil {
		e.Logger.Fatal(err)
	}
}
