package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
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
