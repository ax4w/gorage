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
	if eval("( 'William' == 'William' && 2 == 2 ) || 85.5 >= 90.0") != "t" {
		t.Fatalf("Should return true")
	}
	if eval("1 != 1") != "f" {
		t.Fatalf("Should be false")
	}
	if eval("( 'Hi' == 'hi' ) || ( 1 == 1 && ( 5 != 5 !& ( t == f ) ) ) && 1.0 < 1.1") != "t" {
		t.Fatalf("Should be true")
	}
	if eval("2023-11-19 == 2023-11-19") != "t" {
		t.Fatalf("Should be true")
	}
	if eval("2023-11-19 == 2023-11-19") != "t" {
		t.Fatalf("Should be true")
	}
}
