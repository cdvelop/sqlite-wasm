//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	_ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestError_Constraints(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO test VALUES (1)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO test VALUES (1)`)
	if err == nil {
		t.Fatal("expected primary key constraint violation")
	}
}

func TestError_Syntax(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`INVALID SQL`)
	if err == nil {
		t.Fatal("expected syntax error")
	}
}
