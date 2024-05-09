package lexer

import "fmt"

type LexerError struct {
	Err error
}

func (e LexerError) Unwrap() error {
	return e.Err
}

func (e LexerError) Error() string {
	return fmt.Sprintf("lexer: %s", e.Err)
}
