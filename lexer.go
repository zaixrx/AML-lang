package main

import (
	"fmt"
	"strconv"
)

const EOF_RUNE = '\000';

type Token struct {
	Type TokenType
	Lexeme []rune 
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
	AND
	CLASS
	ELSE
	FALSE
	PROD	
	FOR
	IF
	NULL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
);

// subject to change
var keywords_map = map[string]TokenType {
	"and": AND,
	"class": CLASS,
	"else": ELSE,
	"prod": PROD,
	"for": FOR,
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
	case PROD:
		return "PROD"
	case FOR:
		return "FOR"
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

type Scanner struct {
	filename string
	source []rune 
	tokens []*Token

	start uint
	current uint
	line uint
};

func NewScanner(filename string, source string) *Scanner {
	return &Scanner {
		filename: filename,
		source: []rune(source),
		tokens: make([]*Token, 0),
		start: 0,
		current: 0,
		line: 1,
	};
}

// atomic
func (s *Scanner) add_token_literal(tt TokenType, literal any) *Token {
	token := &Token{
		Type: tt,
		Lexeme: s.source[s.start : s.current],
		Literal: literal,
		Line: s.line,
	};
	s.tokens = append(s.tokens, token);
	return token;
}

// molecular
func (s *Scanner) add_token(tt TokenType) *Token {
	return s.add_token_literal(tt, nil);
}

// atomic
func (s *Scanner) eof() bool {
	return s.current >= uint(len(s.source));
}

// atomic
func (s *Scanner) generate_error(description string) error {
	return fmt.Errorf("%s:%d:%d unexpected token\ndescription: %s", s.filename, s.line, s.current, description);
}

// atomic: lookahead with one character
func (s *Scanner) peek_rune() rune {
	if s.eof() {
		return EOF_RUNE;
	}
	return s.source[s.current];
}

// atomic
func (s *Scanner) consume_rune() rune {
	if s.eof() {
		return EOF_RUNE;
	}
	r := s.source[s.current];
	s.current++;
	return r;
}

// molecular
func (s *Scanner) expect_rune(r rune) bool {
	if s.peek_rune() == r {
		s.current++;
		return true
	}
	return false
}

func (s *Scanner) consume_string() (string, error) {
	var r rune;
	for r = s.consume_rune(); !(r == '"' || r == EOF_RUNE); r = s.consume_rune() {
		// TODO: handle invalid string characters
		if (r == '\n') {
			s.line++;
		}
	}
	if (r == EOF_RUNE) {
		return "", s.generate_error("unterminated string");
	}
	return string(s.source[s.start:s.current]), nil
}

func (s *Scanner) consume_integer() (int, error) {
	for r := s.peek_rune(); IsNum(r) && r != EOF_RUNE; r = s.peek_rune() {
		fmt.Println(string(r), s.start, s.current);
		s.consume_rune();
	}
	literal, err := strconv.Atoi(string(s.source[s.start:s.current]));
	if err != nil {
		return 0, s.generate_error(err.Error());
	}
	return literal, nil
}

func (s *Scanner) consume_identifier() (string, error) {
	for r := s.peek_rune(); IsAlphaNum(r) && r != EOF_RUNE; r = s.peek_rune() {
		s.consume_rune();
	}
	return string(s.source[s.start:s.current]), nil
}

// cellular
func (s *Scanner) scan_curr() error {
	s.start = s.current;
	char := s.consume_rune();
	switch char {
		case '(': { s.add_token(LEFT_PAREN); break; } 
		case ')': { s.add_token(RIGHT_PAREN); break; }
		case '{': { s.add_token(LEFT_BRACE); break; }
		case '}': { s.add_token(RIGHT_BRACE); break; }
		case ',': { s.add_token(COMMA); break; }
		case '.': { s.add_token(DOT); break; }
		case '-': { s.add_token(MINUS); break; }
		case '+': { s.add_token(PLUS); break; }
		case ';': { s.add_token(SEMICOLON); break; }
		case '*': { s.add_token(STAR); break; }
		case '!': {
			tt := BANG;
			if s.expect_rune('=') {
				tt = BANG_EQUAL;
			}
			s.add_token(tt) 
			break;
		}
		case '=': {
			tt := EQUAL;
			if s.expect_rune('=') {
				tt = EQUAL_EQUAL;
			}
			s.add_token(tt) 
			break;
		}
		case '>': {
			tt := GREATER;
			if s.expect_rune('=') {
				tt = GREATER_EQUAL;
			}
			s.add_token(tt) 
			break;
		}
		case '<': {
			tt := LESS;
			if s.expect_rune('=') {
				tt = LESS_EQUAL;
			}
			s.add_token(tt) 
			break;
		}
		case '/': {
			if s.expect_rune('/') {
				// ignore all of the following text
				for r := s.peek_rune(); !(r == '\n' || s.eof()); r = s.peek_rune() {
					s.consume_rune();
				}
			} else {
				s.add_token(SLASH);
			}
			break;
		}
		case ' ':
		case '\t':
		case '\r': {
			break;
		}
		case '\n': {
			s.line++;
			break;
		}
		case '"': {
			literal ,err := s.consume_string()
			if err != nil {
				return err;
			}
			s.add_token_literal(STRING, literal);
			break;
		}
		default: {
			if IsNum(char) {
				s.current--;
				var literal float64 = 0.0;
				integer, err := s.consume_integer();
				if err != nil {
					return err;
				}
				literal += float64(integer);
				if (s.peek_rune() == '.') {
					s.consume_rune();
					if IsNum(s.peek_rune()) {
						rollback := s.start;
						s.start = s.current;
						integer, err = s.consume_integer();
						s.start = rollback;
						if err != nil {
							return err;
						}
						literal += Normalize(integer);
					} else {
						return s.generate_error("expected number after .");
					}
				}
				s.add_token_literal(NUMBER, literal);
			} else if IsAlpha(char) {
				s.current--;
				literal, err := s.consume_identifier();
				if err != nil {
					return err;
				}
				tt := IDENTIFIER;
				if keyword, pres := keywords_map[literal]; pres {
					tt = keyword;
				}
				s.add_token_literal(tt, literal);
			} else {
				return s.generate_error(fmt.Sprintf("NOTE: learn the fucking language idiot, got %v", char));
			}
		}
	}
	return nil;
}

// organelle
func (s *Scanner) Scan() ([]*Token, error) {
	for !s.eof() {
		if err := s.scan_curr(); err != nil {
			return nil, err
		}
	}
	return s.tokens, nil;
}
