package main

import "fmt";

type Bool bool
func (b Bool) String() string { return fmt.Sprintf("%t", b) }

type Number float64
func (n Number) String() string { return fmt.Sprintf("%g", n) }

type Str string
func (s Str) String() string { return string(s) }

type Null struct{}
func (n Null) String() string { return "null" }

func IsAlpha(val rune) bool {
	return 'a' <= val && val <= 'z' || 
	       'A' <= val && val <= 'Z' || val == '_';
}

func IsNum(val rune) bool {
	return '0' <= val && val <= '9';
}

func IsAlphaNum(val rune) bool {
	return IsAlpha(val) || IsNum(val);
}

func Normalize(val int) float64 {
	var fval = float64(val);
	for fval > 1 {
		fval /= 10;
	}
	return fval;
}
