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
)

type GorageColumn struct {
	Name     string
	Datatype int
}

type GorageTable struct {
	Name    string
	Columns []GorageColumn
	Rows    [][]interface{}
	host    *Gorage
}

func (g *GorageTable) getColAndIndexByName(name string) (*GorageColumn, int) {
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
func (g *GorageTable) RemoveColumn(name string) *GorageTable {
	c, idx := g.getColAndIndexByName(name)
	if c == nil {
		return g
	}
	g.Columns = append(g.Columns[:idx], g.Columns[idx+1:]...)
	//remove cells
	for i := 0; i < len(g.Rows); i++ {
		//cpy := g.Rows[i]
		//g.Rows[i] = []interface{}{}
		g.Rows[i] = append(g.Rows[i][:idx], g.Rows[i][idx+1:]...)
	}
	return g
}

/*
name is the name of the column. The datatype can be choosen from the provieded and implemented datatypes (f.e. INT,STRING)
*/
func (g *GorageTable) AddColumn(name string, datatype int) *GorageTable {
	if v, _ := g.getColAndIndexByName(name); v == nil {
		g.Columns = append(g.Columns, GorageColumn{
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
f is the eval string. See github README.md for examples
*/
func (g *GorageTable) Where(f string) *GorageTable {
	res := &GorageTable{
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
		if e := eval(toTree(parse(q))[0]); e != nil && string(e.value) == "t" {
			res.Rows = append(res.Rows, v)
		}
	}
	return res
}

/*
data is a map, where the key is the column and the interace is the value.
the datatype of the interface needs to match the datatype, which the column represents
*/
func (g *GorageTable) Update(data map[string]interface{}) {
	rt := g.host.FromTable(g.Name) // we need to get the table again to do persistent changes to it in memory
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
					panic("No matching column found or mismatch datatype")
				}
				rt.Rows[i][index] = val
				if g.host.Log {
					gprint("Update", "Updated cell")
				}
			}
		}
	}
}

/*
Deletes Rows
*/
func (g *GorageTable) Delete() {
	realTable := g.host.FromTable(g.Name) // we need to get the table again to do persistent changes to it in memory
	if realTable == nil {
		panic("Table not found")
	}
	for idx, o := range realTable.Rows {
		for _, i := range g.Rows {
			if compareRows(o, i) {
				if idx+1 > len(realTable.Rows) {
					realTable.Rows = append(realTable.Rows[idx:])
				} else {
					realTable.Rows = append(realTable.Rows[:idx], realTable.Rows[idx+1:]...)
				}
			}
		}
	}
}

/*
columns is a string array, in which the wanted columns are stored
*/
func (g *GorageTable) Select(columns []string) *GorageTable {
	var columnIdx []int
	tmp := &GorageTable{
		Name:    g.Name,
		Columns: []GorageColumn{},
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

func (g *GorageTable) isDuplicate(hash uint32) bool {
	for _, v := range g.Rows {
		if hash == computeHash(v) {
			return true
		}
	}
	return false
}

/*
The data is the data that shall be inserted. The len(data) needs to match the number of columns.
If a cell shall be empty you can use nil.

*Remember*: You can not compare in nil value, when using the column in a where condition
*/
func (g *GorageTable) Insert(data []interface{}) {
	if len(data) != len(g.Columns) {
		panic(fmt.Errorf("column count and data count are different"))
	}
	if !g.host.AllowDuplicated && g.isDuplicate(computeHash(data)) {
		if g.host.Log {
			gprint("Insert", "Data already exists in Table. Returning")
		}
		return
	}
	for i, v := range g.Columns {
		validateDatatype(data[i], v)
	}
	g.Rows = append(g.Rows, data)
}
