package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseIfStmt() (*ast.IfStatement, *errors.SyntaxError) {
	p.advance() // if

	if _, err := p.expect(lexer.OParen); err != nil {
		return nil, err
	}

	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.CParen); err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.OBrace); err != nil {
		return nil, err
	}

	body := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	if _, err := p.expect(lexer.CBrace); err != nil {
		return nil, err
	}

	ifStmt := &ast.IfStatement{
		Condition: condition,
		Body:      body,
		ElseIf:    []*ast.IfStatement{}, // Initialize ElseIf as an empty slice
	}

	// Check for else or else if
	for p.at().Type == lexer.Else {
		p.advance() // else

		if p.at().Type == lexer.If {
			// Else if
			elseIfStmt, err := p.parseIfStmt()
			if err != nil {
				return nil, err
			}
			ifStmt.ElseIf = append(ifStmt.ElseIf, elseIfStmt)
		} else if p.at().Type == lexer.OBrace {
			// Else
			p.advance() // {
			elseBody := []ast.Stmt{}

			for !p.isEOF() && p.at().Type != lexer.CBrace {
				stmt, err := p.parseStmt()
				if err != nil {
					return nil, err
				}
				elseBody = append(elseBody, stmt)
			}

			if _, err := p.expect(lexer.CBrace); err != nil {
				return nil, err
			}

			ifStmt.Else = elseBody
			break
		} else {
			return nil, errors.NewSyntaxError("Unexpected token after else", lexer.Stringify(p.at().Type), p.at().Start, p.at().Difference)
		}
	}

	return ifStmt, nil
}
