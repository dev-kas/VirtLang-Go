package parser

import (
	"strconv"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parsePrimaryExpr() (ast.Expr, *errors.SyntaxError) {
	start := p.at()
	tk := start.Type
	var value interface{}

	switch tk {
	case lexer.Identifier:
		return &ast.Identifier{
			Symbol: p.advance().Literal,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil

	case lexer.Number:
		value = p.advance().Literal
		parsedValue, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return nil, &errors.SyntaxError{
				Expected: "Number",
				Got:      value.(string),
				Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
				End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
			}
		}
		return &ast.NumericLiteral{
			Value: parsedValue,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil

	case lexer.OParen:
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		p.expect(lexer.CParen)
		return expr, nil

	case lexer.OBracket:
		return p.parseArrayLiteral()

	case lexer.String:
		value = p.advance().Literal
		return &ast.StringLiteral{
			Value: value.(string),
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil

	case lexer.WhileLoop:
		return p.parseWhileLoop()

	case lexer.Comment:
		p.advance()
		var result *ast.Expr

		if p.isEOF() {
			result = nil
		} else {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			result = &expr
		}

		return *result, nil

	case lexer.Try:
		return p.parseTryCatch()

	case lexer.Return:
		return p.parseReturnStmt()

	case lexer.Break:
		return p.parseBreakStmt()

	case lexer.Continue:
		return p.parseContinueStmt()

	default:
		return nil, &errors.SyntaxError{
			Expected: "Primary Expression",
			Got:      lexer.Stringify(tk),
			Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
			End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
		}
	}
}
