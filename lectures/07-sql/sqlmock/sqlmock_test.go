package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT name FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	_, err = db.Exec("SELECT name FROM users WHERE id = ?", 1)
	require.NotNil(t, err)
	require.Equal(t, err, sql.ErrNoRows)

	require.NoError(t, tx.Commit())

	require.NoError(t, mock.ExpectationsWereMet())
}
