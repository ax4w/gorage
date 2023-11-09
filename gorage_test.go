package Gorage

import (
	"os"
	"strconv"
	"testing"
)

func TestCreateFile(t *testing.T) {
	if fileExists("./test") {
		err := os.Remove("./test")
		if err != nil {
			t.Fatalf("Error removing old test file")
			return
		}
	}
	CreateNewGorage("./test", false, false)
	_, err := os.Stat("./test")
	if os.IsNotExist(err) {
		t.Fatalf("File was not created")
	}
}

func TestCreateTable(t *testing.T) {
	g := OpenGorage("./test")
	table := g.CreateTable("User")
	if table != nil {
		table.AddColumn("FirstName", STRING).
			AddColumn("LastName", STRING).
			AddColumn("Age", INT).
			AddColumn("IQ", FLOAT)
	} else {
		t.Fatalf("Table was not created")
	}
	g.Save()
}

func TestInsert(t *testing.T) {
	g := OpenGorage("./test")
	userTable := g.FromTable("User")
	userTable.Insert([]interface{}{"James", "aa", 2, 85.5})
	userTable.Insert([]interface{}{"Carl", "aa", 3, 90.5})
	res := g.
		FromTable("User").
		Where(":FirstName = 'James'")
	if len(res.Rows) != 1 {
		t.Fatalf("Row was not inserted")
	}
	g.Save()
}

func TestUpdate(t *testing.T) {
	g := OpenGorage("./test")
	g.FromTable("User").
		Where(":FirstName == 'James'").
		Update(map[string]interface{}{
			"FirstName": "William",
		})
	res := g.
		FromTable("User").
		Where(":FirstName == 'William'")
	if len(res.Rows) != 1 {
		t.Fatalf("Row was not inserted")
	}
	g.Save()
}

func TestWhere(t *testing.T) {
	g := OpenGorage("./test")
	userTable := g.
		FromTable("User").
		Where(":FirstName == 'William' && :Age == 2").
		Select([]string{"FirstName", "LastName", "Age"})
	if len(userTable.Rows) != 1 {
		t.Fatalf("More than expected")
	}
	for _, v := range userTable.Rows {
		for _, j := range v {
			switch j.(type) {
			case string:
				println(j.(string))
			case float64:
				println(strconv.FormatFloat(j.(float64), 'f', -1, 64))
			case float32:
				println("UNREACHABLE")
			case int:
				println(j.(int))
			default:
				if j == nil {
					println("nil")
				}
			}
		}
	}
	g.Save()
}

func TestDelete(t *testing.T) {
	g := OpenGorage("./test")
	g.
		FromTable("User").
		Where(":FirstName = 'Carl'").
		Delete()
	r := g.
		FromTable("User").
		Where(":FirstName = 'Carl'")
	if len(r.Rows) != 0 {
		t.Fatalf("Delete did not work")
	}
	g.Save()
}

func TestRemoveColumn(t *testing.T) {
	g := OpenGorage("./test")
	table := g.FromTable("User")
	l := len(table.Columns)
	table.RemoveColumn("IQ")
	if l-len(table.Columns) != 1 {
		t.Fatalf("Expeced one, got a nother value")
	}
	g.Save()
}
