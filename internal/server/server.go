package server

import (
	"context"
	"github.com/Lidne/praktika_MAI/config"
	_ "github.com/Lidne/praktika_MAI/docs"
	"github.com/Lidne/praktika_MAI/internal/product"
	productRepo "github.com/Lidne/praktika_MAI/internal/product/repository"
	"github.com/Lidne/praktika_MAI/internal/sell"
	sellRepo "github.com/Lidne/praktika_MAI/internal/sell/repository"
	"github.com/Lidne/praktika_MAI/internal/user"
	userRepo "github.com/Lidne/praktika_MAI/internal/user/repository"
	"github.com/Lidne/praktika_MAI/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	maxHeaderBytes  = 1 << 20
	gzipLevel       = 5
	stackSize       = 1 << 10 // 1 KB
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "2M"
	kafkaGroupID    = "products_group"
)

// server
type server struct {
	log      logger.Logger
	cfg      *config.Config
	tracer   opentracing.Tracer
	dbclient *pgxpool.Pool
	echo     *echo.Echo
}

type Services struct {
	user    user.UserRepository
	product product.ProductRepository
	sell    sell.SellRepository
	ctx     context.Context
}

func NewServices(pool *pgxpool.Pool, ctx context.Context) *Services {
	return &Services{
		user:    userRepo.NewUserRepo(pool),
		product: productRepo.NewProductRepo(pool),
		sell:    sellRepo.NewSellRepo(pool),
		ctx:     ctx,
	}
}

// NewServer constructor
func NewServer(log logger.Logger, cfg *config.Config, tracer opentracing.Tracer, db *pgxpool.Pool) *server {
	return &server{log: log, cfg: cfg, tracer: tracer, dbclient: db, echo: echo.New()}
}

// Run Start server
func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services := NewServices(s.dbclient, ctx)
	api := s.echo.Group("/api")
	api.GET("/users", services.getUsers)
	api.GET("/users/:id", services.getUserById)
	api.GET("/products", services.getProducts)
	api.GET("/products/:id", services.getProductById)
	api.GET("/sales", services.getSales)
	api.GET("/sales/:id", services.getSellById)
	api.GET("/sales/interval", services.getSalesDate)

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := s.echo.Start(s.cfg.Http.Port); err != nil {
		s.log.Error(err)
		cancel()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.log.Errorf("ctx.Done: %v", done)
	}

	if err := s.echo.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "echo.Server.Shutdown")
	}

	/*if err := metricsServer.Shutdown(ctx); err != nil {
		s.log.Errorf("metricsServer.Shutdown: %v", err)
	}*/
	//grpcServer.GracefulStop()
	s.log.Info("Server Exited Properly")

	return nil
}

