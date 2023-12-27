package eval

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

func evaluate(f *token) *token {
	if f.left == nil && f.right == nil {
		return f
	}
	l := evaluate(f.left)
	m := f
	r := evaluate(f.right)
	switch string(m.value) {
	case "<":
		return less(l, r)
	case ">":
		return greater(l, r)
	case ">=":
		return greaterThan(l, r)
	case "<=":
		return lessThan(l, r)
	case "==":
		return equal(l, r)
	case "!=":
		return notEqual(l, r)
	case "&&":
		return and(l, r)
	case "||":
		return or(l, r)
	case "!&":
		return nand(l, r)
	case "!|":
		return nor(l, r)
	}
	panic("UNREACHABLE: NOT A VALID OPERATOR " + string(m.value))
	return nil
}

func Eval(f string) bool {
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
