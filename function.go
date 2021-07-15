package mbti

import (
	"errors"
	"fmt"
	"unicode"
)

const (
	focusInternal = 'I'
	focusExternal = 'E'
)

const (
	kindFeeling   = 'F'
	kindThinking  = 'T'
	kindSensation = 'S'
	kindIntuition = 'N'
)

const (
	tacticJudging     = 'J'
	tacticProspecting = 'P'
)

func invertFocus(focus rune) rune {
	if focus == focusInternal {
		return focusExternal
	}

	return focusInternal
}

func invertKind(function rune) rune {
	if function == kindFeeling {
		return kindThinking
	}

	if function == kindThinking {
		return kindFeeling
	}

	if function == kindSensation {
		return kindIntuition
	}

	return kindSensation
}

type Function struct {
	focus rune
	kind  rune
}

func (f *Function) IsJudging() bool {
	return f.kind == kindFeeling || f.kind == kindThinking
}

func (f *Function) IsProspecting() bool {
	return f.kind == kindSensation || f.kind == kindIntuition
}

func (f *Function) IsIntroverted() bool {
	return f.focus == focusInternal
}

func (f *Function) IsExtroverted() bool {
	return f.focus == focusExternal
}

func (f *Function) InvertFocus() Function {
	return Function{
		focus: invertFocus(f.focus),
		kind:  f.kind,
	}
}

func (f *Function) String() string {
	if f.IsIntroverted() {
		return string(f.kind) + string(unicode.ToLower(focusInternal))
	}

	return string(f.kind) + string(unicode.ToLower(focusExternal))
}

var ErrInvalidFunctionsString = errors.New("invalid function string")

func isValidFunction(function *Function) bool {
	return function.IsJudging() != function.IsProspecting() && function.IsIntroverted() != function.IsExtroverted()
}

func functionFromString(s string) (Function, error) {
	function := Function{
		focus: unicode.ToUpper(rune(s[1])),
		kind:  unicode.ToUpper(rune(s[0])),
	}

	if !isValidFunction(&function) {
		return Function{}, fmt.Errorf("%w %q", ErrInvalidFunctionsString, s)
	}

	return function, nil
}

func FunctionsCountInString(s string) int {
	inputLength := len(s)
	if inputLength%2 == 1 {
		return 0
	}

	count := 0
	for i := 0; i < inputLength; i += 2 {
		kindRune := unicode.ToUpper(rune(s[i]))
		focusRune := unicode.ToUpper(rune(s[i+1]))

		switch kindRune {
		case kindSensation, kindFeeling, kindIntuition, kindThinking:
		default:
			return 0
		}

		switch focusRune {
		case focusInternal, focusExternal:
		default:
			return 0
		}

		count += 1
	}

	return count
}

func FunctionsFromString(s string) ([]Function, error) {
	inputLength := len(s)
	if inputLength%2 == 1 {
		return nil, fmt.Errorf("%w: functions string must have even length", ErrInvalidFunctionsString)
	}

	funcs := make([]Function, 0, inputLength/2)
	for i := 0; i < inputLength; i += 2 {
		fn, err := functionFromString(s[i : i+2])
		if err != nil {
			return nil, err
		}

		funcs = append(funcs, fn)
	}

	return funcs, nil
}
