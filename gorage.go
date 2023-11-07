package Gorage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Gorage struct {
	Path   string
	Tables []GorageTable
}

func (g *Gorage) FromTable(name string) *GorageTable {
	k := -1
	for i, v := range g.Tables {
		if v.Name == name {
			k = i
			//return &v
			break
		}
	}
	return &g.Tables[k]
}

func (g *Gorage) TableExists(name string) bool {
	for _, v := range g.Tables {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (g *Gorage) AddTable(name string, columns []string) {
	if g.TableExists(name) {
		return
	}
	t := GorageTable{
		Name:    name,
		Columns: columns,
		Rows:    [][]interface{}{},
	}
	g.Tables = append(g.Tables, t)
}

func (g *Gorage) Save() {
	err := os.Truncate(g.Path, 0)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%p\n", g.FromTable("User"))
	file, _ := json.MarshalIndent(g, "", " ")
	err = os.WriteFile(g.Path, file, 0644)
	if err != nil {
		panic(err.Error())
	}
}

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
	return &g
}

func CreateNewGorage(path string) *Gorage {
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
			Path:   path,
			Tables: []GorageTable{},
		}
		file, _ := json.MarshalIndent(g, "", " ")
		err = os.WriteFile(path, file, 0644)
		if err != nil {
			panic(err.Error())
		}
	}
	return OpenGorage(path)
}
