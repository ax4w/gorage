package Gorage

import (
	"strings"
)

type token struct {
	value []byte
	left  *token
	right *token
}

var (
	//keywords    = []string{"&&", "!&", "||", "!|", "==", "!="}
	strongSplit = []string{"&&", "!&", "||", "!|"}
)

func compareByteArray(b1, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}
	for i, _ := range b1 {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
}

func eval(f *token) *token {
	if f.left == nil && f.right == nil {
		return f
	}
	l := eval(f.left)
	m := f
	r := eval(f.right)
	switch string(m.value) {
	case "==":
		if compareByteArray(l.value, r.value) {
			return &token{value: []byte("t")}
		}
		return &token{value: []byte("f")}
	case "!=":
		if !compareByteArray(l.value, r.value) {
			return &token{value: []byte("t")}
		}
		return &token{value: []byte("f")}
	case "&&":
		if compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t")) {
			return &token{value: []byte("t")}
		}
		return &token{value: []byte("f")}
	case "||":
		if compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t")) {
			return &token{value: []byte("t")}
		}
		return &token{value: []byte("f")}
	case "!&":
		if !(compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t"))) {
			return &token{value: []byte("t")}
		}
		return &token{value: []byte("f")}
	case "!|":
		if !(compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t"))) {
			return &token{value: []byte("t")}
		}
		return &token{value: []byte("f")}
	}
	println(string(m.value))
	println("UNREACHABLE")
	return nil
}

func toTree(nodes []*token) []*token {
	var op string
	var query []*token
	var base []*token

	buildQuery := func() {
		nq := &token{
			value: []byte(op),
			left:  query[0],
			right: query[1],
		}
		query = []*token{}
		op = ""
		query = append(query, nq)
	}

	for i := 0; i < len(nodes); i++ {
		isOp := false
		for _, k := range strongSplit {
			if compareByteArray([]byte(k), nodes[i].value) {
				op = k
				isOp = true
			}
		}
		if string(nodes[i].value) == "(" && !isOp {
			openCount := 1
			var tmp []*token
			i += 1
			for ; i < len(nodes); i++ {
				if string(nodes[i].value) == "(" {
					openCount++
				}
				if string(nodes[i].value) == ")" {
					openCount--
					if openCount == 0 {
						break
					}

				}
				tmp = append(tmp, nodes[i])
			}
			if i == len(nodes) {
				panic("No ) found")
			}
			query = append(query, toTree(tmp)[0])
			if len(query) == 2 {
				buildQuery()
			}
			continue
		}
		if !isOp {
			base = append(base, nodes[i])
		}
		if len(base) == 3 {
			ne := &token{
				value: base[1].value,
				left:  &token{value: base[0].value},
				right: &token{value: base[2].value},
			}
			query = append(query, ne)
			if len(query) == 2 {
				buildQuery()
			}
			base = []*token{}
		}
	}
	return query
}

func traverseTree(t *token) {
	if t == nil {
		return
	}
	traverseTree(t.left)
	println(string(t.value))
	traverseTree(t.right)
}

func parse(f string) []*token {
	var nodes []*token
	split := strings.SplitAfter(f, "'")
	var r []string
	for _, v := range split {
		if strings.Contains(v, "'") {
			v = strings.ReplaceAll(v, "'", "")
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			r = append(r, v)
			continue
		}
		s := strings.Split(v, " ")
		for _, k := range s {
			k = strings.TrimSpace(k)
			if len(k) == 0 {
				continue
			}
			r = append(r, k)
		}
	}
	for _, v := range r {
		nodes = append(nodes, &token{
			value: []byte(v),
			left:  nil,
			right: nil,
		})
	}
	return nodes
}
