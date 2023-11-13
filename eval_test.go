package Gorage

import (
	"testing"
)

/*
			&
		=		=
	a	  2   b		3
*/

// (a & b) & (c & d)
// a = true & b != true & c = true

// t
//

func TestEval(t *testing.T) {
	//println(float64(int(rune('a'))))
	//f := fmt.Sprintf("'%s' = 'hallo welt' & %d = 3 & %d = 4", "hallo welt", 3, 4)

	p := parse("( 'William' == 'William' && 2 == 2 ) || 85.5 >= 90.0")
	tr := toTree(p)
	//traverseTree(tr[0])
	//println(len(tr))
	//println(len(tr))
	println(string(eval(tr[0]).value))
	//traverseTree(tr[0])

	//_printTree(&tr[0])
	//_printTree(&tr[0])
}
