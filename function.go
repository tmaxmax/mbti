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
	KindFeeling   = 'F'
	KindThinking  = 'T'
	KindSensation = 'S'
	KindIntuition = 'N'
)

const (
	tacticJudging     = 'J'
	tacticProspecting = 'P'
)

func invertFocus(focus rune) rune {
	switch focus {
	case focusInternal:
		return focusExternal
	default:
		return focusInternal
	}
}

func invertKind(function rune) rune {
	switch function {
	case KindFeeling:
		return KindThinking
	case KindThinking:
		return KindFeeling
	case KindSensation:
		return KindIntuition
	default:
		return KindSensation
	}
}

type Function struct {
	focus rune
	kind  rune
}

func (f Function) IsJudging() bool {
	return f.kind == KindFeeling || f.kind == KindThinking
}

func (f Function) IsProspecting() bool {
	return f.kind == KindSensation || f.kind == KindIntuition
}

func (f Function) IsIntroverted() bool {
	return f.focus == focusInternal
}

func (f Function) IsExtroverted() bool {
	return f.focus == focusExternal
}

func (f Function) invertFocus() Function {
	return Function{
		focus: invertFocus(f.focus),
		kind:  f.kind,
	}
}

func (f Function) String() string {
	if f.IsIntroverted() {
		return string(f.kind) + string(unicode.ToLower(focusInternal))
	}

	return string(f.kind) + string(unicode.ToLower(focusExternal))
}

func (f Function) Kind() rune {
	return f.kind
}

var ErrInvalidFunctionsString = errors.New("invalid function string")

func isValidFunction(function Function) bool {
	return function.IsJudging() != function.IsProspecting() && function.IsIntroverted() != function.IsExtroverted()
}

func functionFromStringUnchecked(s string) Function {
	return Function{
		focus: unicode.ToUpper(rune(s[1])),
		kind:  unicode.ToUpper(rune(s[0])),
	}
}

func functionFromString(s string) (Function, error) {
	function := functionFromStringUnchecked(s)
	if !isValidFunction(function) {
		return Function{}, fmt.Errorf("%w %q", ErrInvalidFunctionsString, s)
	}

	return function, nil
}

func FunctionCountInString(s string) int {
	inputLength := len(s)
	if inputLength%2 == 1 {
		return 0
	}

	count := 0
	for i := 0; i < inputLength; i += 2 {
		if fn := functionFromStringUnchecked(s[i : i+2]); !isValidFunction(fn) {
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
