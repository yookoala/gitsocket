package main

import "fmt"

// statement represents executable command statement
// or contain any parse / syntax error
type statement struct {
	ctx  *gitContext
	cmd  string
	argv []string
	err  error
}

func (c *statement) run() error {
	switch c.cmd {
	case "hardpull":
		return c.ctx.HardPull()
	}
	return fmt.Errorf("Statement: unknown command %#v", c.cmd)
}

// parse parse the lexer into commands to run
func parse(ctx *gitContext, l *lexer) (cout <-chan *statement) {
	ch := make(chan *statement)
	go func() {
		defer close(ch)

		tokens, tokenStrs := []item{}, make([]string, 0)
		for token := l.nextItem(); token.typ != itemEOF; token = l.nextItem() {
			tokens = append(tokens, token)
			if token.typ == itemText {
				tokenStrs = append(tokenStrs, token.val)
			}
		}
		switch cmd := tokenStrs[0]; cmd {
		case "hardpull":
			ch <- &statement{
				ctx:  ctx,
				cmd:  cmd,
				argv: tokenStrs[1:],
			}
		default:
			ch <- &statement{
				ctx: ctx,
				err: fmt.Errorf("Parse error: unknown command %#v", cmd),
			}
		}
	}()
	return ch
}
