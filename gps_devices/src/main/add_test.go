package main

import (
	"testing"
)

type Tests struct {
	a, b, res int
}

var tests = []Tests{
	{a: 1, b: 2, res: 3},
	{a: 10, b: 20, res: 30},
	{a: 5, b: 2, res: 7},
	{a: 3, b: 6, res: 10},
}

func Test_add_sum(t *testing.T) {
	for _, test := range tests {
		if Add_sum(test.a, test.b) != test.res {
			t.Errorf("%d , %d not euqal to %d", test.a, test.b, test.res)
		}
	}
}
