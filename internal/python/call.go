package python

import (
	"fmt"
	"strings"
)

type KeywordArgument struct {
	Name  string
	Value Expression
}

func Call(function Expression, args []Expression, kwargs []KeywordArgument) (Expression, error) {
	if function.code == "" {
		return Expression{}, fmt.Errorf("call function must not be empty")
	}

	parts := make([]string, 0, len(args)+len(kwargs))
	for _, arg := range args {
		if arg.code == "" {
			return Expression{}, fmt.Errorf("call argument must not be empty")
		}
		parts = append(parts, arg.code)
	}

	seen := make(map[string]struct{}, len(kwargs))
	for _, kwarg := range kwargs {
		if err := ValidateIdentifier(kwarg.Name); err != nil {
			return Expression{}, fmt.Errorf("invalid keyword argument: %w", err)
		}

		if _, ok := seen[kwarg.Name]; ok {
			return Expression{}, fmt.Errorf("duplicate keyword argument %q", kwarg.Name)
		}
		seen[kwarg.Name] = struct{}{}

		if kwarg.Value.code == "" {
			return Expression{}, fmt.Errorf("keyword argument %q must not be empty", kwarg.Name)
		}

		parts = append(parts, kwarg.Name+"="+kwarg.Value.code)
	}

	return Expression{code: function.code + "(" + strings.Join(parts, ", ") + ")"}, nil
}

func CallNamed(function string, args []Expression, kwargs []KeywordArgument) (Expression, error) {
	ref, err := Reference(function)
	if err != nil {
		return Expression{}, err
	}

	return Call(ref, args, kwargs)
}
