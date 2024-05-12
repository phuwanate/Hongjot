package transaction

import (
	"database/sql"
	"net/http"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handlerTransaction struct {
	flag config.FeatureFlag
	db   *sql.DB
}

const (
	cStmt = `INSERT INTO transaction (date, amount, category, transaction_type, spender_id, note, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	uStmt = `UPDATE transaction SET date = $1, amount = $2, category = $3, transaction_type = $4, spender_id = $5, note = $6, image_url = $7 WHERE id = $8;`
)

func NewHandler(cfg config.FeatureFlag, db *sql.DB) *handlerTransaction {
	return &handlerTransaction{cfg, db}
}

func (h handlerTransaction) Create(c echo.Context) error {

	logger := mlog.L(c)
	ctx := c.Request().Context()
	var trBody TransactionReqBody
	err := c.Bind(&trBody)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var insertTransactionId string

	err = h.db.QueryRowContext(ctx, cStmt, trBody.Date, trBody.Amount, trBody.Category, trBody.TransactionType, trBody.SpenderID, trBody.Note, trBody.ImageURL).Scan(&insertTransactionId)

	transaction := Transaction{
		ID:              insertTransactionId,
		Date:            trBody.Date,
		Amount:          trBody.Amount,
		Category:        trBody.Category,
		TransactionType: trBody.TransactionType,
		SpenderID:       trBody.SpenderID,
		Note:            trBody.Note,
		ImageURL:        trBody.ImageURL,
	}

	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create successfully", zap.String("id", insertTransactionId))
	return c.JSON(http.StatusCreated, transaction)
}

func (h handlerTransaction) Update(c echo.Context) error {

	logger := mlog.L(c)
	ctx := c.Request().Context()
	var trBody TransactionReqBody
	err := c.Bind(&trBody)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id := c.Param("id")
	// var updatedTransaction Transaction
	_ = h.db.QueryRowContext(ctx, uStmt, trBody.Date, trBody.Amount, trBody.Category, trBody.TransactionType, trBody.SpenderID, trBody.Note, trBody.ImageURL, id).Scan()

	transaction := Transaction{
		ID:              id,
		Date:            trBody.Date,
		Amount:          trBody.Amount,
		Category:        trBody.Category,
		TransactionType: trBody.TransactionType,
		SpenderID:       trBody.SpenderID,
		Note:            trBody.Note,
		ImageURL:        trBody.ImageURL,
	}

	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create successfully", zap.String("id", id))
	return c.JSON(http.StatusOK, transaction)
}
