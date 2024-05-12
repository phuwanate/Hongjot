package transaction

import (
	"context"
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	flag   config.FeatureFlag
	storer TxDetailStorer
}

type TxDetailStorer interface {
	GetTransactionDetailBySpenderId(ctx context.Context, id string, offset int, limit int) (TransactionWithDetail, error)
	GetTransactionSummaryBySpenderId(ctx context.Context, id string) (TransactionSummary, error)
}

func New(cfg config.FeatureFlag, storer TxDetailStorer) *handler {
	return &handler{cfg, storer}
}

//=========================================================
//GET /api/v1/spenders/{id}/transactions
/*
{
	"transactions": [
		{
			"id": 1,
			"date": "2024-04-30T09:00:00.000Z",
			"amount": 1000,
			"category": "Food",
			"transaction_type": "expense",
			"spender_id": 1,
			"note": "Lunch",
			"image_url": "https://example.com/image1.jpg"
		},
		{
			"id": 2,
			"date": "2024-04-29T19:00:00.000Z",
			"amount": 2000,
			"category": "Transport",
			"transaction_type": "income",
			"spender_id": 1,
			"note": "Salary",
			"image_url": "https://example.com/image2.jpg"
		}
	],
	"summary": {
		"total_income": 2000,
		"total_expenses": 1000,
		"current_balance": 1000
	},
	"pagination": {
		"current_page": 1,
		"total_pages": 1,
		"per_page": 10
	}
}
*/

func (h handler) GetTransactionDetailBySpenderIdHandler(c echo.Context) error {

	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	//Get page number
	rawPage := c.QueryParam("page")
	page := 1
	if rawPage != "" {
		var err error
		page, err = strconv.Atoi(rawPage)
		if err != nil {
			logger.Error("bad request", zap.Error(err))
			return c.JSON(http.StatusBadRequest, "Please check your page number")
		}
	}

	//Get limit
	rawLimit := c.QueryParam("limit")
	limit := 10
	if rawLimit != "" {
		var err error
		limit, err = strconv.Atoi(rawLimit)
		if err != nil {
			logger.Error("bad request", zap.Error(err))
			return c.JSON(http.StatusBadRequest, "Please check your page limit")
		}
	}

	txDetail, err := h.storer.GetTransactionDetailBySpenderId(ctx, id, page, limit)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Please check server logs")
	}

	txSum, errTxSum := h.storer.GetTransactionSummaryBySpenderId(ctx, id)
	if errTxSum != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Please check server logs")
	}
	txDetail.Summary = txSum
	return c.JSON(http.StatusOK, txDetail)

}

type Postgres struct {
	Db *sql.DB
}

func (p *Postgres) GetTransactionDetailBySpenderId(ctx context.Context, id string, page int, limit int) (TransactionWithDetail, error) {

	//Query
	//SELECT * FROM transaction WHERE spender_id = id
	skip := (page - 1) * limit
	//OFFSET 0
	//LIMIT 10
	rows, err := p.Db.QueryContext(ctx, `SELECT id,date,amount,category, transaction_type,spender_id, note, image_url FROM transaction WHERE spender_id = $1 OFFSET $2 LIMIT $3`, id, skip, limit)
	if err != nil {

		return TransactionWithDetail{}, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(&tx.ID, &tx.Date, &tx.Amount, &tx.Category, &tx.TransactionType, &tx.SpenderID, &tx.Note, &tx.ImageURL)
		if err != nil {
			return TransactionWithDetail{}, err
		}
		txs = append(txs, tx)
	}

	//Count total pages
	var total int
	errCountTx := p.Db.QueryRowContext(ctx, `SELECT COUNT(*) FROM transaction WHERE spender_id = $1`, id).Scan(&total)
	if errCountTx != nil {
		return TransactionWithDetail{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return TransactionWithDetail{
		Transactions: txs,
		Summary:      TransactionSummary{},
		Pagination:   PaginationInfo{page, totalPages, limit},
	}, nil
}

// =========================================================
// GET /api/v1/spenders/{id}/transactions/summary
func (h handler) GetTransactionSummaryBySpenderIdHandler(c echo.Context) error {

	logger := mlog.L(c)
	ctx := c.Request().Context()

	id := c.Param("id")

	txSummary, err := h.storer.GetTransactionSummaryBySpenderId(ctx, id)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Please check server logs")
	}

	return c.JSON(http.StatusOK, txSummary)
}

func (p *Postgres) GetTransactionSummaryBySpenderId(ctx context.Context, id string) (TransactionSummary, error) {

	//Query
	//SELECT * FROM transaction WHERE spender_id = id
	rows, err := p.Db.QueryContext(ctx, `SELECT SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) AS total_income, SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) AS total_expenses, SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE -amount END) AS current_balance FROM transaction WHERE spender_id = $1`, id)
	if err != nil {
		return TransactionSummary{}, err
	}
	defer rows.Close()

	var totalIncome float64
	var totalExpenses float64
	var currentBalance float64
	for rows.Next() {
		err := rows.Scan(&totalIncome, &totalExpenses, &currentBalance)
		if err != nil {
			return TransactionSummary{}, err
		}
	}

	return TransactionSummary{
		TotalIncome:    totalIncome,
		TotalExpenses:  totalExpenses,
		CurrentBalance: totalIncome - totalExpenses,
	}, nil
}
