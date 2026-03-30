//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	"github.com/cdvelop/sqlite-wasm/driver"
)

func TestBackup(t *testing.T) {
	src, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	_, err = src.Exec(`CREATE TABLE test (val TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = src.Exec(`INSERT INTO test VALUES (?)`, "A")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := src.Conn(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = conn.Raw(func(driverConn any) error {
		_, ok := driverConn.(interface {
			NewBackup(string) (*driver.Backup, error)
		})
		if !ok {
			t.Fatal("driver.Conn does not implement NewBackup")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
