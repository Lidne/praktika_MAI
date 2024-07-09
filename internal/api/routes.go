package api

import (
	_ "github.com/Lidne/praktika_MAI/docs"
	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Group) {
	statisticsGroup := e.Group("/statistics")
	statisticsGroup.GET("/sales", getSales)
	statisticsGroup.GET("/users", getUsers)
	statisticsGroup.GET("/products", getProducts)
}
