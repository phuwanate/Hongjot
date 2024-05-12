package api

import (
	"database/sql"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/eslip"
	"github.com/KKGo-Software-engineering/workshop-summer/api/health"
	"github.com/KKGo-Software-engineering/workshop-summer/api/mlog"
	"github.com/KKGo-Software-engineering/workshop-summer/api/spender"
	"github.com/KKGo-Software-engineering/workshop-summer/api/transaction"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	*echo.Echo
}

type Postgres struct {
	Db *sql.DB
}

func New(db *sql.DB, cfg config.Config, logger *zap.Logger) *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(mlog.Middleware(logger))

	v1 := e.Group("/api/v1")

	v1.GET("/slow", health.Slow)
	v1.GET("/health", health.Check(db))
	v1.POST("/upload", eslip.Upload)
	// For pre-commit
	v1.GET("/transactions", transaction.GetTransactionsHandler(db))

	{
		h := spender.New(cfg.FeatureFlag, db)
		v1.GET("/spenders", h.GetAll)
		v1.POST("/spenders", h.Create)
		v1.GET("/spenders/:id", h.GetByID)
	}

	{
		h := transaction.New(cfg.FeatureFlag, &transaction.Postgres{Db: db})
		v1.GET("/spenders/:id/transactions", h.GetTransactionDetailBySpenderIdHandler)
		v1.GET("/spenders/:id/transactions/summary", h.GetTransactionSummaryBySpenderIdHandler)
	}

	{
		h := transaction.NewHandler(cfg.FeatureFlag, db)
		v1.POST("/transactions", h.Create)
		v1.PUT("/transactions/:id", h.Update)
	}

	return &Server{e}
}
