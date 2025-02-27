package namedflags

import (
	"reflect"
	"testing"
)

type goodNF struct {
	Frodo,
	Sam bool
}

type badTypedNF struct {
	Frodo bool
	Sam   string
}

type badTooLargeNF struct {
	A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z,
	Aa, Bb, Cc, Dd, Ee, Ff, Gg, Hh, Ii, Jj, Kk, Ll, Mm, Nn, Oo, Pp, Qq, Rr, Ss, Tt, Uu, Vv, Ww, Xx, Yy, Zz,
	Aaa, Bbb, Ccc, Ddd, Eee, Fff, Ggg, Hhh, Iii, Jjj, Kkk, Lll, Mmm, Nnn, Ooo, Ppp, Qqq, Rrr, Sss, Ttt, Uuu, Vvv, Www, Xxx, Yyy,
	Zzz bool
}

func TestIsValid(t *testing.T) {
	type testCase struct {
		name       string
		in         reflect.Value
		want_error bool
	}
	testCases := [...]testCase{
		{
			"valid struct",
			reflect.ValueOf(&goodNF{}).Elem(),
			false,
		},
		{
			"not all fields on struct are bool",
			reflect.ValueOf(&badTypedNF{}).Elem(),
			true,
		},
		{
			"too many fields on struct",
			reflect.ValueOf(&badTooLargeNF{}).Elem(),
			true,
		},
		{
			"bool instead of struct",
			reflect.ValueOf(new(bool)).Elem(),
			true,
		},
		{
			"int32 instead of struct",
			reflect.ValueOf(new(int32)).Elem(),
			true,
		},
		{
			"string instead of struct",
			reflect.ValueOf(new(string)).Elem(),
			true,
		},
	}
	for _, tc := range testCases {
		if err := isValid(tc.in); (err == nil) && tc.want_error {
			t.Errorf(`case "%s": expected error but didn't get one`, tc.name)
		} else if (err != nil) && !(tc.want_error) {
			t.Errorf(`case "%s": got unexpected error "%v"`, tc.name, err)
		}
	}
}

func TestFromInt(t *testing.T) {
	type goodTestCase struct {
		in   uint
		want goodNF
	}
	goodTestCases := [...]goodTestCase{
		// single flag set
		{1, goodNF{Frodo: true}},
		// multiple flags set
		{3, goodNF{Frodo: true, Sam: true}},
	}
	for _, tc := range goodTestCases {
		got, err := FromInt[goodNF](tc.in)
		if err != nil {
			t.Errorf("Got unexpected error %v", err)
		}
		if got != tc.want {
			t.Errorf("got %v, wanted %v", got, tc.want)
		}
	}
}

func TestToInt(t *testing.T) {
	type goodTestCase struct {
		in   goodNF
		want uint
	}
	goodTestCases := [...]goodTestCase{
		// single flag set
		{goodNF{Frodo: true}, 1},
		// multiple flags set
		{goodNF{Frodo: true, Sam: true}, 3},
	}
	for _, tc := range goodTestCases {
		res, err := ToInt(tc.in)
		if err != nil {
			t.Errorf("Got unexpected error %v", err)
		}
		if res != tc.want {
			t.Errorf("got %v, wanted %v", res, tc.want)
		}
	}
}
