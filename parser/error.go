package parser

import (
	"errors"
	"fmt"

	"github.com/jimmykodes/jk/token"
)

var ErrParserError = errors.New("parser error")

func newParseError(line int, message string, args ...any) error {
	return fmt.Errorf("%w at line %d: %s", ErrParserError, line, fmt.Sprintf(message, args...))
}

func invalidTokenError(line int, expected, got token.Type) error {
	return newParseError(line, "invalid token. expected: %s - got: %s", expected, got)
}

type ParseError struct {
	error
	Token token.Token
}
