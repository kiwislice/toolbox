package mssql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a successful case
func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// our sqlmock expectations
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = recordStats(db, 2, 3); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// a failing case
func TestShouldRollbackStatUpdatesOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnError(fmt.Errorf("update failed"))
	mock.ExpectRollback()

	// now we execute our method
	if err = recordStats(db, 2, 3); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Methods to test
func recordStats(db *sql.DB, userID, productID int64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
		return
	}

	if _, err = tx.Exec("INSERT INTO product_viewers (user_id, product_id) VALUES (?, ?)", userID, productID); err != nil {
		return
	}

	return
}

func TestSelectAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	cdb := &CustomDB{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "one").
		AddRow(2, "two")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM test;")).WillReturnRows(rows)

	result, err := SelectAll(cdb, "test")
	require.NoError(t, err)

	expected := []map[string]any{
		{"id": "1", "name": "one"},
		{"id": "2", "name": "two"},
	}
	assert.Equal(t, expected, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
func TestCompareTableData(t *testing.T) {
	dbL, mockL, err := sqlmock.New()
	require.NoError(t, err)
	defer dbL.Close()
	cdbL := &CustomDB{DB: dbL}

	dbR, mockR, err := sqlmock.New()
	require.NoError(t, err)
	defer dbR.Close()
	cdbR := &CustomDB{DB: dbR}

	rowsL := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "one").
		AddRow(2, "two")

	rowsR := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "one").
		AddRow(2, "two")

	mockL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM test;")).WillReturnRows(rowsL)
	mockR.ExpectQuery(regexp.QuoteMeta("SELECT * FROM test;")).WillReturnRows(rowsR)

	same := compareTableData(cdbL, cdbR, "test")
	assert.True(t, same)

	assert.NoError(t, mockL.ExpectationsWereMet())
	assert.NoError(t, mockR.ExpectationsWereMet())
}
