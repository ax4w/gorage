package Gorage

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateFile(t *testing.T) {
	CreateNewGorage("./test.json")
	_, err := os.Stat("./test.json")
	if os.IsNotExist(err) {
		t.Fatalf("File was not created")
	}
}

func TestCreateTable(t *testing.T) {
	g := OpenGorage("./test.json")
	g.AddTable("User", []string{"FirstName", "LastName", "Age"})
	if !g.TableExists("User") {
		t.Fatalf("Table was not created")
	}
	g.Save()
}

func TestInsert(t *testing.T) {
	g := OpenGorage("./test.json")
	userTable := g.FromTable("User")
	fmt.Printf("k %p\n", userTable)
	userTable.Insert([]interface{}{"Lars goofy", "Oeli", 5})
	g.Save()

}

func TestDelete(t *testing.T) {
	g := OpenGorage("./test.json")
	g.
		FromTable("User").
		Where("s:FirstName = 'Lars goofy' & i:Age = '5'").
		Delete()
	_ = g.
		FromTable("User").
		Where("s:FirstName = 'Lars goofy' & i:Age = '5'").
		Select([]string{"FirstName", "LastName", "Age"})
	g.Save()
}

func TestWhere(t *testing.T) {
	g := OpenGorage("./test.json")
	userTable := g.
		FromTable("User").
		Where("s:FirstName = 'Lars goofy' & i:Age = '5'").
		Select([]string{"FirstName", "LastName", "Age"})
	for _, v := range userTable.Rows {
		for _, j := range v {
			switch j.(type) {
			case string:
				println(j.(string))
			case float64:
				println(int(j.(float64)))
			case int:
				println(j.(int))
			}

		}
	}
	g.Save()
}
