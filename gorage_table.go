package Gorage

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	INT     = 0
	STRING  = 1
	BOOLEAN = 2
	FLOAT   = 3
	DATE    = 4
)

type Column struct {
	Name     string
	Datatype int
}

type transaction struct {
	q *queue
}

type Table struct {
	Name    string
	Columns []Column
	Rows    [][]interface{}
	host    *Gorage
	t       transaction
	p       bool
}

func (g *Table) isDuplicate(hash uint32) bool {
	for _, v := range g.Rows {
		if hash == computeHash(v) {
			return true
		}
	}
	return false
}

func (g *Table) getColByType(t int) *Column {
	if len(g.Columns) == 0 {
		return nil
	}
	for _, v := range g.Columns {
		if v.Datatype == t {
			return &v
		}
	}
	return nil
}

func (g *Table) sendExit() {
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action: actionExit,
		})
	} else {
		g.t.q.append(&data{
			action: actionExit,
		})
	}
}

func (g *Table) getColAndIndexByName(name string) (*Column, int) {
	if len(g.Columns) == 0 {
		return nil, -1
	}
	for i, v := range g.Columns {
		if name == v.Name {
			if g.host.Log {
				gprint("getColAndIndexByName", "Column: "+name+" found")
			}
			return &v, i
		}
	}
	if g.host.Log {
		gprint("getColAndIndexByName", "Column: "+name+"  was not found")
	}
	return nil, -1
}

/*
the name is the column name
*/
func (g *Table) RemoveColumn(name string) *Table {
	if !g.p {
		return g.removeColumn(name)
	}
	c := make(chan *Table)
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action:  actionDeleteColumn,
			payload: []interface{}{name},
			c:       c,
		})
	} else {
		g.t.q.append(&data{
			action:  actionDeleteColumn,
			payload: []interface{}{name},
			c:       c,
		})
	}
	for {
		select {
		case e := <-c:
			close(c)
			return e
		}
	}

}

func (g *Table) removeColumn(name string) *Table {
	c, idx := g.getColAndIndexByName(name)
	if c == nil {
		return g
	}
	g.Columns = append(g.Columns[:idx], g.Columns[idx+1:]...)
	//remove cells
	for i := 0; i < len(g.Rows); i++ {
		g.Rows[i] = append(g.Rows[i][:idx], g.Rows[i][idx+1:]...)
	}
	return g
}

/*
name is the name of the column. The datatype can be choosen from the provieded and implemented datatypes (f.e. INT,STRING)
*/
func (g *Table) AddColumn(name string, datatype int) *Table {
	if !g.p {
		return g.addColumn(name, datatype)
	}
	c := make(chan *Table)
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action:  actionAddColumn,
			payload: []interface{}{name, datatype},
			c:       c,
		})
	} else {
		g.t.q.append(&data{
			action:  actionAddColumn,
			payload: []interface{}{name, datatype},
			c:       c,
		})
	}
	for {
		select {
		case e := <-c:
			close(c)
			return e
		}
	}
}

func (g *Table) addColumn(name string, datatype int) *Table {
	if v, _ := g.getColAndIndexByName(name); v == nil {
		g.Columns = append(g.Columns, Column{
			name,
			datatype,
		})
		if g.host.Log {
			gprint("AddColumn", "Column: "+name+" added")
		}
		if len(g.Rows) != 0 {
			if g.host.Log {
				gprint("AddColumn", "Table has Rows, filling up the holes")
			}
			for i := 0; i < len(g.Rows); i++ {
				g.Rows[i] = append(g.Rows[i], nil)
			}
			if g.host.Log {
				gprint("AddColumn", "Filled up the holes")
			}
		}

	} else {
		if g.host.Log {
			gprint("AddColumn", "Column: "+name+" was not added. Duplicate?")
		}
	}
	return g
}

/*
f is the evaluate string. See github README.md for examples
*/
func (g *Table) Where(f string) *Table {
	if !g.p {
		return g.where(f)
	}
	c := make(chan *Table)
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action:  actionWhere,
			payload: []interface{}{f},
			c:       c,
		})
	} else {
		g.t.q.append(&data{
			action:  actionWhere,
			payload: []interface{}{f},
			c:       c,
		})
	}
	for {
		select {
		case e := <-c:
			close(c)
			return e
		}
	}
}
func (g *Table) where(f string) *Table {
	res := &Table{
		Name:    g.Name,
		Columns: g.Columns,
		host:    g.host,
		Rows:    [][]interface{}{},
	}
	split := strings.Split(f, " ")
	for _, v := range g.Rows {
		var tmp []string
		for _, k := range split {
			y := strings.Split(k, ":")
			if len(y) > 1 {
				col, colIdx := g.getColAndIndexByName(y[1])
				if col == nil {
					panic("Column not found")
				}
				switch col.Datatype {
				case DATE:
					if v[colIdx] == nil {
						k = fmt.Sprintf("f")
					}
					k = fmt.Sprintf("%s", v[colIdx])
					break
				case STRING:
					if v[colIdx] == nil {
						k = fmt.Sprintf("f")
					}
					k = fmt.Sprintf("'%s'", v[colIdx])
					break
				case FLOAT, INT:
					if v[colIdx] == nil {
						k = fmt.Sprintf("f")
					} else {
						switch v[colIdx].(type) {
						case float32:
							println("IS THIS EVEN BEING USED ON MOST MACHINES?")
							break
						case float64:
							k = strconv.FormatFloat(v[colIdx].(float64), 'f', -1, 64)
							break
						default:
							k = fmt.Sprintf("%d", v[colIdx])
						}
					}

				default:
					k = fmt.Sprintf("%s", v[colIdx])
					break
				}
			}
			tmp = append(tmp, k)
		}
		q := strings.Join(tmp, " ")
		if eval(q) {
			res.Rows = append(res.Rows, v)
		}
	}
	return res
}

