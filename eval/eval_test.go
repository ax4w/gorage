package eval

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
	if !Eval("( 'William' == 'William' && 2 == 2 ) || 85.5 >= 90.0") {
		t.Fatalf("Should return true")
	}
	if Eval("1 != 1") {
		t.Fatalf("Should be false")
	}
	if !Eval("( 'Hi' == 'hi' ) || ( 1 == 1 && ( 5 != 5 !& ( t == f ) ) ) && 1.0 < 1.1") {
		t.Fatalf("Should be true")
	}
	if !Eval("2023-11-19 == 2023-11-19") {
		t.Fatalf("Should be true")
	}
	if !Eval("2022-11-19 <= 2023-11-19") {
		t.Fatalf("Should be true")
	}
}
