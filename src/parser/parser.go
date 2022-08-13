package parser

import (
	"fmt"
	"mnk/src/ast"
	"mnk/src/lexer"
	"mnk/src/token"
	"strconv"
)

// ============================== TYPES

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     //*
	PREFIX      // -X or !X
	CALL        // myFunction()
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l              *lexer.Lexer
	errors         []string
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// ============================== PARSER COMMON METHODS

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// read 2 tokens to set both curToken and nextToken
	p.nextToken()
	p.nextToken()

	// prefixParseFns init
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// infixParseFns init
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// ============================== PARSER prefixParseFn HOOKS

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    // defer untrace(trace("parseIntegerLiteral"))

	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    // defer untrace(trace("parsePrefixExpression"))

	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    // defer untrace(trace("parseInfixExpression"))

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()

	p.nextToken()

		expression.Right = p.parseExpression(precedence)

	return expression
}

// ============================== PARSER ParseProgram METHOD

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}              // build root node
	program.Statements = []ast.Statement{} // build []statements nodes

	// iterate tokens until EOF
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// ============================== PARSER parseStatement DISPATCHER METHOD

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// ============================== PARSER parseExpressionStatement METHOD

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    // defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    // defer untrace(trace("parseExpression"))

	// try to find a prefix hook function for the current token
	prefix := p.prefixParseFns[p.curToken.Type]

	// if not found print error and return nil
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// else execute it
	leftExp := prefix()

	// until not SEMICOLON AND
	// until given precedence arg is < that of the next token
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// try to find an infix hook function for the next token
		infix := p.infixParseFns[p.peekToken.Type]

		// if not found return leftExp (prefix)
		if infix == nil {
			return leftExp
		}

		// else advance AND
		p.nextToken()
		// execute it (w/ arg leftExp)
		leftExp = infix(leftExp)
	}

	return leftExp
}

// ============================== PARSER parseReturnStatement METHOD

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken} // build ReturnStatement node w/ token.RETURN

	p.nextToken()

	// skip expressions for now

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ============================== PARSER parseLetStatement METHOD

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken} // build LetStatement node w/ token.LET

	// expect identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// build Identifier node
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// expect =
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// skip expressions for now

	// expect ;
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
