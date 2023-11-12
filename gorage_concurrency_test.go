package Gorage

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
	if fileExists("./concurrent") {
		err := os.Remove("./concurrent")
		if err != nil {
			t.Fatalf("Error removing old test file")
			return
		}
	}
	g := CreateNewGorage("./concurrent", false, false)
	//CreateTable should not be used concurrent UNLESS it is waited on properly.
	//But it still can fail in concurrent usage. NOT RECOMMENDED
	var group sync.WaitGroup
	for i := 0; i < 10; i++ {
		g.CreateTable(fmt.Sprintf("Table%d", i)).AddColumn("Number", INT)
	}
	g.Save()
	table1 := g.FromTable("Table1")
	for i := 0; i < 5; i++ {
		group.Add(1)
		go func(k int) {
			defer group.Done()
			table1.Insert([]interface{}{k})
		}(i)
	}
	group.Wait()
	if len(table1.Rows) != 5 {
		t.Fatalf("Not all rows created")
	}
	for i := 0; i < 5; i++ {
		group.Add(1)
		go func(k int) {
			defer group.Done()
			table1 = table1.Update(map[string]interface{}{"Number": 1337})
		}(i)
	}
	group.Wait()
	g.Save()
	for _, v := range table1.Select([]string{"Number"}).Rows {
		if len(v) != 1 || v[0] != 1337 {
			t.Fatalf("update did not work")
		}
	}

	for i := 0; i < 5; i++ {
		group.Add(1)
		go func(k int) {
			defer group.Done()
			table1.Where(":Number == 1337").Delete()
		}(i)
	}
	group.Wait()
	g.Save()
	if len(table1.Rows) != 0 {
		t.Fatalf("Delete did not work")
	}
	//RemoveTable does not support concurrency, but it should work
	for i := 0; i < 5; i++ {
		group.Add(1)
		go func(k int) {
			defer group.Done()
			g.RemoveTable("Table1")
		}(i)
	}
	group.Wait()
	g.Save()
	if len(g.Tables) != 9 {
		t.Fatalf("RemoveTable did not work")
	}

}