/*
data is a map, where the key is the column and the interace is the value.
the datatype of the interface needs to match the datatype, which the column represents
*/

func (g *Table) Update(d map[string]interface{}) *Table {
	if !g.p {
		return g.update(d)
	}
	c := make(chan *Table)
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action:  actionUpdate,
			payload: []interface{}{d},
			c:       c,
		})
	} else {
		g.t.q.append(&data{
			action:  actionUpdate,
			payload: []interface{}{d},
			c:       c,
		})
	}
	for {
		select {
		case e := <-c:
			return e
		}
	}
}

func (g *Table) Wait() {
	for g.t.q.Head() != nil {

	}
}

func (g *Table) update(data map[string]interface{}) *Table {
	//g.Lock()

	rtCopy := g.host.copyTable(g.Name)
	var rt *Table
	for i, v := range g.host.Tables {
		if v.Name == g.Name {
			rt = &g.host.Tables[i]
		}
	}
	for _, v := range g.Rows {
		for i, r := range rt.Rows {
			if computeHash(v) != computeHash(r) {
				if g.host.Log {
					gprint("Update", fmt.Sprintf("Hash not matching %d != %d", computeHash(v), computeHash(r)))
				}
				continue
			}
			for key, val := range data {
				c, index := rt.getColAndIndexByName(key)
				if c == nil || !validateDatatype(val, *c) {
					g.host.copyTableToTable(g.Name, &rtCopy)
					return rt
				}
				rt.Rows[i][index] = val
				if g.host.Log {
					gprint("Update", "Updated cell")
				}
			}
		}
	}
	return rt

}

/*
Deletes Rows
*/
func (g *Table) Delete() {
	if !g.p {
		g.delete()
		return
	}
	c := make(chan *Table)
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action:  actionDelete,
			payload: []interface{}{g},
			c:       c,
		})
	} else {
		g.t.q.append(&data{
			action:  actionDelete,
			payload: []interface{}{g},
			c:       c,
		})
	}
}

func (g *Table) delete() *Table {
	realTable := g.host.FromTable(g.Name) // we need to get the table again to do persistent changes to it in memory

	if realTable == nil {
		panic("Table not found")
	}
	for idx, o := range realTable.Rows {
		for _, i := range g.Rows {
			if compareRows(o, i) {
				if idx > len(realTable.Rows) {
					return g
				}
				if idx == len(realTable.Rows) {
					realTable.Rows = append(realTable.Rows[:idx-1])
				} else {
					realTable.Rows = append(realTable.Rows[:idx], realTable.Rows[idx+1:]...)
				}
			}
		}
	}
	return g
}

/*
columns is a string array, in which the wanted columns are stored
*/

func (g *Table) Select(columns []string) *Table {
	if !g.p {
		return g._select(columns)
	}
	c := make(chan *Table)
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action:  actionSelect,
			payload: []interface{}{columns},
			c:       c,
		})
	} else {
		g.t.q.append(&data{
			action:  actionSelect,
			payload: []interface{}{columns},
			c:       c,
		})
	}
	for {
		select {
		case e := <-c:
			close(c)
			return e
		}
	}
}
func (g *Table) _select(columns []string) *Table {
	var columnIdx []int
	tmp := &Table{
		Name:    g.Name,
		Columns: []Column{},
		host:    g.host,
		Rows:    [][]interface{}{},
	}
	for _, v := range columns {
		col, i := g.getColAndIndexByName(v)
		tmp.AddColumn(col.Name, col.Datatype)
		columnIdx = append(columnIdx, i)
	}
	for _, v := range g.Rows {
		var t []interface{}
		for _, i := range columnIdx {
			if i >= len(v) {
				if g.host.Log {
					gprint("Select", "temp column index is out of bounds. skipping")
				}
				continue
			}
			t = append(t, v[i])
		}

		tmp.Rows = append(tmp.Rows, t)
	}

	return tmp
}

/*
The data is the data that shall be inserted. The len(data) needs to match the number of columns.
If a cell shall be empty you can use nil.

*Remember*: You can not compare in nil value, when using the column in a where condition
*/
func (g *Table) Insert(d []interface{}) {
	if !g.p {
		g.insert(d)
		return
	}
	c := make(chan *Table)
	if g.t.q == nil {
		g.t.q = newQueue(&data{
			action:  actionInsert,
			payload: d,
			c:       c,
		})
	} else {
		g.t.q.append(&data{
			action:  actionInsert,
			payload: d,
			c:       c,
		})
	}
	for {
		select {
		case <-c:
			return
		}
	}
}
func (g *Table) insert(data []interface{}) *Table {
	if len(data) != len(g.Columns) {
		panic(fmt.Errorf("column count and data count are different"))
	}
	if !g.host.AllowDuplicated && g.isDuplicate(computeHash(data)) {
		if g.host.Log {
			gprint("Insert", "Data already exists in Table. Returning")
		}
		return g
	}
	for i, v := range g.Columns {
		if !validateDatatype(data[i], v) {
			return g
		}
	}
	g.Rows = append(g.Rows, data)
	return g
}
