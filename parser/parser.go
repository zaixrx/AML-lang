// TODO: support non-cascaded error handling and syntax error recovery
// in case of entering panic mode (to resume parsing other tokens)

// TODO: Add error productions to handle each binary operator appearing without a left-hand operand. 
// aka: detect a binary operator appearing at the beginning of an expression.
// Report that as an error, but also parse and discard a right-hand operand with the appropriate precedence.
package parser

import (
	"fmt"
	"aml/lexer"
)

type Value = any;

type Parser struct {
	current int;
	tokens []*lexer.Token;
};

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser {
		tokens: tokens,
		current: 0,
	};
}

func (p *Parser) eof(offset uint) bool {
	return p.current + int(offset) >= len(p.tokens);
}

func (p *Parser) prev() *lexer.Token {
	return p.tokens[p.current - 1];
}

func (p *Parser) expect(tts ...lexer.TokenType) bool {
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

func (p *Parser) declarestmt() (Stmt, error) {
	if p.expect(lexer.VAR) {
		if !p.expect(lexer.IDENTIFIER) {
			return nil, p.generate_expect_error("IDENTIFIER after 'var' in variable declartion");
		}
		id := p.prev();
		if !p.expect(lexer.EQUAL) {
			return nil, p.generate_expect_error("'=' after IDENTIFIER in variable declartion");
		}
		target, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if !p.expect(lexer.SEMICOLON) {
			return nil, p.generate_expect_error("';' at the end of the statement.");
		}
		return &VarDeclarationStmt{
			Name: string(id.Lexeme),
			Asset: target,
		}, nil;
	}
	return p.printstmt();
}

func (p *Parser) printstmt() (Stmt, error) {
	if p.expect(lexer.PRINT) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if !p.expect(lexer.SEMICOLON) {
			return nil, p.generate_expect_error("';' at the end of the statement.");
		}
		return &PrintStmt{
			Asset: expr,
		}, nil;
	}
	return p.statement();
}

func (p *Parser) statement() (Stmt, error) {
	expr, err := p.expression();
	if err != nil {
		return nil, err;
	}
	if !p.expect(lexer.SEMICOLON) {
		return nil, p.generate_expect_error("';' at the end of the statement.");
	}
	return &ExprStmt{
		InnerExpr: expr,
	}, nil;
}

// expression -> equality
func (p *Parser) expression() (Expr, error) {
	return p.expressions();
}

// expressions -> assign "," assign | expressions* 
func (p *Parser) expressions() (Expr, error) {
	expr, err := p.assign();
	if err != nil {
		return nil, err;
	}
	if p.expect(lexer.COMMA) {
		operator := p.prev();
		right, err := p.assign();
		if err != nil {
			return nil, err;
		}
		return &BinaryExpr{
			LOperand: expr,
			Operator: operator,
			ROperand: right,
		}, nil;
	}
	return expr, nil;
}

func (p *Parser) assign() (Expr, error) {
	if p.expect(lexer.IDENTIFIER) {
		id := p.prev();
		if !p.expect(lexer.EQUAL) {
			return nil, p.generate_expect_error("'=' after IDENTIFIER in assignment");
		}
		src, err := p.ternary();
		if err != nil {
			return nil, err;
		}
		return &AssignExpr{
			To: string(id.Lexeme),
			From: src,
		}, nil;
	}
	return p.ternary();
}

// ternay -> equality "?" equality ":" ternary*;
// example: (a > b) ? a : b > a ? b : 0;
func (p *Parser) ternary() (Expr, error) {
	expr, err := p.equality();
	if err != nil {
		return nil, err;
	}
	if p.expect(lexer.QUESTION) {
		iftrue, err := p.equality();
		if err != nil {
			return nil, err;
		}
		if !p.expect(lexer.COLON) {
			return nil, p.generate_expect_error("':' in ternay operator");
		}
		iffalse, err := p.ternary();
		if err != nil {
			return nil, err;
		}
		return &TernaryExpr{
			Cond: expr,
			Iftrue: iftrue,
			Iffalse: iffalse,
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
	for p.expect(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.prev()
		right, err := p.comparison();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			LOperand: expr,
			Operator: operator,
			ROperand: right,
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
	for p.expect(lexer.LESS, lexer.GREATER, lexer.LESS_EQUAL, lexer.GREATER_EQUAL) {
		operator := p.prev();
		right, err := p.term();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			LOperand: expr,
			Operator: operator,
			ROperand: right,
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
	for p.expect(lexer.PLUS, lexer.MINUS) {
		operator := p.prev();
		right, err := p.term();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			LOperand: expr,
			Operator: operator,
			ROperand: right,
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
	for p.expect(lexer.STAR, lexer.SLASH) {
		operator := p.prev();
		right, err := p.term();
		if err != nil {
			return nil, err;
		}
		expr = &BinaryExpr {
			LOperand: expr,
			Operator: operator,
			ROperand: right,
		};
	}
	return expr, nil;
}

// unary   -> primary | (("!" | "-") unary)*
func (p *Parser) unary() (Expr, error) {
	if p.expect(lexer.BANG, lexer.MINUS) {
		operand, err := p.unary();
		if err != nil {
			return nil, err;
		}
		operator := p.prev();
		return &UnaryExpr {
			Operand: operand,
			Operator: operator,
		}, nil;
	}
	return p.primary();
}

// primary -> IDENTIFIER | STRING | NUMBER | "true" | "false" | "null" | "(" expression ")"
func (p *Parser) primary() (Expr, error) {
	if p.expect(lexer.TRUE) {
		return &LiteralExpr{
			ValueLiteral: true,
		}, nil
	} else if p.expect(lexer.FALSE) {
		return &LiteralExpr{
			ValueLiteral: false,
		}, nil
	} else if p.expect(lexer.NULL) {
		return &LiteralExpr{
			ValueLiteral: nil,
		}, nil
	} else if p.expect(lexer.IDENTIFIER, lexer.STRING, lexer.NUMBER) {
		return &LiteralExpr {
			ValueLiteral: p.prev().Literal,
		}, nil
	} else if p.expect(lexer.LEFT_PAREN) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if p.expect(lexer.RIGHT_PAREN) {
			return &GroupingExpr{
				InnerExpr: expr,
			}, nil;
		}
		return nil, p.generate_expect_error(")");
	}
	return nil, p.generate_expect_error("valid token");
}

func (p *Parser) generate_expect_error(expected string) error {
	return fmt.Errorf("Parser Error: expected %s\n", expected);
}

func (p *Parser) Parse() ([]Stmt, error) {
	stmts := make([]Stmt, 0);
	for !p.eof(0) {
		stmt, err := p.declarestmt();
		if err != nil {
			return nil, err;
		}
		stmts = append(stmts, stmt);
	}
	return stmts, nil;
}
