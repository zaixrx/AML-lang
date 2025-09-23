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
	tokens []lexer.Token;
};

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser {
		tokens: tokens,
		current: 0,
	};
}

func (p *Parser) eof(offset uint) bool {
	return p.current + int(offset) >= len(p.tokens);
}

func (p *Parser) prev() lexer.Token {
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

func (p *Parser) expect_serie(tokens *[]lexer.Token, tts ...lexer.TokenType) bool {
	if p.eof(uint(len(tts))) {
		return false;
	}
	for i, tt := range tts {
		if token := p.tokens[p.current + i]; token.Type != tt {
			return false;
		}
	}
	p.current += len(tts);
	if tokens != nil {
		*tokens = p.tokens[p.current-len(tts):p.current];
	}
	return true;
}

// block -> "{" declarativestmt "}"
func (p *Parser) consume_block() ([]Stmt, error) {
	stmts := make([]Stmt, 0);
	for !p.expect(lexer.RIGHT_BRACE) {
		if p.eof(0) {
			return nil, p.generate_expect_error("} at the end of the block");
		}
		stmt, err := p.declarative_statement();
		if err != nil {
			return nil, err;
		}
		stmts = append(stmts, stmt);
	}
	return stmts, nil;
}

// condstmt -> "if" expr stmt ("else" condstmt | stmt)?
func (p *Parser) consume_if(branches *[]ConditionalBranch) error {
	if !p.expect(lexer.LEFT_PAREN) {
		return p.generate_expect_error("( in if condition");
	}
	cond, err := p.expression();
	if err != nil {
		return err;
	}
	if !p.expect(lexer.RIGHT_PAREN) {
		return p.generate_expect_error(") in if condition");
	}
	ndstmt, err := p.statement();
	if err != nil {
		return err;
	}
	*branches = append(*branches, ConditionalBranch{
		Condition: cond,
		NDStmt: ndstmt,
	});
	if p.expect(lexer.ELSE) {
		if p.expect(lexer.IF) {
			return p.consume_if(branches);
		}
		ndstmt, err = p.statement();
		if err != nil {
			return err;
		}
		*branches = append(*branches, ConditionalBranch{
			Condition: nil,
			NDStmt: ndstmt,
		});
	}
	return nil;
}

// func -> IDENTIFIER "(" params? ")" block
func (p *Parser) consume_func() (*Func, error) {
	if !p.expect(lexer.IDENTIFIER) {
		return nil, p.generate_expect_error("IDENTIFIER in function signature");
	}
	name := p.prev();
	if !p.expect(lexer.LEFT_PAREN) {
		return nil, p.generate_expect_error("'(' in function signature");
	}
	params := make([]lexer.Token, 0);
	if !p.expect(lexer.RIGHT_PAREN) {
		if err := p.consume_func_params(&params); err != nil {
			return nil, err;
		}
		if !p.expect(lexer.RIGHT_PAREN) {
			return nil, p.generate_expect_error("')' in function signature");
		}
	}
	if !p.expect(lexer.LEFT_BRACE) {
		return nil, p.generate_expect_error("'{' to start function body")
	}
	body, err := p.consume_block();
	if err != nil {
		return nil, err;
	}
	return &Func{
		Name: name,
		Params: params,
		Body: body,
	}, nil;
}

// params -> IDENTIFIER | (IDENTIFIER "," params)
func (p *Parser) consume_func_params(params *[]lexer.Token) error {
	if !p.expect(lexer.IDENTIFIER) {
		return p.generate_expect_error("IDENTIFIER as a parameter")
	}
	*params = append(*params, p.prev());
	if p.expect(lexer.COMMA) {
		return p.consume_func_params(params);
	}
	return nil;
}

// args -> expression | (expression "," args)
func (p *Parser) consume_func_args(params *[]Expr) error {
	val, err := p.expression();
	if err != nil {
		return err;
	}
	*params = append(*params, val);
	if p.expect(lexer.COMMA) {
		return p.consume_func_args(params);
	}
	return nil;
}

// recursive decent start
func (p *Parser) declarative_statement() (Stmt, error) {
	// var -> "var" IDENTIFIER ("=" expression)?
	if p.expect(lexer.VAR) {
		var (
			asset Expr = nil;
			err error
		);
		if !p.expect(lexer.IDENTIFIER) {
			return nil, p.generate_expect_error("IDENTIFIER in variable declartion");
		}
		id := p.prev();
		if p.expect(lexer.EQUAL) {
			asset, err = p.expression();
			if err != nil {
				return nil, err;
			}
		}
		if !p.expect(lexer.SEMICOLON) {
			return nil, p.generate_expect_error("';' at the end of the statement.");
		}
		return VarDeclarationStmt{
			Name: id,
			Asset: asset,
		}, nil;
	}
	// funcdecl -> "func" func
	if p.expect(lexer.FUNC) {
		fn, err := p.consume_func();
		if err != nil {
			return nil, err;
		}
		return (*FuncDeclarationStmt)(fn), nil;
	}
	return p.statement();
}

func (p *Parser) statement() (Stmt, error) {
	if p.expect(lexer.IF) {
		branches := make([]ConditionalBranch, 0);
		if err := p.consume_if(&branches); err != nil {
			return nil, err;
		}
		return ConditionalStmt {
			Branches: branches,
		}, nil;
	}
	// whileloop -> "while" expression statement
	if p.expect(lexer.WHILE) {
		if !p.expect(lexer.LEFT_PAREN) {
			return nil, p.generate_expect_error("( in while loop condition");
		}
		cond, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if !p.expect(lexer.RIGHT_PAREN) {
			return nil, p.generate_expect_error(") in while loop condition");
		}
		ndstmt, err := p.statement();
		if err != nil {
			return nil, err;
		}
		return WhileStmt{
			Cond: cond,
			NDStmt: ndstmt,
		}, nil;
	}
	// forloop -> "for" "(" declarative_statement? ";" expession? ";" expression? ")" statement
	if p.expect(lexer.FOR) {
		var (
			init Stmt = nil;
			cond Expr = nil;
			step Expr = nil;
			err error = nil;
		);
		if !p.expect(lexer.LEFT_PAREN) {
			return nil, p.generate_expect_error("( in for loop header");
		}
		if !p.expect(lexer.SEMICOLON) {
			init, err = p.declarative_statement()
			if err != nil {
				return nil, err;
			}
		}
		if !p.expect(lexer.SEMICOLON) {
			cond, err = p.expression();
			if err != nil {
				return nil, err;
			}
			if !p.expect(lexer.SEMICOLON) {
				return nil, p.generate_expect_error("; in for loop header");
			}
		}
		if !p.expect(lexer.RIGHT_PAREN) {
			step, err = p.expression();
			if err != nil {
				return nil, err;
			}
			if !p.expect(lexer.RIGHT_PAREN) {
				return nil, p.generate_expect_error(") after for loop header");
			}
		}
		ndstmt, err := p.statement();
		if err != nil {
			return nil, err;
		}
		return ForStmt {
			Init: init,
			Cond: cond,
			Step: step,
			NDStmt: ndstmt,
		}, nil;
	}
	if p.expect(lexer.LEFT_BRACE) {
		stmts, err := p.consume_block();
		if err != nil {
			return nil, err;
		}
		return BlockStmt{
			Stmts: stmts,
		}, nil;
	}
	// return -> "return" expression? ";"
	if p.expect(lexer.RETURN) {
		var ( expr Expr = nil; err error = nil; )
		if !p.expect(lexer.SEMICOLON) {
			expr, err = p.expression();
			if err != nil {
				return nil, err;
			}
			if !p.expect(lexer.SEMICOLON) {
				return nil, p.generate_expect_error("';' at the end of the return statement");
			}
		}
		return ReturnStmt{
			Asset: expr,
		}, nil;
	}
	if p.expect(lexer.BREAK) {
		if !p.expect(lexer.SEMICOLON) {
			return nil, p.generate_expect_error("';' at the end of 'break'");
		}
		return BreakStmt{}, nil;
	}
	if p.expect(lexer.CONTINUE) {
		if !p.expect(lexer.SEMICOLON) {
			return nil, p.generate_expect_error("';' at the end of 'continue'");
		}
		return ContinueStmt{}, nil;
	}
	// printstmt -> "print" expression ("," expression)* ";"
	if p.expect(lexer.PRINT) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		assets := make([]Expr, 1);
		assets[0] = expr;
		for p.expect(lexer.COMMA) {
			expr, err := p.expression();
			if err != nil {
				return nil, err;
			}
			assets = append(assets, expr);
		}
		if !p.expect(lexer.SEMICOLON) {
			return nil, p.generate_expect_error("';' at the end of the statement.");
		}
		return PrintStmt{
			Assets: assets,
		}, nil;
	}
	// exprstmt -> expression ";"
	expr, err := p.expression();
	if err != nil {
		return nil, err;
	}
	if !p.expect(lexer.SEMICOLON) {
		return nil, p.generate_expect_error("';' at the end of the statement.");
	}
	return ExprStmt{
		InnerExpr: expr,
	}, nil;
}

