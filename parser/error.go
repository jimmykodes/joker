package parser

import (
	"fmt"

	"github.com/jimmykodes/joker/token"
)

type ParserError struct {
	Token   token.Token
	Message string
	Err     error
}

func (e ParserError) Unwrap() error {
	return e.Err
}

func (e ParserError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("parser: token %v: %s: %s", e.Token, e.Message, e.Err)
	}
	return fmt.Sprintf("parser: token: %v: %s", e.Token, e.Message)
}
