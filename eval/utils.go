package eval

import (
	"fmt"
	"strconv"
	"time"
)

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
	t2, err := time.Parse("2006-01-02", d2)
	if err != nil {
		panic("Error parsing dates")
	}
	td2 := t2.Unix()
	td1 := t1.Unix()
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

func checkIfCompatible(r, l *token) {
	if !(r.tokenType == tokenTypeInt && l.tokenType == tokenTypeInt ||
		l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeFloat ||
		l.tokenType == tokenTypeInt && r.tokenType == tokenTypeFloat ||
		l.tokenType == tokenTypeFloat && r.tokenType == tokenTypeInt ||
		l.tokenType == tokenTypeDate && r.tokenType == tokenTypeDate ||
		l.tokenType == tokenTypeChar && r.tokenType == tokenTypeInt ||
		l.tokenType == tokenTypeInt && r.tokenType == tokenTypeChar ||
		l.tokenType == tokenTypeString && r.tokenType == tokenTypeString ||
		l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
		panic(fmt.Sprintf("cant do operation for tokens %s , %s", l.value, r.value))
	}
}
