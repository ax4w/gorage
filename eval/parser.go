package eval

import (
	"strconv"
	"strings"
)

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
		//handle braces
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
