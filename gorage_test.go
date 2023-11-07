package Gorage

import (
	"os"
	"strconv"
	"testing"
)

func TestCreateFile(t *testing.T) {
	CreateNewGorage("./test.json", false, true)
	_, err := os.Stat("./test.json")
	if os.IsNotExist(err) {
		t.Fatalf("File was not created")
	}
}

func TestCreateTable(t *testing.T) {
	g := OpenGorage("./test.json")
	table := g.CreateTable("User")
	if table != nil {
		table.AddColumn("FirstName", STRING).
			AddColumn("LastName", STRING).
			AddColumn("Age", INT).
			AddColumn("IQ", FLOAT)
	}
	if !g.TableExists("User") {
		t.Fatalf("Table was not created")
	}
	g.Save()
}

func TestInsert(t *testing.T) {
	g := OpenGorage("./test.json")
	userTable := g.FromTable("User")
	userTable.Insert([]interface{}{"Lars goofy", "Oeli", 1, nil})
	g.Save()
}

func TestDelete(t *testing.T) {
	g := OpenGorage("./test.json")
	g.
		FromTable("User").
		Where(":FirstName = 'Lars goofy' & :Age = 5").
		Delete()
	_ = g.
		FromTable("User").
		Where(":FirstName = 'Lars goofy' & :Age = 5").
		Select([]string{"FirstName", "LastName", "Age"})
	g.Save()
}
func TestWhere(t *testing.T) {
	g := OpenGorage("./test.json")
	userTable := g.
		FromTable("User").
		Where(":FirstName = 'Lars goofy' & :Age = 5 & :IQ != 85.5").
		Select([]string{"FirstName", "LastName", "Age", "IQ"})

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
