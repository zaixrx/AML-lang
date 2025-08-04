package lexer

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
