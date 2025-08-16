package lexer

import "fmt";

const EOF_RUNE = '\000';

type Token struct {
	Type TokenType
	Lexeme string
	Literal any
	Line uint
};

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type.ToString(), string(t.Lexeme), t.Literal);
}

type TokenType uint;
const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	QUESTION
	COLON

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	TRUE
	FALSE
	NULL
	AND
	OR
	IF
	ELSE
	WHILE
	FOR
	BREAK
	CONTINUE
	VAR
	FUNC	
	RETURN
	CLASS
	THIS
	SUPER
	PRINT

	EOF
);

// subject to change
var keywords_map = map[string]TokenType {
	"and": AND,
	"class": CLASS,
	"else": ELSE,
	"func": FUNC,
	"for": FOR,
	"break": BREAK,
	"continue": CONTINUE,
	"if": IF,
	"null": NULL,
	"or": OR,
	"return": RETURN,
	"super": SUPER,
	"print": PRINT,
	"this": THIS,
	"true": TRUE,
	"false": FALSE,
	"var": VAR,
	"while": WHILE,
};

func (tt TokenType) ToString() string {
	switch tt {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case QUESTION:
		return "QUESTION"
	case COLON:
		return "COLON"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FUNC:
		return "FUNC"
	case FOR:
		return "FOR"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case IF:
		return "IF"
	case NULL:
		return "NULL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
