package main

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestLexer_next(t *testing.T) {
	reader, writer := io.Pipe()
	lex := &lexer{
		input:   bufio.NewReader(reader),
		state:   lexText,
		readbuf: make([]byte, 0, 3), // test with small buffer
		items:   make(chan item),
	}

	out := make(chan rune)
	go func() {
		defer close(out)
		for {
			ch := lex.next()
			if ch == eof {
				break
			}
			out <- ch
		}
	}()
	go func() {
		writer.Write([]byte("hello 你好"))
		writer.Write([]byte(" 世界 world\x00"))
	}()

	result := make([]rune, 0, 1024)
	for ch := range out {
		result = append(result, ch)
	}

	if want, have := "hello 你好 世界 world", string(result); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestLexer_backup(t *testing.T) {

	reader, writer := io.Pipe()
	lex := &lexer{
		input:   bufio.NewReader(reader),
		state:   lexText,
		readbuf: make([]byte, 0, 3), // test with small buffer
		items:   make(chan item),
	}

	go func() {
		writer.Write([]byte("hello 你好"))
		writer.Write([]byte(" 世界 world\x00"))
	}()

	first := make([]rune, 0, 40)
	for i := 0; i < 7; i++ {
		first = append(first, lex.next())
	}
	lex.backup()

	if want, have := "hello 你", string(first); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := 6, lex.pos; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	final := make([]rune, 0, 40)
	for {
		ch := lex.next()
		if ch == eof {
			break
		}
		final = append(final, ch)
	}
	if want, have := "hello 你好 世界 world", string(lex.readbuf); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "你好 世界 world", string(final); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestLexer(t *testing.T) {
	tests := []struct {
		Input string
		Items []item
	}{
		{
			"hello world",
			[]item{
				{
					typ: itemText,
					val: "hello",
				},
				{
					typ: itemSpace,
					val: " ",
				},
				{
					typ: itemText,
					val: "world",
				},
			},
		},
	}

	for _, test := range tests {
		lex := newLex(strings.NewReader(test.Input))
		i, hasEOF := 0, false
		for ; ; i++ {
			token := lex.nextItem()
			if token.typ == itemEOF {
				// do not use EOF
				hasEOF = true
				break
			}
			if i >= len(test.Items) {
				t.Errorf("more tokens than expected. got %#v", token)
				break
			}
			if want, have := test.Items[i], token; want != have {
				t.Errorf("[%d] expected %#v, got %#v", i, want, have)
			}
		}

		if want, have := true, hasEOF; want != have {
			t.Errorf("expected hasEOF=%#v; got %#v", want, have)
		}
		if want, have := len(test.Items), i; want != have {
			t.Errorf("less tokens than expected. want %#v got %#v", want, have)
			break
		}

	}
}
