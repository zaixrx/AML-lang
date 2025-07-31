package main

import (
	"fmt"
	"os"
)

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
	FUN
	FOR
	IF
	NIL
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

func (tt TokenType) ToString() string {
	return "TODO";
}

type Token struct {
	Type TokenType
	Lexeme []rune 
	Literal any
	Line uint
};

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type.ToString(), string(t.Lexeme), t.Literal);
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

const EOF_RUNE = '\000';

// atomic
func (s *Scanner) eof() bool {
	return s.current >= uint(len(s.source));
}

// atomic
func (s *Scanner) generate_error() error {
	return fmt.Errorf("%s:%d:%d unexpected token", s.filename, s.line, s.current);
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

func (s *Scanner) consume_string() error {
	literal := make([]rune, 0);
	var r rune;
	for r = s.consume_rune(); !(r == '"' || r == EOF_RUNE); r = s.consume_rune() {
		literal = append(literal, r);
	}
	if (r == EOF_RUNE) {
		return s.generate_error();
	}
	s.add_token_literal(STRING, literal);
	return nil
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
			if err := s.consume_string(); err != nil {
				return err;
			}
			break;
		}
		default: {
			return s.generate_error();
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

func main() {
	if (len(os.Args) != 2) {
		fmt.Printf("usage: %s <file.amel>\n", os.Args[0]);
		return;
	}
	filename := os.Args[1];
	fbytes, err := os.ReadFile(filename);
	if err != nil {
		fmt.Println("could not read file", filename);
		fmt.Println(err);
		return;
	}
	scanner := NewScanner(filename, string(fbytes))
	tokens, err := scanner.Scan();
	if err != nil {
		fmt.Println(err);
		return;
	}
	for _, token := range tokens {
		fmt.Printf("{\n    lexme: %s,\n}\n", string(token.Lexeme));
	}
}
