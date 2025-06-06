package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalVarAssignment(node *ast.VarAssignmentExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	if node.Assignee.GetType() == ast.IdentifierNode {
		varname := node.Assignee.(*ast.Identifier).Symbol

		value, err := Evaluate(node.Value, env, dbgr)
		if err != nil {
			return nil, err
		}

		return env.AssignVar(varname, *value)
	} else if node.Assignee.GetType() == ast.MemberExprNode {
		memberExpr := node.Assignee.(*ast.MemberExpr)
		obj, err := Evaluate(memberExpr.Object, env, dbgr)
		if err != nil {
			return nil, err
		}

		if obj.Type != shared.Object && obj.Type != shared.Array {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Cannot access property of non-object (attempting to access properties of %v).", shared.Stringify(obj.Type)),
			}
		}

		if obj.Type == shared.Array {
			indexVal, err := Evaluate(memberExpr.Value, env, dbgr)
			if err != nil {
				return nil, err
			}

			if indexVal.Type != shared.Number {
				return nil, &errors.RuntimeError{
					Message: fmt.Sprintf("Cannot assign to array using non-number index (attempted to use %v).", shared.Stringify(indexVal.Type)),
				}
			}

			index := int(indexVal.Value.(float64))
			array := obj.Value.([]shared.RuntimeValue)

			if index < 0 {
				return nil, &errors.RuntimeError{
					Message: fmt.Sprintf("Index out of bounds: %d", index),
				}
			}

			if index >= len(array) {
				extendedArray := make([]shared.RuntimeValue, index+1)
				copy(extendedArray, array)
				for i := len(array); i <= index; i++ {
					extendedArray[i] = values.MK_NIL()
				}
				obj.Value = extendedArray
				array = extendedArray
				obj.Value = array
			}

			value, err := Evaluate(node.Value, env, dbgr)
			if err != nil {
				return nil, err
			}

			array[index] = *value
			obj.Value = array

			switch memberExpr.Object.GetType() {
			case ast.IdentifierNode:
				varName := memberExpr.Object.(*ast.Identifier).Symbol
				_, err = env.AssignVar(varName, *obj)
				if err != nil {
					return nil, err
				}
			case ast.MemberExprNode:
				memberExpr := node.Assignee.(*ast.MemberExpr)
				obj, err := Evaluate(memberExpr.Object, env, dbgr)
				if err != nil {
					return nil, err
				}

				switch obj.Type {
				case shared.Object:
					var prop *shared.RuntimeValue
					if memberExpr.Computed {
						val, err := Evaluate(memberExpr.Value, env, dbgr)
						if err != nil {
							return nil, err
						}
						prop = val
					} else {
						prop = &shared.RuntimeValue{
							Type:  shared.String,
							Value: memberExpr.Value.(*ast.Identifier).Symbol,
						}
					}

					var key string
					switch v := prop.Value.(type) {
					case string:
						key = v
					case int:
						key = fmt.Sprintf("%v", v)
					default:
						return nil, &errors.RuntimeError{
							Message: fmt.Sprintf("Invalid property key type: %T", prop.Value),
						}
					}

					obj.Value.(map[string]*shared.RuntimeValue)[key] = value
					return value, nil
				}
			}

			return value, nil
		}

		var prop *shared.RuntimeValue
		if memberExpr.Computed {
			val, err := Evaluate(memberExpr.Value, env, dbgr)
			if err != nil {
				return nil, err
			}

			prop = val
		} else {
			prop = &shared.RuntimeValue{
				Type:  shared.String,
				Value: memberExpr.Value.(*ast.Identifier).Symbol,
			}
		}

		var key string

		if prop.Type == shared.ValueType(ast.IdentifierNode) {
			identifier, ok := prop.Value.(*ast.Identifier)
			if !ok {
				return nil, &errors.RuntimeError{
					Message: "Cannot access property of object by non-string key.",
				}
			}
			key = identifier.Symbol
		} else {
			// key = prop.Value.(shared.RuntimeValue).Value.(string)
			key = prop.Value.(string)
		}

		value, err := Evaluate(node.Value, env, dbgr)
		if err != nil {
			return nil, err
		}

		obj.Value.(map[string]*shared.RuntimeValue)[key] = value

		return value, nil
	} else {
		return nil, &errors.RuntimeError{
			Message: "Cannot access property of object by non-string key.",
		}
	}
}
