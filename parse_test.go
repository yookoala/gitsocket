package main

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		cmd   string
		argv  []string
		err   string
	}{
		{
			input: "hello world",
			err:   "Parse error: unknown command \"hello\"",
		},
		{
			input: "hardpull origin master",
			cmd:   "hardpull",
			argv:  []string{"origin", "master"},
		},
	}

	ctx := &gitContext{
		Src:    gitSource{"/some/dir", "upstream", "dev"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	for _, test := range tests {
		reader, writer := io.Pipe()
		go func() {
			writer.Write([]byte(test.input))
			writer.Close()
		}()

		for stmt := range parse(ctx, newLex(reader)) {

			t.Logf("input: %#v", test.input)

			// test error message against test expectation
			if stmt.err == nil && test.err != "" {
				t.Errorf("expected error: %s", test.err)
			} else if stmt.err != nil {
				if want, have := test.err, stmt.err.Error(); want != have {
					t.Errorf("expected stmt.err: %#v; got: %#v", want, have)
				}
			}

			// test cmd string
			if want, have := test.cmd, stmt.cmd; want != have {
				t.Errorf("expected stmt.cmd: %#v; got: %#v", want, have)
			}

			// test argv
			if want, have := len(test.argv), len(stmt.argv); want != have {
				t.Errorf("expected stmt.argv %#v; got %#v", stmt.argv, stmt.argv)
			} else {
				for pos, arg := range test.argv {
					if want, have := arg, stmt.argv[pos]; want != have {
						t.Errorf("expected %#v on stmt.argv[%#v], got %#v",
							want, pos, have)
					}
				}
			}
		}
	}
}
