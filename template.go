package main

import "strings"

// Execute renders template with whiskers like syntax https://github.com/gsf/whiskers.js
func Execute(template string, vars map[string]string) string {
	lookupVar := func(name string, vars map[string]string) string {
		name = strings.TrimSpace(name)
		for key, value := range vars {
			if strings.ToLower(key) == strings.ToLower(name) {
				return value
			}
		}
		return ""
	}

	const (
		textState = iota + 1
		varState
		escapingState
	)
	state := textState
	result := ""
	curVar := ""
	// FIXME: maybe recursive functions is better?
	for _, c := range []rune(template) {
		switch state {
		case textState:
			if c == '\\' {
				state = escapingState
			} else if c == '{' {
				curVar = ""
				state = varState
			} else {
				result += string(c)
			}
		case escapingState:
			result += string(c)
			state = textState
		case varState:
			if c == '}' {
				result += lookupVar(curVar, vars)
				state = textState
			} else {
				curVar += string(c)
			}
		}
	}
	// FIXME: check incorrect states
	return result
}
