package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Richtermnd/ferret/ast"
	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	curToken  token.Token
	l         *lexer.Lexer
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	// --- prefix ---
	p.prefixParseFns[token.NOT] = p.parsePrefixExpression
	p.prefixParseFns[token.SUB] = p.parsePrefixExpression

	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.FLOAT] = p.parseFloatLiteral
	p.prefixParseFns[token.TRUE] = p.parseBooleanLiteral
	p.prefixParseFns[token.FALSE] = p.parseBooleanLiteral

	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression

	// --- infix ---
	p.infixParseFns[token.ADD] = p.parseInfixExpression
	p.infixParseFns[token.SUB] = p.parseInfixExpression
	p.infixParseFns[token.MUL] = p.parseInfixExpression
	p.infixParseFns[token.DIV] = p.parseInfixExpression
	p.infixParseFns[token.REM] = p.parseInfixExpression
	p.infixParseFns[token.REM] = p.parseInfixExpression

	p.infixParseFns[token.EQ] = p.parseInfixExpression
	p.infixParseFns[token.NEQ] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression
	p.infixParseFns[token.GEQ] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.LEQ] = p.parseInfixExpression
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() *ast.Program {
	program := new(ast.Program)
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if !p.HasErrors() {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// Here will be other tokens like var, functions declarations assignment and other
	// Everything other - expressions
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	p.nextToken()
	if !p.curToken.Is(token.IDENT) {
		p.errors = append(p.errors, fmt.Errorf("expected identifier, got %v", p.curToken.Type))
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	if !p.curToken.Is(token.ASSIGN) {
		p.errors = append(p.errors, fmt.Errorf("expected =, got %v", p.curToken.Type))
	}
	p.nextToken()
	stmt.Value = p.parseExpression(token.LOWEST)
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expr = p.parseExpression(token.LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParser, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.errors = append(p.errors, fmt.Errorf("no prefix parsers for %s", p.curToken.Literal))
		return nil
	}

	leftExp := prefixParser()
	for (!p.peekToken.Is(token.LF) || !p.peekToken.Is(token.SEMICOLON)) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	if p.peekToken.Is(token.LF) || p.peekToken.Is(token.SEMICOLON) {
		p.nextToken()
	}
	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("not a valid int %s", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("not a valid float %s", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Is(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	exp.Right = p.parseExpression(token.UNARY)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(token.LOWEST)
	if !p.peekToken.Is(token.RPAREN) {
		p.errors = append(p.errors, fmt.Errorf("no closing )"))
		return nil
	}
	p.nextToken()
	return exp
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curPrecedence() int {
	return p.curToken.Precedence()
}

func (p *Parser) peekPrecedence() int {
	return p.peekToken.Precedence()
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) != 0
}

func (p *Parser) PrintErrors(out io.Writer) {
	for _, err := range p.Errors() {
		fmt.Fprintln(out, err)
	}
}

func (p *Parser) Errors() []error {
	return p.errors
}
