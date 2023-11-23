package Gorage

import (
	"os"
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
	Create("./test", false, false)
	_, err := os.Stat("./test")
	if os.IsNotExist(err) {
		t.Fatalf("File was not created")
	}
}

func TestCreateTable(t *testing.T) {
	g := Open("./test")
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
	g := Open("./test")
	userTable := g.FromTable("User")
	userTable.Insert([]interface{}{"James", "aa", 2, 85.5})
	userTable.Insert([]interface{}{"Carl", "aa", 3, 90.5})
	res := g.
		FromTable("User").
		Where(":FirstName == 'James'")
	if len(res.Rows) != 1 {
		g.Close()
		t.Fatalf("Row was not inserted")
	}
	g.Save()
}

func TestUpdate(t *testing.T) {
	g := Open("./test")
	g.FromTable("User").
		Where(":FirstName == 'James'").
		Update(map[string]interface{}{
			"FirstName": "William",
		})
	res := g.
		FromTable("User").
		Where(":FirstName == 'William'")
	if len(res.Rows) != 1 {
		g.Close()
		t.Fatalf("Row was not inserted")
	}
	g.Save()
}

func TestWhere(t *testing.T) {
	g := Open("./test")
	userTable := g.
		FromTable("User").
		Where("( :FirstName == 'William' && :Age == 2 ) || :IQ >= 90.0").
		Select([]string{"FirstName", "LastName", "Age"})
	if len(userTable.Rows) != 2 {
		g.Close()
		t.Fatalf("More than expected")
	}
	/*for _, v := range userTable.Rows {
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
	}*/
	g.Save()
}

func TestDelete(t *testing.T) {
	g := Open("./test")
	g.
		FromTable("User").
		Where(":FirstName == 'Carl'").
		Delete()
	r := g.
		FromTable("User").
		Where(":FirstName == 'Carl'")
	if len(r.Rows) != 0 {
		t.Fatalf("Delete did not work")
	}
	g.Save()
}

func TestAddColumn(t *testing.T) {
	g := Open("./test")
	g.FromTable("User").AddColumn("Lol", INT)
	g.Save()
}

func TestRemoveColumn(t *testing.T) {
	g := Open("./test")
	table := g.FromTable("User")
	l := len(table.Columns)
	table.RemoveColumn("IQ")
	if l-len(table.Columns) != 1 {
		t.Fatalf("Expeced one, got a nother value")
	}
	g.Save()
}

func TestComplete(t *testing.T) {
	if fileExists("./Social") {
		err := os.Remove("./Social")
		if err != nil {
			t.Fatalf("Error removing old test file")
			return
		}
	}
	gorage := Create("./Social", false, false)
	userTable := gorage.CreateTable("User")
	if userTable == nil {
		return
	}
	userTable.
		AddColumn("Name", STRING).
		AddColumn("Handle", STRING).
		AddColumn("Age", INT)
	gorage.Save()

	userTable.Insert([]interface{}{"Emily", "@Emily", 20})
	userTable.Insert([]interface{}{"Emily", "@Emily_Backup", 20})
	userTable.Insert([]interface{}{"Carl", "@Carl", 23})

	gorage.Save()

	userTable.
		Where(":Handle == '@Emily'").
		Update(map[string]interface{}{
			"Name": "Emily MLG",
		})

	gorage.Save()
	userTable.Where(":Handle == '@Emily_Backup' || :Name == 'Carl'").Delete()
	gorage.Save()
}

func TestCreateMemOnly(t *testing.T) {
	g := CreateMemOnly(true, false)
	tab := g.CreateTable("Test")
	tab.AddColumn("Name", STRING)
	tab.Insert([]interface{}{"Tom"})
	tab.Insert([]interface{}{"Tom"})
	g.Close()
}
