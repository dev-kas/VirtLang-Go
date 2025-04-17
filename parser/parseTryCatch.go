package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseTryCatch() (ast.Expr, *errors.SyntaxError) {
	p.advance() // try
	if _, err := p.expect(lexer.OBrace); err != nil {
		return nil, err
	}

	tryBody := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		tryBody = append(tryBody, stmt)
	}

	if _, err := p.expect(lexer.CBrace); err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.Catch); err != nil {
		return nil, err
	}

	cVar, err := p.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.OBrace); err != nil {
		return nil, err
	}

	catchBody := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		catchBody = append(catchBody, stmt)
	}

	if _, err := p.expect(lexer.CBrace); err != nil {
		return nil, err
	}

	return &ast.TryCatchStmt{
		Try:      tryBody,
		CatchVar: cVar.Literal,
		Catch:    catchBody,
	}, nil
}