// getUsers godoc
//
//	@Summary		Get Users
//	@Tags			Users
//	@Description	Get all users
//	@ID				get-users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"data"
//	@Failure		500	{object}	map[string]interface{}	"error"
//	@Router			/api/users [get]
func (s *Services) getUsers(c echo.Context) error {
	usrs, err := s.user.FindAll(s.ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get users",
			"err":     err,
		})
	}
	res := []echo.Map{}
	for _, usr := range usrs {
		res = append(res, echo.Map{
			"id":         usr.ID,
			"name":       usr.Name,
			"updated_at": usr.UpdatedAt,
			"login":      usr.Login,
			"password":   usr.Password,
			"is_admin":   usr.IsAdmin,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

// getUserById godoc
//
//	@Summary		Get User By ID
//	@Tags			Users
//	@Description	Get a user by ID
//	@ID				get-user-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200	{object}	map[string]interface{}	"data"
//	@Failure		500	{object}	map[string]interface{}	"error"
//	@Router			/api/users/{id} [get]
func (s *Services) getUserById(c echo.Context) error {
	userId := c.Param("id")
	usr, err := s.user.GetByID(s.ctx, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get user",
			"err":     err,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": echo.Map{
			"id":         usr.ID,
			"name":       usr.Name,
			"updated_at": usr.UpdatedAt,
			"login":      usr.Login,
			"password":   usr.Password,
			"is_admin":   usr.IsAdmin,
		},
	})
}

// getSales godoc
//
//	@Summary		Get Sales
//	@Tags			Sales
//	@Description	Get all sales
//	@ID				get-sales
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"data"
//	@Failure		500	{object}	map[string]interface{}	"error"
//	@Router			/api/sales [get]
func (s *Services) getSales(c echo.Context) error {
	slls, err := s.sell.FindAll(s.ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get sales",
			"err":     err,
		})
	}
	res := []echo.Map{}
	for _, sll := range slls {
		res = append(res, echo.Map{
			"id":         sll.ID,
			"updated_at": sll.UpdatedAt,
			"user_id":    sll.UserId,
			"product_id": sll.ProductId,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

// getSellById godoc
//
//	@Summary		Get Sale By ID
//	@Tags			Sales
//	@Description	Get a sale by ID
//	@ID				get-sale-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Sale ID"
//	@Success		200	{object}	map[string]interface{}	"data"
//	@Failure		500	{object}	map[string]interface{}	"error"
//	@Router			/api/sales/{id} [get]
func (s *Services) getSellById(c echo.Context) error {
	sellId := c.Param("id")
	sll, err := s.sell.GetByID(s.ctx, sellId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get sell",
			"err":     err,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": echo.Map{
			"id":         sll.ID,
			"updated_at": sll.UpdatedAt,
			"user_id":    sll.UserId,
			"product_id": sll.ProductId,
		},
	})
}

// getProducts godoc
//
//	@Summary		Get Products
//	@Tags			Products
//	@Description	Get all products
//	@ID				get-products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"data"
//	@Failure		500	{object}	map[string]interface{}	"error"
//	@Router			/api/products [get]
func (s *Services) getProducts(c echo.Context) error {
	products, err := s.product.FindAll(s.ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get products",
			"err":     err,
		})
	}
	res := []echo.Map{}
	for _, product := range products {
		res = append(res, echo.Map{
			"id":         product.ID,
			"updated_at": product.CreatedAt,
			"name":       product.Name,
			"price":      product.Price,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

// getProductById godoc
//
//	@Summary		Get Product By ID
//	@Tags			Products
//	@Description	Get a product by ID
//	@ID				get-product-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Product ID"
//	@Success		200	{object}	map[string]interface{}	"data"
//	@Failure		500	{object}	map[string]interface{}	"error"
//	@Router			/api/products/{id} [get]
func (s *Services) getProductById(c echo.Context) error {
	productId := c.Param("id")
	product, err := s.product.GetByID(s.ctx, productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get sale",
			"err":     err,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": echo.Map{
			"id":         product.ID,
			"updated_at": product.CreatedAt,
			"name":       product.Name,
			"price":      product.Price,
		},
	})
}

type Interval struct {
	interval string `query:"interval"`
}

// getSalesDate godoc
//
//	@Summary		Get Sales by Interval
//	@Tags			Sales
//	@Description	Get sales data filtered by a time interval
//	@ID				get-sales-date
//	@Accept			json
//	@Produce		json
//	@Param			interval	query	string	true	"Time Interval"
//	@Success		200	{object}	map[string]interface{}	"data"
//	@Failure		400	{string}	string	"bad request"
//	@Failure		500	{object}	map[string]interface{}	"error"
//	@Router			/sales/interval [get]
func (s *Services) getSalesDate(c echo.Context) error {
	var interval Interval
	err := c.Bind(&interval)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	slls, err := s.sell.SelectByTime(s.ctx, interval.interval)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get sales",
			"err":     err,
		})
	}
	res := []echo.Map{}
	for _, sll := range slls {
		res = append(res, echo.Map{
			"id":         sll.ID,
			"updated_at": sll.UpdatedAt,
			"user_id":    sll.UserId,
			"product_id": sll.ProductId,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}
