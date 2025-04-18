package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
)

func (p *Parser) parseReturnStmt() (ast.Expr, *errors.SyntaxError) {
	p.advance() // return
	var value ast.Expr

	if p.isEOF() {
		value = nil
	} else {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		value = expr
	}

	return &ast.ReturnStmt{
		Value: value,
	}, nil
}
