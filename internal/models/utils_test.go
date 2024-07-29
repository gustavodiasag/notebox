package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// Because the setup and teardown scripts contain multiple SQL statements, the parameter
	// `multiStatements=true` is used in the DSN. This instructs the database driver to support
	// executing multiple statements in one `db.Exec()` call.
	db, err := sql.Open("mysql", "test_web:pass@/test_notebox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Registers a function which will automatically be called when the current test (or sub-test)
	// that calls `newTestDB` has finished.
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})

	return db
}
