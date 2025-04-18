package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"redditDataCompiler/controllers"
)

type Stats struct {
	controller controllers.Stats
}

func NewStats(controller controllers.Stats) *Stats {
	return &Stats{
		controller: controller,
	}
}

func (h Stats) GetUserWithMostPosts(c echo.Context) error {
	user, err := h.controller.GetUserWithMostPosts(c.Request().Context())
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (h Stats) GetPostWithMostUpvotes(c echo.Context) error {
	post, err := h.controller.GetPostWithMostUpvotes(c.Request().Context())
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.JSON(http.StatusOK, post)
}
