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
	Host    *Gorage `json:"-"`
}

func (g *GorageTable) getColAndIndexByName(name string) (*GorageColumn, int) {
	if len(g.Columns) == 0 {
		return nil, -1
	}
	for i, v := range g.Columns {
		if name == v.Name {
			if g.Host.Log {
				gprint("AddColumn", "Column: "+name+" added")
			}
			return &v, i
		}
	}
	if g.Host.Log {
		gprint("AddColumn", "Column: "+name+"  was not found")
	}
	return nil, -1
}

func (g *GorageTable) AddColumn(name string, datatype int) *GorageTable {
	if v, _ := g.getColAndIndexByName(name); v == nil {
		g.Columns = append(g.Columns, GorageColumn{
			name,
			datatype,
		})
		if g.Host.Log {
			gprint("AddColumn", "Column: "+name+" added")
		}

	} else {
		if g.Host.Log {
			gprint("AddColumn", "Column: "+name+" was not added. Duplicate?")
		}
	}
	return g
}

func (g *GorageTable) Where(f string) *GorageTable {
	res := &GorageTable{
		Name:    g.Name,
		Columns: g.Columns,
		Host:    g.Host,
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
				case INT:
					if v[colIdx] == nil {
						k = fmt.Sprintf("f")
					} else {
						switch v[colIdx].(type) {
						case float32:
							println("UNREACHABLE")
							break
						case float64:
							k = strconv.FormatFloat(v[colIdx].(float64), 'f', -1, 64)
							break
						default:
							k = fmt.Sprintf("%d", v[colIdx])
						}

					}
					break
				case FLOAT:
					if v[colIdx] == nil {
						k = fmt.Sprintf("f")
					} else {
						switch v[colIdx].(type) {
						case float32:
							println("UNREACHABLE")
							break
						case float64:
							k = strconv.FormatFloat(v[colIdx].(float64), 'f', -1, 64)
							break
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

func compareRows(a, b []interface{}) bool {
	return computeHash(a) == computeHash(b)
}

func (g *GorageTable) Delete() {
	k := -1
	for i, v := range g.Host.Tables {
		if v.Name == g.Name {
			k = i
		}
	}
	if k == -1 {
		panic("Table not found")
	}
	realTable := g.Host.Tables[k]
	for idx, o := range realTable.Rows {
		for _, i := range g.Rows {
			if compareRows(o, i) {
				if idx+1 > len(g.Host.Tables[k].Rows) {
					g.Host.Tables[k].Rows = append(
						g.Host.Tables[k].Rows[idx:],
					)
				} else {
					g.Host.Tables[k].Rows = append(
						g.Host.Tables[k].Rows[:idx],
						g.Host.Tables[k].Rows[idx+1:]...,
					)
				}

			}
		}
	}
}

func (g *GorageTable) Select(columns []string) *GorageTable {
	var columnIdx []int
	tmp := &GorageTable{
		Name:    g.Name,
		Columns: []GorageColumn{},
		Host:    g.Host,
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
				if g.Host.Log {
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

func (g *GorageTable) Insert(data []interface{}) {
	if len(data) != len(g.Columns) {
		panic(fmt.Errorf("column count and data count are different"))
	}
	if !g.Host.AllowDuplicated && g.isDuplicate(computeHash(data)) {
		if g.Host.Log {
			gprint("Insert", "Data already exists in Table. Returning")
		}
		return
	}
	for i, v := range g.Columns {
		switch data[i].(type) {
		case int:
			if v.Datatype != INT {
				panic("Mismatch in Datatype")
			}
		case string:
			if v.Datatype != STRING {
				panic("Mismatch in Datatype")
			}
		case bool:
			if v.Datatype != BOOLEAN {
				panic("Mismatch in Datatype")
			}
		case float64:
			if v.Datatype != FLOAT {
				panic("Mismatch in Datatype")
			}
		case float32:
			if v.Datatype != FLOAT {
				panic("Mismatch in Datatype")
			}
		default:
			if data[i] != nil {
				panic("Unknown datatype")
			}
		}
	}
	g.Rows = append(g.Rows, data)
}
