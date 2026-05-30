package python

import (
	"fmt"
	"strings"
)

func Return(value Expression) (Statement, error) {
	if value.code == "" {
		return Statement{}, fmt.Errorf("return value must not be empty")
	}

	return Statement{code: "return " + value.code}, nil
}

func With(context Expression, alias string, body []Statement) (Statement, error) {
	if context.code == "" {
		return Statement{}, fmt.Errorf("with context must not be empty")
	}

	header := "with " + context.code
	if alias != "" {
		if err := ValidateIdentifier(alias); err != nil {
			return Statement{}, err
		}
		header += " as " + alias
	}

	return Block(header, body)
}

func FunctionDef(name string, args []string, decorators []Expression, body []Statement) (Statement, error) {
	if err := ValidateIdentifier(name); err != nil {
		return Statement{}, err
	}

	for _, arg := range args {
		if err := ValidateIdentifier(arg); err != nil {
			return Statement{}, fmt.Errorf("invalid function argument: %w", err)
		}
	}

	header := "def " + name + "(" + strings.Join(args, ", ") + ")"
	block, err := Block(header, body)
	if err != nil {
		return Statement{}, err
	}

	if len(decorators) == 0 {
		return block, nil
	}

	lines := make([]string, 0, len(decorators)+1)
	for _, decorator := range decorators {
		if decorator.code == "" {
			return Statement{}, fmt.Errorf("function decorator must not be empty")
		}
		lines = append(lines, "@"+decorator.code)
	}
	lines = append(lines, block.code)

	return Statement{code: strings.Join(lines, "\n")}, nil
}

func Block(header string, body []Statement) (Statement, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return Statement{}, fmt.Errorf("block header must not be empty")
	}
	if strings.HasSuffix(header, ":") {
		header = strings.TrimSuffix(header, ":")
	}

	bodyCode, err := renderBlockBody(body)
	if err != nil {
		return Statement{}, err
	}

	return Statement{code: header + ":\n" + bodyCode}, nil
}

func renderBlockBody(body []Statement) (string, error) {
	if len(body) == 0 {
		return "\tpass", nil
	}

	lines := make([]string, 0, len(body))
	for _, statement := range body {
		if statement.code == "" {
			return "", fmt.Errorf("block body statement must not be empty")
		}
		lines = append(lines, indent(statement.code))
	}

	return strings.Join(lines, "\n"), nil
}

func indent(code string) string {
	code = strings.TrimRight(code, "\n")
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = "\t"
			continue
		}
		lines[i] = "\t" + line
	}
	return strings.Join(lines, "\n")
}
