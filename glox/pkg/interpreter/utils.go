package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

func isTruthy(value any) bool {
	if value == nil {
		return false
	}

	if value, ok := value.(bool); ok {
		return value
	}

	return true
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}

func stringify(value any) string {
	if value == nil {
		return "<nil>"
	}

	if value, ok := value.(float64); ok {
		text := fmt.Sprintf("%g", value)
		if strings.HasSuffix(text, ".0") {
			return text[:len(text)-2]
		}

		return text
	}

	if value, ok := value.(string); ok {
		return strconv.Quote(value)
	}

	return fmt.Sprintf("%v", value)
}
