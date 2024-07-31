package ast

import (
	"dodo-lang/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foobar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "barbaz"},
					Value: "barbaz",
				},
			},
		},
	}

	if program.String() != "let foo = barbaz;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
