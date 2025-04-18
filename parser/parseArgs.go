package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseArgs() ([]ast.Expr, *errors.SyntaxError) {
	_, err := p.expect(lexer.OParen)
	if err != nil {
		return nil, err
	}
	var args []ast.Expr
	if p.at().Type == lexer.CParen {
		args = []ast.Expr{}
	} else {
		newArgs, err := p.parseArgsList()
		if err != nil {
			return nil, err
		}
		args = newArgs
	}
	_, err = p.expect(lexer.CParen)
	if err != nil {
		return nil, err
	}
	return args, nil
}
