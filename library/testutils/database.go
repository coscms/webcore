package testutils

import (
	"database/sql"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/webx-top/db/lib/factory"
)

type BindSession interface {
	// BindSession sets the *sql.DB the session will use.
	BindSession(*sql.DB) error
}

func SQLMock() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New() // mock sql.DB
	if err != nil {
		panic(err)
	}
	factory.GetCluster(0).Master().(BindSession).BindSession(db)
	return mock
}
