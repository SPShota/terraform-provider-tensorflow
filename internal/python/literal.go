package python

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func Literal(value any) (Expression, error) {
	code, err := literalCode(value)
	if err != nil {
		return Expression{}, err
	}

	return Expression{code: code}, nil
}

func literalCode(value any) (string, error) {
	if value == nil {
		return "None", nil
	}

	switch v := value.(type) {
	case Expression:
		if v.code == "" {
			return "", fmt.Errorf("expression literal must not be empty")
		}
		return v.code, nil
	case string:
		return strconv.Quote(v), nil
	case bool:
		if v {
			return "True", nil
		}
		return "False", nil
	case int:
		return strconv.Itoa(v), nil
	case int8, int16, int32, int64:
		return fmt.Sprintf("%d", v), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return finiteFloat(float64(v), 32)
	case float64:
		return finiteFloat(v, 64)
	case json.Number:
		if _, err := v.Float64(); err != nil {
			return "", fmt.Errorf("invalid JSON number literal %q", v)
		}
		return v.String(), nil
	case []any:
		return sequenceLiteral("[", "]", v)
	case map[string]any:
		return mapLiteral(v)
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		items := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			items = append(items, rv.Index(i).Interface())
		}
		return sequenceLiteral("[", "]", items)
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return "", fmt.Errorf("map literal keys must be strings, got %s", rv.Type().Key())
		}

		values := make(map[string]any, rv.Len())
		for _, key := range rv.MapKeys() {
			values[key.String()] = rv.MapIndex(key).Interface()
		}
		return mapLiteral(values)
	}

	return "", fmt.Errorf("unsupported literal type %T", value)
}

func finiteFloat(value float64, bitSize int) (string, error) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return "", fmt.Errorf("float literal must be finite")
	}

	return strconv.FormatFloat(value, 'g', -1, bitSize), nil
}

func sequenceLiteral(open, close string, values []any) (string, error) {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		code, err := literalCode(value)
		if err != nil {
			return "", err
		}
		parts = append(parts, code)
	}

	return open + strings.Join(parts, ", ") + close, nil
}

func mapLiteral(values map[string]any) (string, error) {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		code, err := literalCode(values[key])
		if err != nil {
			return "", err
		}
		parts = append(parts, strconv.Quote(key)+": "+code)
	}

	return "{" + strings.Join(parts, ", ") + "}", nil
}
