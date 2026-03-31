//go:build !wasm

package driver_test

import (
	"database/sql"
	"testing"

	"github.com/cdvelop/sqlite-wasm/driver/vtab"
)

type testModule struct{}

func (m *testModule) Create(ctx vtab.Context, args []string) (vtab.Table, error) {
	return &testTable{}, ctx.Declare("CREATE TABLE x(a, b)")
}

func (m *testModule) Connect(ctx vtab.Context, args []string) (vtab.Table, error) {
	return &testTable{}, ctx.Declare("CREATE TABLE x(a, b)")
}

type testTable struct{}

func (v *testTable) Open() (vtab.Cursor, error) {
	return &testCursor{0}, nil
}
func (v *testTable) BestIndex(i *vtab.IndexInfo) error { return nil }
func (v *testTable) Disconnect() error                 { return nil }
func (v *testTable) Destroy() error                    { return nil }

type testCursor struct {
	rowid int64
}

func (c *testCursor) Column(i int) (vtab.Value, error) {
	switch i {
	case 0:
		return c.rowid, nil
	case 1:
		return "test", nil
	}
	return nil, nil
}

func (c *testCursor) Filter(idxNum int, idxStr string, vals []vtab.Value) error {
	c.rowid = 1
	return nil
}
func (c *testCursor) Next() error {
	c.rowid++
	return nil
}
func (c *testCursor) Eof() bool {
	return c.rowid > 2
}
func (c *testCursor) Rowid() (int64, error) {
	return c.rowid, nil
}
func (c *testCursor) Close() error {
	return nil
}

func TestVTab_Register(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = vtab.RegisterModule(db, "test_vtab", &testModule{})
	if err != nil {
		t.Fatal(err)
	}

	// Re-open to pick up registered module (as per RegisterModule doc)
	db2, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()

	_, err = db2.Exec("CREATE VIRTUAL TABLE x USING test_vtab")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db2.Query("SELECT * FROM x")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 rows, got %d", count)
	}
}
