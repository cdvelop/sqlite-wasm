//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	_ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestTx_CommitAndRollback(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE test (val TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(`INSERT INTO test (val) VALUES (?)`, "A")
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	var count int
	db.QueryRow(`SELECT count(*) FROM test`).Scan(&count)
	if count != 1 {
		t.Fatalf("after commit: got count %d, want 1", count)
	}

	tx, err = db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(`INSERT INTO test (val) VALUES (?)`, "B")
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}

	db.QueryRow(`SELECT count(*) FROM test`).Scan(&count)
	if count != 1 {
		t.Fatalf("after rollback: got count %d, want 1", count)
	}
}
