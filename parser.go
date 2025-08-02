package main

import "fmt"

type Expr interface {
	implement();
};
type BinaryExpr struct {
	left  Expr
	operator *Token
	right Expr	
};
type UnaryExpr struct {
	operand Expr
	operator *Token
};
type LiteralExpr struct {
	value any;
}
func (bin *BinaryExpr) implement() {}
func (un *UnaryExpr) implement() {}
func (lit *LiteralExpr) implement() {}

type Parser struct {
	current int;
	tokens []*Token;
};

func NewParser(tokens []*Token) *Parser {
	return &Parser {
		tokens: tokens,
		current: 0,
	};
}

func (p *Parser) eof(offset uint) bool {
	return p.current + int(offset) >= len(p.tokens);
}

func (p *Parser) prev() *Token {
	return p.tokens[p.current - 1];
}

func (p *Parser) expect(tts ...TokenType) bool {
	if p.eof(0) {
		return false;
	}
	for _, tt := range tts {
		if token := p.tokens[p.current]; token.Type == tt {
			p.current++;
			return true;
		}
	}
	return false;
}

// expression -> equality
func (p *Parser) expression() (Expr, error) {
	return p.equality();
}

// equality   -> comparison (("!=" | "==") comparison)*
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison();
	if err != nil {
		return nil, err;
	}
	for p.expect(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.prev()
		right, err := p.comparison();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			left: expr,
			operator: operator,
			right: right,
		};
	}
	return expr, nil;
}

// comparison -> term (("<" | ">" | "<=" | ">=") term)*
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term();
	if err != nil {
		return nil, err;
	}
	for p.expect(LESS, GREATER, LESS_EQUAL, GREATER_EQUAL) {
		operator := p.prev();
		right, err := p.term();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			left: expr,
			operator: operator,
			right: right,
		};
	}
	return expr, nil;
}

// term -> factor | (("+" | "-") factor)*
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor();
	if err != nil {
		return nil, err;
	}
	for p.expect(PLUS, MINUS) {
		operator := p.prev();
		right, err := p.term();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			left: expr,
			operator: operator,
			right: right,
		};
	}
	return expr, nil;
}

// factor -> unary | (("*" | "/") unary)*
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary();
	if err != nil {
		return nil, err;
	}
	for p.expect(STAR, SLASH) {
		operator := p.prev();
		right, err := p.term();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			left: expr,
			operator: operator,
			right: right,
		};
	}
	return expr, nil;
}

// unary   -> primary | (("!" | "-") unary)*
func (p *Parser) unary() (Expr, error) {
	if p.expect(BANG, MINUS) {
		operand, err := p.unary();
		if err != nil {
			return nil, err;
		}
		operator := p.prev();
		return &UnaryExpr {
			operand: operand,
			operator: operator,
		}, nil;
	}
	return p.primary();
}

// primary -> IDENTIFIER | STRING | NUMBER | "true" | "false" | "null" | "(" expression ")"
func (p *Parser) primary() (Expr, error) {
	if p.expect(TRUE) {
		return &LiteralExpr{
			value: true,
		}, nil
	} else if p.expect(FALSE) {
		return &LiteralExpr{
			value: false,
		}, nil
	} else if p.expect(NULL) {
		return &LiteralExpr{
			value: nil,
		}, nil
	} else if p.expect(STRING, NUMBER) {
		return &LiteralExpr {
			value: p.prev().Lexeme,
		}, nil
	} else if p.expect(LEFT_PAREN) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if p.expect(RIGHT_PAREN) {
			return expr, nil;
		}
		return nil, p.generate_expect_error(")");
	}
	return nil, p.generate_expect_error("valid token");
}

func (p *Parser) generate_expect_error(expected string) error {
	return fmt.Errorf("Parser Error: expected %s got %s\n", expected, string(p.tokens[p.current].Lexeme));
}

func (p *Parser) Parse() ([]Expr, error) {
	exprs := make([]Expr, 0);
	for !p.eof(0) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		exprs = append(exprs, expr);
	}
	return exprs, nil;
}
