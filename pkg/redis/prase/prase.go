package prase

import (
	"errors"
	"fmt"
	"unicode"
)

func ParseArgs(str string) ([]string, error) {
	const (
		stateUnTouched = iota
		stateTouched
		stateInQuotes
		stateQuotesEnd
	)

	var quotes int32
	start := 0
	state := stateUnTouched
	args := make([]string, 0, 3)

	for i, c := range str {
		switch state {
		case stateUnTouched:
			if unicode.IsSpace(rune(c)) {
			} else if c == '\'' || c == '"' {
				state = stateInQuotes
				start = i + 1
				quotes = c
			} else {
				state = stateTouched
				start = i
			}
		case stateTouched:
			if unicode.IsSpace(rune(c)) {
				args = append(args, str[start:i])
				state = stateUnTouched
			} else if c == '\'' || c == '"' {
				return nil, errors.New(fmt.Sprintf("Unexpected: %c", c))
			} else {
			}

		case stateInQuotes:
			if c == quotes {
				args = append(args, str[start:i])
				state = stateQuotesEnd
			}
		case stateQuotesEnd:
			if !unicode.IsSpace(rune(c)) {
				return nil, errors.New(fmt.Sprintf("Unexpected: %c", c))
			}
			state = stateUnTouched
		}
	}

	switch state {
	case stateUnTouched:
	case stateTouched:
		args = append(args, str[start:])
	case stateInQuotes:
		return nil, errors.New("Unexpected: end")
	case stateQuotesEnd:
	}
	if len(args) == 0 {
		return nil, errors.New("Empty command")
	}
	return args, nil
}