// expression -> assign 
func (p *Parser) expression() (Expr, error) {
	return p.assign();
}

// expressions -> assign "," assign | expressions* 
// func (p *Parser) expressions() (Expr, error) {
// 	expr, err := p.assign();
// 	if err != nil {
// 		return nil, err;
// 	}
// 	if p.expect(lexer.COMMA) {
// 		operator := p.prev();
// 		right, err := p.assign();
// 		if err != nil {
// 			return nil, err;
// 		}
// 		return BinaryExpr{
// 			LOperand: expr,
// 			Operator: operator,
// 			ROperand: right,
// 		}, nil;
// 	}
// 	return expr, nil;
// }

// assign -> IDENTIFIER "=" assign
func (p *Parser) assign() (Expr, error) {
	tokens := []lexer.Token{};
	if p.expect_serie(&tokens, lexer.IDENTIFIER, lexer.EQUAL) {
		src, err := p.assign();
		if err != nil {
			return nil, err;
		}
		return AssignExpr{
			Name: tokens[0],
			Asset: src,
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
		return TernaryExpr{
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
		expr = BinaryExpr {
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
		expr = BinaryExpr {
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
		expr = BinaryExpr {
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
		expr = BinaryExpr {
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
		operator := p.prev();
		operand, err := p.unary();
		if err != nil {
			return nil, err;
		}
		return UnaryExpr {
			Operand: operand,
			Operator: operator,
		}, nil;
	}
	return p.call();
}

// call -> IDENTIFIER ( "(" funcparams ")" )+
func (p *Parser) call() (Expr, error) {
	expr, err := p.primary();
	if err != nil {
		return nil, err;
	}
	for ;; {
		if p.expect(lexer.LEFT_PAREN) {
			args := make([]Expr, 0);
			if !p.expect(lexer.RIGHT_PAREN) {
				if err := p.consume_func_args(&args); err != nil {
					return nil, err;
				}
				if !p.expect(lexer.RIGHT_PAREN) {
					return nil, p.generate_expect_error("')' in function call");
				}
			}
			expr = FuncCall{
				Callee: expr,
				Args: args,
			};
		} else {
			break;
		}
	}
	return expr, nil;
}

// primary -> IDENTIFIER | STRING | NUMBER | "true" | "false" | "null" | "(" expression ")"
func (p *Parser) primary() (Expr, error) {
	if p.expect(lexer.TRUE) {
		return LiteralExpr{
			ValueLiteral: true,
		}, nil
	} else if p.expect(lexer.FALSE) {
		return LiteralExpr{
			ValueLiteral: false,
		}, nil
	} else if p.expect(lexer.NULL) {
		return LiteralExpr{
			ValueLiteral: nil,
		}, nil
	} else if p.expect(lexer.STRING, lexer.NUMBER) {
		return LiteralExpr {
			ValueLiteral: p.prev().Literal,
		}, nil
	} else if p.expect(lexer.IDENTIFIER) {
		return VariableExpr{
			Name: p.prev(),
		}, nil;
	} else if p.expect(lexer.LEFT_PAREN) {
		expr, err := p.expression();
		if err != nil {
			return nil, err;
		}
		if p.expect(lexer.RIGHT_PAREN) {
			return GroupingExpr{
				InnerExpr: expr,
			}, nil;
		}
		return nil, p.generate_expect_error(")");
	}
	return nil, p.generate_expect_error("valid token");
}
// recursive decent end

func (p *Parser) generate_expect_error(expected string) error {
	tok := p.tokens[p.current];
	return fmt.Errorf("Parser Error: expected %s found %s at line %d\n", expected, tok.Lexeme, tok.Line);
}

func (p *Parser) Parse() ([]Stmt, error) {
	stmts := make([]Stmt, 0);
	for !p.eof(0) {
		stmt, err := p.declarative_statement();
		if err != nil {
			return nil, err;
		}
		stmts = append(stmts, stmt);
	}
	return stmts, nil;
}
