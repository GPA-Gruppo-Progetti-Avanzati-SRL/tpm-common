package ceLexer

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/lexer"

func NewLexer(text string) (*lexer.L, error) {
	l := lexer.New(text, ZeroState)
	l.Start()
	return l, nil
}
