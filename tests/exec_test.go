//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	_ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestExec_CreateAndInsert(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO users (name) VALUES (?)`, "alice")
	if err != nil {
		t.Fatal(err)
	}

	var name string
	err = db.QueryRow(`SELECT name FROM users WHERE id=1`).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if name != "alice" {
		t.Fatalf("got %q, want %q", name, "alice")
	}
}
