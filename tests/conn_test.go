//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	_ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestOpen_InMemory(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestOpen_InvalidPath(t *testing.T) {
	db, err := sql.Open("sqlite", "/nonexistent/path/db.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	// Ping should fail for invalid path
	if err := db.Ping(); err == nil {
		t.Fatal("expected error for invalid path")
	}
}
