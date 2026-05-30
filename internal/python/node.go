package python

import "fmt"

type Expression struct {
	code string
}

func (e Expression) Code() string {
	return e.code
}

type Statement struct {
	code string
}

func (s Statement) Code() string {
	return s.code
}

func RawExpression(code string) (Expression, error) {
	if code == "" {
		return Expression{}, fmt.Errorf("expression code must not be empty")
	}

	return Expression{code: code}, nil
}

func Reference(name string) (Expression, error) {
	if err := ValidateDottedIdentifier(name); err != nil {
		return Expression{}, err
	}

	return Expression{code: name}, nil
}

func Attribute(receiver Expression, name string) (Expression, error) {
	if receiver.code == "" {
		return Expression{}, fmt.Errorf("attribute receiver must not be empty")
	}

	if err := ValidateIdentifier(name); err != nil {
		return Expression{}, err
	}

	return Expression{code: receiver.code + "." + name}, nil
}

func Assign(name string, value Expression) (Statement, error) {
	if err := ValidateIdentifier(name); err != nil {
		return Statement{}, err
	}

	if value.code == "" {
		return Statement{}, fmt.Errorf("assignment value must not be empty")
	}

	return Statement{code: name + " = " + value.code}, nil
}

func ExpressionStatement(expr Expression) (Statement, error) {
	if expr.code == "" {
		return Statement{}, fmt.Errorf("expression statement must not be empty")
	}

	return Statement{code: expr.code}, nil
}
