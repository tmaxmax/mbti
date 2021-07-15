package mbti

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Personality struct {
	primary   Function
	auxiliary Function
	tertiary  Function
	inferior  Function
}

func (p *Personality) Unconscious() *Personality {
	return &Personality{
		primary:   p.primary.InvertFocus(),
		auxiliary: p.auxiliary.InvertFocus(),
		tertiary:  p.tertiary.InvertFocus(),
		inferior:  p.inferior.InvertFocus(),
	}
}

func (p *Personality) Subconscious() *Personality {
	return &Personality{
		primary:   p.inferior,
		auxiliary: p.tertiary,
		tertiary:  p.auxiliary,
		inferior:  p.primary,
	}
}

func (p *Personality) SuperEgo() *Personality {
	return &Personality{
		primary:   p.inferior.InvertFocus(),
		auxiliary: p.tertiary.InvertFocus(),
		tertiary:  p.auxiliary.InvertFocus(),
		inferior:  p.primary.InvertFocus(),
	}
}

func (p *Personality) String() string {
	ret := &strings.Builder{}

	var extrovertedFunction Function

	if p.primary.IsExtroverted() {
		ret.WriteRune(focusExternal)

		extrovertedFunction = p.primary
	} else {
		ret.WriteRune(focusInternal)

		extrovertedFunction = p.auxiliary
	}

	if p.primary.IsProspecting() {
		ret.WriteRune(p.primary.kind)
		ret.WriteRune(p.auxiliary.kind)
	} else {
		ret.WriteRune(p.auxiliary.kind)
		ret.WriteRune(p.primary.kind)
	}

	if extrovertedFunction.IsProspecting() {
		ret.WriteRune(tacticProspecting)
	} else {
		ret.WriteRune(tacticJudging)
	}

	return ret.String()
}

func (p *Personality) Functions() []Function {
	return []Function{p.primary, p.auxiliary, p.tertiary, p.inferior}
}

var ErrInvalidFunctions = errors.New("invalid functions")

func fromValidDominantFunctions(primary, auxiliary Function) *Personality {
	return &Personality{
		primary:   primary,
		auxiliary: auxiliary,
		tertiary: Function{
			focus: primary.focus,
			kind:  invertKind(auxiliary.kind),
		},
		inferior: Function{
			focus: auxiliary.focus,
			kind:  invertKind(primary.kind),
		},
	}
}

func FromDominantFunctions(primary, auxiliary Function) (*Personality, error) {
	if !isValidFunction(&primary) || !isValidFunction(&auxiliary) {
		return nil, ErrInvalidFunctions
	}

	if primary.focus == auxiliary.focus || primary.IsJudging() == auxiliary.IsJudging() {
		return nil, fmt.Errorf("%w: primary %q and auxiliary %q can't form a personality type", ErrInvalidFunctions, primary.String(), auxiliary.String())
	}

	return fromValidDominantFunctions(primary, auxiliary), nil
}

var ErrInvalidIndicatorString = errors.New("invalid indicator string")

func getIndicatorRunes(indicator string) (rune, rune, rune, rune, error) {
	if strings.HasSuffix(indicator, "-A") || strings.HasSuffix(indicator, "-T") {
		indicator = indicator[:len(indicator)-2]
	}

	if len(indicator) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("%w %q", ErrInvalidIndicatorString, indicator)
	}

	focusRune := unicode.ToUpper(rune(indicator[0]))
	perceivingRune := unicode.ToUpper(rune(indicator[1]))
	judgingRune := unicode.ToUpper(rune(indicator[2]))
	tacticsRune := unicode.ToUpper(rune(indicator[3]))

	return focusRune, perceivingRune, judgingRune, tacticsRune, nil
}

func IsIndicatorString(indicator string) bool {
	focusRune, perceivingRune, judgingRune, tacticsRune, err := getIndicatorRunes(indicator)

	return err == nil &&
		(focusRune == focusInternal || focusRune == focusExternal) &&
		(perceivingRune == kindIntuition || perceivingRune == kindSensation) &&
		(judgingRune == kindThinking || judgingRune == kindFeeling) &&
		(tacticsRune == tacticJudging || tacticsRune == tacticProspecting)
}

func FromIndicator(indicator string) (*Personality, error) {
	focusRune, perceivingRune, judgingRune, tacticsRune, err := getIndicatorRunes(indicator)
	if err != nil {
		return nil, err
	}

	primary, auxiliary := Function{focus: focusRune}, Function{}

	if primary.IsIntroverted() == primary.IsExtroverted() {
		return nil, fmt.Errorf("%w: %q is not a valid focus letter", ErrInvalidIndicatorString, string(focusRune))
	}

	auxiliary.focus = invertFocus(primary.focus)

	if tacticsRune == tacticJudging {
		if primary.IsExtroverted() {
			primary.kind = judgingRune
			auxiliary.kind = perceivingRune
		} else {
			auxiliary.kind = judgingRune
			primary.kind = perceivingRune
		}
	} else if tacticsRune == tacticProspecting {
		if primary.IsExtroverted() {
			primary.kind = perceivingRune
			auxiliary.kind = judgingRune
		} else {
			auxiliary.kind = perceivingRune
			primary.kind = judgingRune
		}
	}

	return FromDominantFunctions(primary, auxiliary)
}
