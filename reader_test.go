package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		in       string
		template string
		out      string
		wantErr  bool
	}{
		{
			name:     "success",
			in:       "col2,col1,col3\ntwo,one,three\n2,1,3\n",
			template: `INSERT INTO foo values ({{.col1}}, {{.col2}}, {{.col3}});`,
			out: "INSERT INTO foo values (one, two, three);\n" +
				"INSERT INTO foo values (1, 2, 3);\n",
		},
		{
			name:     "undefined column is empty string",
			in:       "col2,col1,col3\ntwo,one\n",
			template: `{{.col1}} + {{.col2}} + {{.col3}} + {{.unknown}} + {{._4}}`,
			out:      "one + two +  +  + \n",
		},
		{
			name: "but in strict mode undefined column is error",
			cfg: Config{
				StrictMode: true,
			},
			in:       "col2,col1,col3\ntwo,one\n",
			template: `INSERT INTO foo values ({{.col1}}, {{.col2}}, {{.col3}})`,
			wantErr:  true,
		},
		{
			name: "in strict mode undefined number column is error",
			cfg: Config{
				StrictMode: true,
			},
			in:       "col2,col1,col3\ntwo,one\n",
			template: `{{._3}}`,
			wantErr:  true,
		},
		{
			name:     "support variables with spaces and in lower case",
			in:       "First name,Last name\nJohn,Smith\n",
			template: `{{.First_name}} {{.Last_name}} and {{.first_name}} {{.LAST_NAME}} again and link by number {{._1}} {{._2}}`,
			out:      "John Smith and John Smith again and link by number John Smith\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			magi, err := New(tt.template, tt.cfg)
			if err != nil {
				t.Fatal(err)
			}

			err = magi.ReadAndExecute(strings.NewReader(tt.in), buf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error but got %v\n", err)
			}

			got := buf.String()
			if got != tt.out {
				t.Fatalf("expected %q, but got %q", tt.out, got)
			}
		})
	}
}
