//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	_ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestStmt_ExecAndQuery(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := db.Prepare(`INSERT INTO test (val) VALUES (?)`)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec("hello")
	if err != nil {
		t.Fatal(err)
	}

	stmt2, err := db.Prepare(`SELECT val FROM test WHERE id = ?`)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt2.Close()

	var val string
	err = stmt2.QueryRow(1).Scan(&val)
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Fatalf("got %q, want %q", val, "hello")
	}
}
