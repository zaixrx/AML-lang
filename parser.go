// TODO: support non-cascaded error handling and syntax error recovery
// in case of entering panic mode (to resume parsing other tokens)

// TODO: Add error productions to handle each binary operator appearing without a left-hand operand. 
// aka: detect a binary operator appearing at the beginning of an expression.
// Report that as an error, but also parse and discard a right-hand operand with the "appropriate precedence".
package main

import (
	"fmt"
	"strings"
)

type Expr interface {
	String() string;
};
type TernaryExpr struct {
	cond Expr;
	truthy Expr;
	falsy Expr;
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
func (tir *TernaryExpr) String() string {
	var b strings.Builder;
	fmt.Fprintf(&b, "[ cond: %s ", tir.cond.String());
	fmt.Fprintf(&b, "truthy: %s ", tir.truthy.String());
	fmt.Fprintf(&b, "falsy: %s ", tir.falsy.String());
	fmt.Fprint(&b, "]");
	return b.String();
}
func (bin *BinaryExpr) String() string {
	var b strings.Builder;
	fmt.Fprintf(&b, "[ %s ", bin.operator.String());
	if bin.left != nil {
		fmt.Fprintf(&b, "left: %s ", bin.left.String());
	}
	if bin.right != nil {
		fmt.Fprintf(&b, "right: %s ", bin.right.String());
	}
	fmt.Fprint(&b, "]");
	return b.String();
}
func (un *UnaryExpr) String() string {
	var b strings.Builder;
	fmt.Fprintf(&b, "[ %s ", un.operator.String());
	if un.operand != nil {
		fmt.Fprintf(&b, "left: %s ", un.operand.String());
	}
	fmt.Fprint(&b, "]");
	return b.String();
}
func (lit *LiteralExpr) String() string {
	return fmt.Sprintf("%v", lit.value);
}

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
	return p.ternary();
}

// ternay -> equality "?" expressions ":" ternary;
// example: (a > 0) ? a, b, c : (b > 0 ? c, b, a : b);
func (p *Parser) ternary() (Expr, error) {
	expr, err := p.expressions();
	if err != nil {
		return nil, err;
	}
	if p.expect(QUESTION) {
		truthy, err := p.expressions();
		if err != nil {
			return nil, err;
		}
		if !p.expect(COLON) {
			return nil, p.generate_expect_error("':' in ternay operator");
		}
		falsy, err := p.ternary();
		if err != nil {
			return nil, err;
		}
		return &TernaryExpr{
			cond: expr,
			truthy: truthy,
			falsy: falsy,
		}, nil;
	}
	return expr, nil;
}

// expressions -> equality "," equality
func (p *Parser) expressions() (Expr, error) {
	expr, err := p.equality();
	if err != nil {
		return nil, err;
	}
	if p.expect(COMMA) {
		operator := p.prev();
		exprs, err := p.equality();
		if err != nil {
			return nil, err;
		}
		return &BinaryExpr{
			left: expr,
			operator: operator,
			right: exprs,
		}, nil;
	}
	return expr, nil;
}

// equality -> comparison (("!=" | "==") comparison)*
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
	} else if p.expect(IDENTIFIER, STRING, NUMBER) {
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
	return fmt.Errorf("Parser Error: expected %s\n", expected);
}

func (p *Parser) Parse() ([]Expr, error) {
	exprs := make([]Expr, 0);
	for !p.eof(0) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if !p.expect(SEMICOLON) {
			return nil, p.generate_expect_error(";");
		}
		exprs = append(exprs, expr);
	}
	return exprs, nil;
}
