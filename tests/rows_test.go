//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	_ "github.com/cdvelop/sqlite-wasm/driver"
)

func TestRows_ScanTypes(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE test (i INTEGER, f FLOAT, s TEXT, b BLOB)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO test VALUES (?, ?, ?, ?)`, 42, 3.14, "hello", []byte{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query(`SELECT * FROM test`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("expected row")
	}

	var i int
	var f float64
	var s string
	var b []byte
	err = rows.Scan(&i, &f, &s, &b)
	if err != nil {
		t.Fatal(err)
	}

	if i != 42 || f != 3.14 || s != "hello" || len(b) != 3 || b[0] != 1 {
		t.Fatalf("unexpected row data: i=%d f=%f s=%q b=%v", i, f, s, b)
	}
}
