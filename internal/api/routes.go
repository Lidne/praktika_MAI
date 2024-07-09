package api

import (
	_ "github.com/Lidne/praktika_MAI/docs"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Services struct {
}

func NewRouter(e *echo.Group) {
	statisticsGroup := e.Group("/statistics")
	statisticsGroup.GET("/sales", getSales)
	statisticsGroup.GET("/users", getUsers)
	statisticsGroup.GET("/products", getProducts)
}

// getUsers godoc
// @Summary Get users
// @Description Retrieve a list of users
// @ID get-users
// @Produce json
// @Success 200 {object} map[string]interface{} "List of users"
// @Router /api/statistics/users [get]
func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get users",
	})
}

func getSales(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get sales",
	})
}

func getProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get products",
	})
}
