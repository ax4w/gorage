package Gorage

import (
	"fmt"
	"strings"
)

type GorageTable struct {
	Name    string
	Columns []string
	Rows    [][]interface{}
}

//i:name
//s:name

//s:name = 'Tom' | i:alter = 10
func (g *GorageTable) Where(f string) *GorageTable {
	res := &GorageTable{
		Name:    "",
		Columns: g.Columns,
		Rows:    [][]interface{}{},
	}
	split := strings.Split(f, " ")
	m := make(map[string]int)
	for i, v := range g.Columns {
		m[v] = i
	}
	for _, v := range g.Rows {
		var tmp []string
		for _, k := range split {
			if strings.Contains(k, "s:") {
				y := strings.Split(k, ":")
				colIdx := m[y[1]]
				k = fmt.Sprintf("'%s'", v[colIdx])
			} else if strings.Contains(k, "i:") || strings.Contains(k, "b:") {
				y := strings.Split(k, ":")
				colIdx := m[y[1]]
				k = fmt.Sprintf("%s", v[colIdx])
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

func (g *GorageTable) Select(columns []string) *GorageTable {
	var columnIdx []int
	tmp := &GorageTable{
		Name:    "",
		Columns: columns,
		Rows:    [][]interface{}{},
	}
	for i, v := range g.Columns {
		for _, k := range columns {
			if v == k {
				columnIdx = append(columnIdx, i)
			}
		}
	}
	for _, v := range g.Rows {
		var t []interface{}
		for _, i := range columnIdx {
			if i >= len(v) {
				continue
			}
			t = append(t, v[i])
		}
		tmp.Rows = append(tmp.Rows, t)
	}
	return tmp
}

func (g *GorageTable) Insert(data []interface{}) {
	if len(data) > len(g.Columns) {
		panic(fmt.Errorf("data has more columns than the table"))
	}
	g.Rows = append(g.Rows, data)
}
