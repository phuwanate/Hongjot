package transaction

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Transaction represents the structure of a transaction record.
type Transaction struct {
	ID              string  `json:"id"`
	Date            string  `json:"date"`
	Amount          float64 `json:"amount"`
	Category        string  `json:"category"`
	TransactionType string  `json:"transaction_type"`
	SpenderID       int     `json:"spender_id"`
	Note            string  `json:"note"`
	ImageURL        string  `json:"image_url"`
}

// ResponseData includes transactions array, summary, and pagination details.
type ResponseData struct {
	Transactions []Transaction      `json:"transactions"`
	Summary      TransactionSummary `json:"summary"`
	Pagination   PaginationInfo     `json:"pagination"`
}

type TransactionSummary struct {
	TotalIncome    float64 `json:"total_income"`
	TotalExpenses  float64 `json:"total_expenses"`
	CurrentBalance float64 `json:"current_balance"`
}

type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	PerPage     int `json:"per_page"`
}

type TransactionWithDetail struct {
	Transactions []Transaction      `json:"transactions"`
	Summary      TransactionSummary `json:"summary"`
	Pagination   PaginationInfo     `json:"pagination"`
}

type TransactionReqBody struct {
	Date            string  `json:"date"`
	Amount          float64 `json:"amount"`
	Category        string  `json:"category"`
	TransactionType string  `json:"transaction_type"`
	SpenderID       int     `json:"spender_id"`
	Note            string  `json:"note"`
	ImageURL        string  `json:"image_url"`
}

// For pre-commit
// GetTransactionsHandler returns a handler function to fetch transactions with optional pagination and filtering.
func GetTransactionsHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse and validate pagination parameters
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 {
			limit = 10
		}

		// Fetch and validate filter parameters
		date := c.QueryParam("date")
		amount := c.QueryParam("amount")
		category := c.QueryParam("category")
		transactionType := c.QueryParam("transaction_type")

		// SQL query construction with filters
		query := `SELECT id, date, amount, category, transaction_type, spender_id, note, image_url FROM "transaction"`
		whereClauses := []string{"TRUE"}
		if date != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("DATE(date) = '%s'", date))
		}
		if amount != "" {
			if _, err := strconv.ParseFloat(amount, 64); err == nil {
				whereClauses = append(whereClauses, fmt.Sprintf("amount = %s", amount))
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid amount format")
			}
		}
		if category != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("category = '%s'", category))
		}
		if transactionType != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("transaction_type = '%s'", transactionType))
		}

		// Execute the filtered query
		filteredQuery := fmt.Sprintf("%s WHERE %s ORDER BY id LIMIT %d OFFSET %d", query, strings.Join(whereClauses, " AND "), limit, (page-1)*limit)
		rows, err := db.Query(filteredQuery)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch transactions: %s", err.Error()))
		}
		defer rows.Close()

		// Process query results
		var transactions []Transaction
		var totalIncome, totalExpenses float64
		for rows.Next() {
			var t Transaction
			// var date sql.NullTime
			var amount sql.NullFloat64
			if err := rows.Scan(&t.ID, &t.Date, &amount, &t.Category, &t.TransactionType, &t.SpenderID, &t.Note, &t.ImageURL); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error scanning transaction: %s", err.Error()))
			}

			// Populate valid data fields
			if amount.Valid {
				t.Amount = amount.Float64
			}

			transactions = append(transactions, t)
			// Update income and expenses based on transaction type
			if strings.ToLower(t.TransactionType) == "income" {
				totalIncome += t.Amount
			} else if strings.ToLower(t.TransactionType) == "expense" {
				totalExpenses += t.Amount
			}
		}

		// Handle post-query errors
		if err = rows.Err(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error fetching transactions: %s", err.Error()))
		}

		// Calculate total pages for pagination
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM \"transaction\" WHERE %s", strings.Join(whereClauses, " AND "))
		var totalRecords int
		if err = db.QueryRow(countQuery).Scan(&totalRecords); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to count transactions: %s", err.Error()))
		}
		totalPages := (totalRecords + limit - 1) / limit

		// Prepare and return the API response
		response := ResponseData{
			Transactions: transactions,
			Summary: TransactionSummary{
				TotalIncome:    totalIncome,
				TotalExpenses:  totalExpenses,
				CurrentBalance: totalIncome - totalExpenses,
			},
			Pagination: PaginationInfo{
				CurrentPage: page,
				TotalPages:  totalPages,
				PerPage:     limit,
			},
		}

		return c.JSON(http.StatusOK, response)
	}
}
