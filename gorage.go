package Gorage

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type Gorage struct {
	AllowDuplicated bool
	Log             bool
	Path            string
	Tables          []Table
}

/*
Select a table from the loaded Gorage
*/

func (g *Gorage) FromTable(name string) *Table {
	k := -1
	for i, v := range g.Tables {
		if v.Name == name {
			k = i
			break
		}
	}
	c := make(chan *Table)
	if g.Tables[k].t.q == nil {
		g.Tables[k].t.q = newQueue(&data{
			action: actionFromTable,
			c:      c,
		})
	} else {
		g.Tables[k].t.q.append(&data{
			action: actionFromTable,
			c:      c,
		})
	}

	for {
		select {
		case <-c:
			close(c)
			return &g.Tables[k]
		}
	}
}

func (g *Gorage) copyTableToTable(name string, t *Table) {
	if !g.TableExists(name) {
		return
	}
	for i, v := range g.Tables {
		if v.Name == name {
			g.Tables[i].Columns = t.Columns
			g.Tables[i].Rows = t.Rows
		}
	}
}

func (g *Gorage) copyTable(name string) Table {
	if !g.TableExists(name) {
		return Table{}
	}
	for _, v := range g.Tables {
		if v.Name == name {
			t := Table{
				host: v.host,
				p:    v.p,
				t: transaction{
					q: v.t.q.n,
				},
			}
			for _, c := range v.Columns {
				t.Columns = append(t.Columns, c)
			}
			for _, r := range v.Rows {
				var a []interface{}
				for _, t1 := range r {
					a = append(a, t1)
				}
				t.Rows = append(t.Rows, a)
			}
			return t
		}
	}
	return Table{}
}

func (g *Gorage) RemoveTable(name string) *Gorage {
	if !g.TableExists(name) {
		return g
	}
	for i, v := range g.Tables {
		if v.Name == name {
			c := make(chan *Table)
			if g.Tables[i].t.q == nil {
				g.Tables[i].t.q = newQueue(&data{
					action: actionDeleteTable,
					c:      c,
				})
			} else {
				g.Tables[i].t.q.append(&data{
					action: actionDeleteTable,
					c:      c,
				})
			}
			select {
			case <-c:
				close(c)
				break
			}
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

func (g *Gorage) CreateTable(name string) *Table {
	if g.TableExists(name) {
		if g.Log {
			gprint("CreateTable", "Table already exists")
		}
		return nil
	}

	t := Table{
		Name:    name,
		host:    g,
		Columns: []Column{},
		Rows:    [][]interface{}{},
		p:       true,
	}
	g.Tables = append(g.Tables, t)
	go transactionManger(&g.Tables[len(g.Tables)-1])
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
func Open(path string) *Gorage {
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
		g.Tables[i].p = true
		go transactionManger(&g.Tables[i])
	}
	return &g
}

/*
Create a new Gorage

path: the path where it should be stored

allowDuplicates: If a table can contain multiple identical datasets

log: If you want to get spammed :^)
*/
func Create(path string, allowDuplicates, log bool) *Gorage {
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
			Tables:          []Table{},
		}
		file, _ := json.MarshalIndent(g, "", "	")
		err = os.WriteFile(path, file, 0644)
		if err != nil {
			panic(err.Error())
		}
	}
	return Open(path)
}

func (g *Gorage) Close() {
	var w sync.WaitGroup
	for i, _ := range g.Tables {
		g.Tables[i].sendExit()
		go func(i int) {
			defer w.Done()
			w.Add(1)
			for g.Tables[i].t.q.Head() != nil {
			}
		}(i)
	}
	w.Wait()
	g.Save()
}
