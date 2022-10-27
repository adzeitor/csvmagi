package main

import "testing"

func TestTemplate_Execute(t *testing.T) {
	tests := []struct {
		name     string
		template string
		vars     map[string]string
		want     string
	}{
		{
			name:     "success",
			template: "{foo} is foo, but {bar} is bar",
			vars: map[string]string{
				"foo": "4",
				"bar": "5",
			},
			want: "4 is foo, but 5 is bar",
		},
		{
			name:     "lookup is case insensitive",
			template: "{foo}, {Foo}, {FOO}, {fOo}, {fOO}",
			vars: map[string]string{
				"Foo": "42",
			},
			want: "42, 42, 42, 42, 42",
		},
		{
			name:     "lookup works with reserved symbols",
			template: "{foo+1} {bar-1} {bar 3}",
			vars: map[string]string{
				"foo+1": "42",
				"bar-1": "43",
				"bar 3": "44",
			},
			want: "42 43 44",
		},
		{
			name:     `can escape { if needed with \`,
			template: `I want to render this \{definitely not var} as is`,
			vars: map[string]string{
				"definitely not var": "SHOULD NOT RENDERED",
			},
			want: `I want to render this {definitely not var} as is`,
		},
		{
			name:     `works properly with space suffixed or prefixed vars`,
			template: `{ foo} { bar } {qux }`,
			vars: map[string]string{
				"foo": "1",
				"bar": "2",
				"qux": "3",
			},
			want: `1 2 3`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Execute(tt.template, tt.vars); got != tt.want {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
