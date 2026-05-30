package python

import (
	"fmt"
	"strings"
	"unicode"
)

var keywords = map[string]struct{}{
	"False": {}, "None": {}, "True": {}, "and": {}, "as": {}, "assert": {},
	"async": {}, "await": {}, "break": {}, "class": {}, "continue": {},
	"def": {}, "del": {}, "elif": {}, "else": {}, "except": {}, "finally": {},
	"for": {}, "from": {}, "global": {}, "if": {}, "import": {}, "in": {},
	"is": {}, "lambda": {}, "nonlocal": {}, "not": {}, "or": {}, "pass": {},
	"raise": {}, "return": {}, "try": {}, "while": {}, "with": {}, "yield": {},
}

func ValidateIdentifier(name string) error {
	if name == "" {
		return fmt.Errorf("identifier must not be empty")
	}

	for i, r := range name {
		if i == 0 {
			if r != '_' && !unicode.IsLetter(r) {
				return fmt.Errorf("identifier %q must start with a letter or underscore", name)
			}
			continue
		}

		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("identifier %q contains invalid character %q", name, r)
		}
	}

	if _, ok := keywords[name]; ok {
		return fmt.Errorf("identifier %q is a Python keyword", name)
	}

	return nil
}

func ValidateDottedIdentifier(name string) error {
	if name == "" {
		return fmt.Errorf("dotted identifier must not be empty")
	}

	parts := strings.Split(name, ".")
	for _, part := range parts {
		if err := ValidateIdentifier(part); err != nil {
			return fmt.Errorf("invalid dotted identifier %q: %w", name, err)
		}
	}

	return nil
}
