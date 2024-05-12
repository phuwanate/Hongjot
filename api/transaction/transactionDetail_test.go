package transaction

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactionDetailBySpenderId(t *testing.T) {
	t.Run("get transaction detail by spender id", func(t *testing.T) {
		//create a new echo instance
		e := echo.New()
		defer e.Close()

		//create a new http request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/spenders/1/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		StubTxDetailStorer := MockStubData()

		h := New(config.FeatureFlag{}, StubTxDetailStorer)
		err := h.GetTransactionDetailBySpenderIdHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"transactions": [
				{
					"id": "1",
					"date": "2024-04-30T09:00:00.000Z",
					"amount": 1000,
					"category": "Food",
					"transaction_type": "expense",
					"spender_id": 1,
					"note": "Lunch",
					"image_url": "https://example.com/image1.jpg"
				},
				{
					"id": "2",
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

		}`, rec.Body.String())

	})

	t.Run("get transaction detail by spender id and invalid page", func(t *testing.T) {
		//create a new echo instance
		e := echo.New()
		defer e.Close()

		/*
			page=NotInt&limit=10
		*/
		//create a new http request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/spenders/1/transactions?page=NotInt", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		StubTxDetailStorer := MockStubData()

		h := New(config.FeatureFlag{}, StubTxDetailStorer)
		err := h.GetTransactionDetailBySpenderIdHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `"Please check your page number"`, rec.Body.String())
	})

	t.Run("get transaction detail by spender id and invalid limit", func(t *testing.T) {
		//create a new echo instance
		e := echo.New()
		defer e.Close()

		/*
			page=NotInt&limit=10
		*/
		//create a new http request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/spenders/1/transactions?page=1&limit=NotInt", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		StubTxDetailStorer := MockStubData()

		h := New(config.FeatureFlag{}, StubTxDetailStorer)
		err := h.GetTransactionDetailBySpenderIdHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `"Please check your page limit"`, rec.Body.String())
	})
}

func MockStubData() StubTxDetailStorer {
	return StubTxDetailStorer{
		txDetail: TransactionWithDetail{
			Transactions: []Transaction{
				{
					ID:              "1",
					Date:            "2024-04-30T09:00:00.000Z",
					Amount:          1000,
					Category:        "Food",
					TransactionType: "expense",
					SpenderID:       1,
					Note:            "Lunch",
					ImageURL:        "https://example.com/image1.jpg",
				},
				{
					ID:              "2",
					Date:            "2024-04-29T19:00:00.000Z",
					Amount:          2000,
					Category:        "Transport",
					TransactionType: "income",
					SpenderID:       1,
					Note:            "Salary",
					ImageURL:        "https://example.com/image2.jpg",
				},
			},
			Summary: TransactionSummary{},
			Pagination: PaginationInfo{
				CurrentPage: 1,
				TotalPages:  1,
				PerPage:     10,
			},
		},
		txSummary: TransactionSummary{
			TotalIncome:    2000,
			TotalExpenses:  1000,
			CurrentBalance: 1000,
		},
	}
}

func TestGetTransactionSummaryBySpenderId(t *testing.T) {
	t.Run("get transaction summary by spender id", func(t *testing.T) {
		//create a new echo instance
		e := echo.New()
		defer e.Close()

		//create a new http request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/spenders/1/transactions/summary", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		StubTxDetailStorer := StubTxDetailStorer{
			txSummary: TransactionSummary{
				TotalIncome:    2000,
				TotalExpenses:  1000,
				CurrentBalance: 1000,
			},
		}

		h := New(config.FeatureFlag{}, StubTxDetailStorer)
		err := h.GetTransactionSummaryBySpenderIdHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"total_income": 2000,
			"total_expenses": 1000,
			"current_balance": 1000
		}`, rec.Body.String())

	})
}

type StubTxDetailStorer struct {
	txDetail  TransactionWithDetail
	txSummary TransactionSummary
}

func (s StubTxDetailStorer) GetTransactionDetailBySpenderId(ctx context.Context, id string, offset int, limit int) (TransactionWithDetail, error) {
	return s.txDetail, nil
}

func (s StubTxDetailStorer) GetTransactionSummaryBySpenderId(ctx context.Context, id string) (TransactionSummary, error) {
	return s.txSummary, nil
}

//=================================================================================================
// SQL Mock

func TestGetTransactionDetailBySpenderIdWithSQLMock(t *testing.T) {
	t.Run("get transaction detail by spender id", func(t *testing.T) {
		//create a new echo instance
		e := echo.New()
		defer e.Close()

		//create a new http request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/spenders/1/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "spender_id", "note", "image_url"}).
			AddRow("1", "2024-04-30T09:00:00.000Z", 1000, "Food", "expense", 1, "Lunch", "https://example.com/image1.jpg").
			AddRow("2", "2024-04-29T19:00:00.000Z", 2000, "Transport", "income", 1, "Salary", "https://example.com/image2.jpg")
		mock.ExpectQuery(`SELECT id,date,amount,category, transaction_type,spender_id, note, image_url FROM transaction WHERE spender_id = $1 OFFSET $2 LIMIT $3`).WithArgs("1", 0, 10).WillReturnRows(rows)

		rowCount := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(`SELECT COUNT(*) FROM transaction WHERE spender_id = $1`).WithArgs("1").WillReturnRows(rowCount)

		rowsSummary := sqlmock.NewRows([]string{"total_income", "total_expenses", "current_balance"}).
			AddRow(2000, 1000, 1000)
		mock.ExpectQuery(`SELECT SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) AS total_income, SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) AS total_expenses, SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE -amount END) AS current_balance FROM transaction WHERE spender_id = $1`).WithArgs("1").WillReturnRows(rowsSummary)

		h := New(config.FeatureFlag{}, &Postgres{Db: db})
		err := h.GetTransactionDetailBySpenderIdHandler(c)

		tmp := rec.Body.String()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"transactions": [
				{
					"id": "1",
					"date": "2024-04-30T09:00:00.000Z",
					"amount": 1000,
					"category": "Food",
					"transaction_type": "expense",
					"spender_id": 1,
					"note": "Lunch",
					"image_url": "https://example.com/image1.jpg"
				},
				{
					"id": "2",
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
		}`, tmp)
	})

}

func TestGetTransactionSummaryBySpenderIdWithSQLMock(t *testing.T) {
	t.Run("get transaction summary by spender id", func(t *testing.T) {
		//create a new echo instance
		e := echo.New()
		defer e.Close()

		//create a new http request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/spenders/1/transactions/summary", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		rows := sqlmock.NewRows([]string{"total_income", "total_expenses", "current_balance"}).
			AddRow(2000, 1000, 1000)
		mock.ExpectQuery(`SELECT SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) AS total_income, SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) AS total_expenses, SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE -amount END) AS current_balance FROM transaction WHERE spender_id = $1`).WithArgs("1").WillReturnRows(rows)

		h := New(config.FeatureFlag{}, &Postgres{Db: db})
		err := h.GetTransactionSummaryBySpenderIdHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"total_income": 2000,
			"total_expenses": 1000,
			"current_balance": 1000
		}`, rec.Body.String())
	})
}
