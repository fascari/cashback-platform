package validator

import "strings"

type Errors map[string]string

func (e Errors) Error() string {
	parts := make([]string, 0, len(e))
	for field, msg := range e {
		parts = append(parts, field+": "+msg)
	}
	return strings.Join(parts, "; ")
}
