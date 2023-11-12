package Gorage

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type Gorage struct {
	sync.Mutex
	AllowDuplicated bool
	Log             bool
	Path            string
	Tables          []GorageTable
}

/*
Select a table from the loaded Gorage
*/

func (g *Gorage) FromTable(name string) *GorageTable {
	k := -1
	for i, v := range g.Tables {
		if v.Name == name {
			k = i
			break
		}
	}
	return &g.Tables[k]
}

func (g *Gorage) RemoveTable(name string) *Gorage {
	if !g.TableExists(name) {
		return g
	}
	for i, v := range g.Tables {
		if v.Name == name {
			g.Tables = append(g.Tables[:i], g.Tables[i+1:]...)
		}
	}
	return g
}

/*
Check if a given table exists
*/
func (g *Gorage) TableExists(name string) bool {
	for _, v := range g.Tables {
		if v.Name == name {
			return true
		}
	}
	return false
}

/*
Create a table.

Two tables with the same name in the same gorage are NOT possible
*/

func (g *Gorage) CreateTable(name string) *GorageTable {
	if g.TableExists(name) {
		if g.Log {
			gprint("CreateTable", "Table already exists")
		}
		return nil
	}
	t := GorageTable{
		Name:    name,
		host:    g,
		Columns: []GorageColumn{},
		Rows:    [][]interface{}{},
	}
	g.Tables = append(g.Tables, t)

	return &g.Tables[len(g.Tables)-1]
}

/*
Save the loaded gorage
*/
func (g *Gorage) Save() {
	err := os.Truncate(g.Path, 0)
	if err != nil {
		panic(err.Error())
	}
	file, _ := json.MarshalIndent(g, "", " ")
	err = os.WriteFile(g.Path, file, 0644)
	if err != nil {
		panic(err.Error())
	}
}

/*
Open a gorage from a path
*/
func OpenGorage(path string) *Gorage {
	f, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	var g Gorage
	err = json.Unmarshal(b, &g)
	if err != nil {
		panic(err.Error())
	}
	for i, _ := range g.Tables {
		g.Tables[i].host = &g
	}
	go func() {
		for _, v := range g.Tables {
			if c := v.getColByType(TIMEOUT); c != nil {

			}
		}
	}()
	return &g
}

/*
Create a new Gorage

path: the path where it should be stored

allowDuplicates: If a table can contain multiple identical datasets

log: If you want to get spammed :^)
*/
func CreateNewGorage(path string, allowDuplicates, log bool) *Gorage {
	if !fileExists(path) {
		f, err := os.Create(path)
		if err != nil {
			panic(err.Error())
		}
		err = f.Close()
		if err != nil {
			panic(err.Error())
		}
		g := Gorage{
			Log:             log,
			AllowDuplicated: allowDuplicates,
			Path:            path,
			Tables:          []GorageTable{},
		}
		file, _ := json.MarshalIndent(g, "", "	")
		err = os.WriteFile(path, file, 0644)
		if err != nil {
			panic(err.Error())
		}
	}
	return OpenGorage(path)
}
