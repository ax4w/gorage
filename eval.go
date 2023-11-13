package Gorage

import (
	"strconv"
	"strings"
)

const (
	tokenTypeInt     = 1
	tokenTypeString  = 2
	tokenTypeFloat   = 3
	tokenTypeBoolean = 4
	tokenTypeChar    = 5
)

type token struct {
	value     []byte
	left      *token
	right     *token
	tokenType int
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

func convertBytesToFloat(v []byte) float64 {
	s := string(v)
	r, err := strconv.ParseFloat(s, 64)
	if err != nil {
		//check if
		if len(s) == 1 { //prob. a char
			return float64(int(rune(s[0]))) //formatted like +9.00..e+001 - not good for comparison
		}
		panic("Value used in >=,<=,<,> is not a number")
	}
	return r
}

func eval(f *token) *token {
	if f.left == nil && f.right == nil {
		return f
	}
	l := eval(f.left)
	m := f
	r := eval(f.right)
	switch string(m.value) {
	case "<":
		if !(r.tokenType == tokenTypeInt && l.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeInt) {
			panic("<= is only supported for int and float")
		}
		lv := convertBytesToFloat(l.value)
		rv := convertBytesToFloat(r.value)
		if lv < rv {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case ">":
		if !(r.tokenType == tokenTypeInt && l.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeInt) {
			panic("<= is only supported for int and float")
		}
		lv := convertBytesToFloat(l.value)
		rv := convertBytesToFloat(r.value)
		if lv > rv {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case ">=":
		if !(r.tokenType == tokenTypeInt && l.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeInt) {
			panic("<= is only supported for int and float")
		}
		lv := convertBytesToFloat(l.value)
		rv := convertBytesToFloat(r.value)
		if lv >= rv {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case "<=":
		if !(r.tokenType == tokenTypeInt && l.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeFloat ||
			l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeInt) {
			panic("<= is only supported for int and float")
		}
		lv := convertBytesToFloat(l.value)
		rv := convertBytesToFloat(r.value)
		if lv <= rv {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case "==":
		if !(l.tokenType == r.tokenType ||
			l.tokenType == tokenTypeChar && r.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeChar) {
			panic("mismatching == types")
		}
		if compareByteArray(l.value, r.value) {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case "!=":
		if !(l.tokenType == r.tokenType ||
			l.tokenType == tokenTypeChar && r.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeChar) {
			panic("mismatching == types")
		}
		if !compareByteArray(l.value, r.value) {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case "&&":
		if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
			panic("&& expects both sides to be a boolean")
		}
		if compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t")) {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case "||":
		if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
			panic("&& expects both sides to be a boolean")
		}
		if compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t")) {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case "!&":
		if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
			panic("&& expects both sides to be a boolean")
		}
		if !(compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t"))) {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	case "!|":
		if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
			panic("&& expects both sides to be a boolean")
		}
		if !(compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t"))) {
			return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
		}
		return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
	}
	panic("UNREACHABLE: NOT A VALID OPERATOR " + string(m.value))
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
				//println(string(nodes[i].value))
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
				left:  &token{value: base[0].value, tokenType: base[0].tokenType},
				right: &token{value: base[2].value, tokenType: base[2].tokenType},
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

func splitForStrings(f string) (r []string) {
	var tmp string
	inString := false
	for _, v := range f {
		if string(v) != "'" {
			tmp += string(v)
		} else {
			if inString {
				tmp += string(v)
				inString = false
				r = append(r, tmp)
				tmp = ""
				continue
			}
			inString = true
			r = append(r, tmp)
			tmp = ""
			tmp += string(v)

		}
	}
	if len(tmp) > 0 {
		r = append(r, tmp)
	}
	return r
}

func parse(f string) []*token {
	var nodes []*token
	split := splitForStrings(f)
	if len(split) == 0 {
		split = strings.Split(f, " ") //no string in f. We can split normal
	}
	for _, v := range split {
		if strings.Contains(v, "'") {
			tokenType := tokenTypeString
			v = strings.ReplaceAll(v, "'", "")
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			//strings and chars both use ' '. The length decides if it's a char or a string
			if len(v) == 1 {
				tokenType = tokenTypeChar
			}
			nodes = append(nodes, &token{
				value:     []byte(v),
				left:      nil,
				right:     nil,
				tokenType: tokenType,
			})
			continue
		}
		s := strings.Split(v, " ")
		for _, k := range s {
			var tokenType int
			k = strings.TrimSpace(k)
			if len(k) == 0 {
				continue
			}
			switch k {
			case "t", "f":
				tokenType = tokenTypeBoolean
				break
			default:
				if _, err := strconv.Atoi(k); err == nil {
					tokenType = tokenTypeInt
				} else if _, err = strconv.ParseFloat(k, 64); err == nil {
					tokenType = tokenTypeFloat
				}
				break
			}

			nodes = append(nodes, &token{
				value:     []byte(k),
				left:      nil,
				right:     nil,
				tokenType: tokenType,
			})
		}
	}
	return nodes
}
