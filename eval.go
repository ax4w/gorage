package Gorage

import (
	"strconv"
	"strings"
	"time"
)

const (
	tokenTypeInt     = 1
	tokenTypeString  = 2
	tokenTypeFloat   = 3
	tokenTypeBoolean = 4
	tokenTypeChar    = 5
	tokenTypeDate    = 6
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

/*
* 	HELPER FUNCTIONS
 */

func validateDate(d string) bool {
	_, err := time.Parse("2006-01-02", d)
	if err != nil {
		return false
	}
	return true
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

// -1 d1 is greater
// 0 equal
// 1 d2 is greater
func compareDates(d1, d2 string) int {
	t1, err := time.Parse("2006-01-02", d1)
	if err != nil {
		panic("Error parsing dates")
	}
	td1 := t1.Unix()
	t2, err := time.Parse("2006-01-02", d2)
	if err != nil {
		panic("Error parsing dates")
	}
	td2 := t2.Unix()
	switch {
	case td1 > td2:
		return -1
	case td1 == td2:
		return 0
	case td2 > td1:
		return 1
	}
	return 0
}

func compareCheck(r, l *token) {
	if !(r.tokenType == tokenTypeInt && l.tokenType == tokenTypeInt ||
		l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeFloat ||
		l.tokenType == tokenTypeInt && r.tokenType == tokenTypeFloat ||
		l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeInt ||
		l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate) {
		panic("< is only supported for int, float and date")
	}
}

/*£
*	-----------------------------------------
 */

func evaluate(f *token) *token {
	if f.left == nil && f.right == nil {
		return f
	}
	l := evaluate(f.left)
	m := f
	r := evaluate(f.right)
	switch string(m.value) {
	case "<":
		compareCheck(r, l)
		if l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate {
			i := compareDates(string(l.value), string(r.value))
			if i == 1 {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		} else {
			lv := convertBytesToFloat(l.value)
			rv := convertBytesToFloat(r.value)
			if lv < rv {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		}
	case ">":
		compareCheck(r, l)
		if l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate {
			i := compareDates(string(l.value), string(r.value))
			if i == -1 {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		} else {
			lv := convertBytesToFloat(l.value)
			rv := convertBytesToFloat(r.value)
			if lv > rv {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		}
	case ">=":
		compareCheck(r, l)
		if l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate {
			i := compareDates(string(l.value), string(r.value))
			if i == -1 || i == 0 {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		} else {
			lv := convertBytesToFloat(l.value)
			rv := convertBytesToFloat(r.value)
			if lv >= rv {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		}
	case "<=":
		compareCheck(r, l)
		if l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate {
			i := compareDates(string(l.value), string(r.value))
			if i == 1 || i == 0 {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		} else {
			lv := convertBytesToFloat(l.value)
			rv := convertBytesToFloat(r.value)
			if lv <= rv {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		}
	case "==":
		if !(l.tokenType == r.tokenType ||
			l.tokenType == tokenTypeChar && r.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeChar) {
			panic("mismatching == types")
		}
		if l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate {
			i := compareDates(string(l.value), string(r.value))
			if i == 0 {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		} else {
			if compareByteArray(l.value, r.value) {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		}
	case "!=":
		if !(l.tokenType == r.tokenType ||
			l.tokenType == tokenTypeChar && r.tokenType == tokenTypeInt ||
			l.tokenType == tokenTypeInt && r.tokenType == tokenTypeChar) {
			panic("mismatching == types")
		}
		if l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate {
			i := compareDates(string(l.value), string(r.value))
			if i != 0 {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		} else {
			if !compareByteArray(l.value, r.value) {
				return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
			}
			return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
		}
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
				if validateDate(k) {
					tokenType = tokenTypeDate
				} else if _, err := strconv.Atoi(k); err == nil {
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

func eval(f string) bool {
	p := parse(f)
	t := toTree(p)
	if len(t) == 0 {
		panic("Error while eval")
	}
	e := evaluate(t[0])
	if e == nil {
		panic("Eval returned nil")
	}
	if string(e.value) == "t" {
		return true
	}
	return false
}
