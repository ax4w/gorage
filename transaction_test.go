package Gorage

import (
	"os"
	"testing"
	"time"
)

func TestRollback(t *testing.T) {
	if fileExists("./transaction") {
		err := os.Remove("./transaction")
		if err != nil {
			t.Fatalf("Error removing old test file")
			return
		}
	}
	g := Create("./transaction", false, false)
	table := g.CreateTable("User")
	if table == nil {
		t.Fatalf("Table not created")
	}
	table.AddColumn("Name", STRING)
	table.AddColumn("Age", INT)
	table.Insert([]interface{}{"Carl", 20})
	table.Update(map[string]interface{}{
		"Name": "Bob",
		"Age":  "30",
	})
	res := table.Select([]string{"Name"})
	rowZero := res.Rows[0]
	if rowZero[0].(string) != "Carl" {
		t.Fatalf("Rollback failed")
	}
	g.Save()
}

func TestAll(t *testing.T) {
	if fileExists("./transaction") {
		err := os.Remove("./transaction")
		if err != nil {
			t.Fatalf("Error removing old test file")
			return
		}
	}
	g := Create("./transaction", false, false)
	table := g.CreateTable("Test")
	go func() {
		time.Sleep(5 * time.Second)
		table.AddColumn("Test1", INT)
	}()
	go func() {
		time.Sleep(5 * time.Second)
		table.AddColumn("Test1", INT)
	}()
	go func() {
		time.Sleep(5 * time.Second)
		table.AddColumn("Test", INT)
	}()
	go func() {
		time.Sleep(3 * time.Second)
		table.AddColumn("Moin", STRING)
	}()
	go func() {
		time.Sleep(4 * time.Second)
		table.Insert([]interface{}{0, "nice"})
	}()
	go func() {
		time.Sleep(2 * time.Second)
		table.AddColumn("Test", INT)
	}()
	time.Sleep(10 * time.Second)
	te := table.Select([]string{"Test", "Moin", "Test1"})
	if len(te.Rows) != 1 {
		t.Fatalf("Expected 1 row")
	}
	r := te.Rows[0]
	_ = r[0].(int)
	_ = r[1].(string)
	g.Close()

}
