package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// itemType represents token type
type itemType int

const (
	itemError itemType = iota // lex error
	itemText                  // dot symbol
	itemSpace                 // space
	itemEOL                   // end of line
	itemEOF                   // end of connection
)

const (
	eof = -iota
)

const (
	charSmallCap = "abcdefghijklmnopqrstuvwxyz"
	charCap      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charNumber   = "0123456789"
	charAlphabet = charSmallCap + charCap
	charAlphaNum = charAlphabet + charNumber
)

// item contains information of token in a selector
type item struct {
	typ itemType
	val string
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer helps tokenize a line protocol
type lexer struct {
	input   *bufio.Reader
	state   stateFn
	readbuf []byte // buffer for backup
	pos     int    // pos of current read
	start   int    // start of current string
	width   int    // width of the last read rune
	items   chan item
}

func newLex(input io.Reader) *lexer {
	return &lexer{
		input:   bufio.NewReader(input),
		state:   lexText,
		readbuf: make([]byte, 0, 1024),
		items:   make(chan item),
	}
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// next returns the next character in the reading byte stream
func (l *lexer) next() (r rune) {

	var w int
	var err error
	fromBuf := false

	// read from rune buffer or from input
	if buflen := len(l.readbuf); buflen > l.pos {
		r, w = utf8.DecodeRune(l.readbuf[l.pos:])
		fromBuf = true
	} else {
		r, w, err = l.input.ReadRune()
	}

	// test error or end of stream
	if err == io.EOF || r == '\x00' {
		l.width = 0
		return eof
	} else if err != nil {
		panic(err)
	}

	if !fromBuf {
		l.readbuf = append(l.readbuf, []byte(string(r))...)
	}

	// append read buffer and increment width, pos
	l.width = w
	l.pos += l.width
	return r
}

// backup steps back one rune and undo 1 read. Can only undo 1 read.
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, string(l.readbuf[l.start:l.pos])}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

// lexText read all text
func lexText(l *lexer) stateFn {
	switch next := l.next(); next {
	case ' ':
		l.emit(itemSpace)
		return lexText
	case eof:
		l.emit(itemEOF)
		break
	default:
		l.acceptRun(charAlphaNum)
		l.emit(itemText)
		return lexText
	}
	return nil
}
