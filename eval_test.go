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
	if runEval("( 'William' == 'William' && 2 == 2 ) || 85.5 >= 90.0") != "t" {
		t.Fatalf("Should return true")
	}
	if runEval("1 != 1") != "f" {
		t.Fatalf("Should be false")
	}
	if runEval("( 'Hi' == 'hi' ) || ( 1 == 1 && ( 5 != 5 !& ( t == f ) ) ) && 1.0 < 1.1") != "t" {
		t.Fatalf("Should be true")
	}
	//traverseTree(tr[0])

	//_printTree(&tr[0])
	//_printTree(&tr[0])
}
