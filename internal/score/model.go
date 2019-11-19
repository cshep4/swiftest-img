package score

import "errors"

var operators = map[string]string{
	"+": "+",
	"-": "-",
	"^": "^",
	"x": "x",
	"X": "x",
	"*": "x",
	"÷": "/",
	"/": "/",
}

var ErrUnsupportedOperator = errors.New("unsupported operator")
