package router

import "strings"

type ValidationErr interface {
	Hint(field string, hint string)
	Err() bool
	error
}

const (
	HintMissingOrZero   = "is missing or zero"
	HintEmptyString     = "is left empty"
	HintIncorrectFormat = "has incorrect format"
	HintMinimumLength8  = "should be at least 8 characters long"
)

type validationErr struct {
	ValidationErr
	hints map[string][]string
}

func NewValidationErr() validationErr {
	return validationErr{
		hints: make(map[string][]string, 0),
	}
}

func (ve validationErr) Hint(field string, hint string) {
	ve.hints[field] = append(ve.hints[field], hint)
}

func (ve validationErr) Err() bool {
	return len(ve.hints) > 0
}

func (ve validationErr) Error() string {
	if len(ve.hints) == 0 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString("validation errors found:\n")
	for field, hints := range ve.hints {
		prefix := "- field '" + field + "' "
		for _, hint := range hints {
			builder.WriteString(prefix)
			builder.WriteString(hint)
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
