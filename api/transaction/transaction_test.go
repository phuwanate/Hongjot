package transaction

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// For pre-commit
func TestGetTransactionsHandler(t *testing.T) {
	// Create a new instance of Echo
	e := echo.New()

	// Mock SQL database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define expectations for SQL mock
	rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "spender_id", "note", "image_url"}).
		AddRow(1, time.Now(), 100.0, "Food", "expense", 1, "Dinner out", "http://example.com/receipt.jpg").
		AddRow(2, time.Now(), 200.0, "Salary", "income", 1, "Monthly salary", "http://example.com/salary.jpg")

	mock.ExpectQuery("^SELECT (.+) FROM \"transaction\" WHERE").WillReturnRows(rows)
	mock.ExpectQuery("^SELECT COUNT\\(\\*\\) FROM \"transaction\" WHERE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Create a request to pass to our handler
	req := httptest.NewRequest(http.MethodGet, "/?page=1&limit=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assign the mocked DB to the handler function
	h := GetTransactionsHandler(db)

	// Run the handler
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, strings.Contains(rec.Body.String(), "Food"))
		assert.True(t, strings.Contains(rec.Body.String(), "Salary"))
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
