package Gorage

import (
	"strings"
)

type fragment struct {
	value []byte
	left  *fragment
	right *fragment
}

var (
	keywords    = []string{"&", "!&", "|", "!|", "=", "!="}
	strongSplit = []string{"&", "!&", "|", "!|"}
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

func eval(f *fragment) *fragment {
	if f.left == nil && f.right == nil {
		return f
	}
	l := eval(f.left)
	m := f
	r := eval(f.right)
	switch string(m.value) {
	case "=":
		if compareByteArray(l.value, r.value) {
			return &fragment{value: []byte("t")}
		}
		return &fragment{value: []byte("f")}
	case "!=":
		if !compareByteArray(l.value, r.value) {
			return &fragment{value: []byte("t")}
		}
		return &fragment{value: []byte("f")}
	case "&":
		if compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t")) {
			return &fragment{value: []byte("t")}
		}
		return &fragment{value: []byte("f")}
	case "|":
		if compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t")) {
			return &fragment{value: []byte("t")}
		}
		return &fragment{value: []byte("f")}
	case "!&":
		if !(compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t"))) {
			return &fragment{value: []byte("t")}
		}
		return &fragment{value: []byte("f")}
	case "!|":
		if !(compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t"))) {
			return &fragment{value: []byte("t")}
		}
		return &fragment{value: []byte("f")}
	}
	println(string(m.value))
	println("UNREACHABLE")
	return nil
}

func toTree(nodes []*fragment) []*fragment {
	var op string
	var query []*fragment
	var base []*fragment
	for i := 0; i < len(nodes); i++ {
		isOp := false
		for _, k := range strongSplit {
			if compareByteArray([]byte(k), nodes[i].value) { //k == v.value {
				op = k
				isOp = true
			}
		}
		if string(nodes[i].value) == "(" && !isOp {
			openCount := 1
			var tmp []*fragment
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
			sub := toTree(tmp)
			//traverseTree(sub[0])
			query = append(query, sub[0])
			if len(query) == 2 {
				nq := &fragment{
					value: []byte(op),
					left:  query[0],
					right: query[1],
				}
				query = []*fragment{}
				op = ""
				query = append(query, nq)
			}
			continue
		}
		if !isOp {
			base = append(base, nodes[i])
		}
		if len(base) == 3 {
			ne := &fragment{
				value: base[1].value,
				left:  &fragment{value: base[0].value},
				right: &fragment{value: base[2].value},
			}
			query = append(query, ne)
			if len(query) == 2 {
				nq := &fragment{
					value: []byte(op),
					left:  query[0],
					right: query[1],
				}
				query = []*fragment{}
				op = ""
				query = append(query, nq)
			}
			base = []*fragment{}
		}
	}
	return query
}

func traverseTree(t *fragment) {
	if t == nil {
		return
	}
	traverseTree(t.left)
	println(string(t.value))
	traverseTree(t.right)
}

func parse(f string) []*fragment {
	var nodes []*fragment
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
		nodes = append(nodes, &fragment{
			value: []byte(v),
			left:  nil,
			right: nil,
		})
	}
	return nodes
}
