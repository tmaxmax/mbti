package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tmaxmax/mbti"

	"github.com/tmaxmax/mbti/pkg/delayed"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	instantOutput := flag.Bool("instantOutput", false, "True if you want output to be shown instantly, without a typewriter-like effect")

	flag.Parse()

	d := delayed.New(delayed.Properties{
		IgnoreDelays:  *instantOutput,
		PrintDuration: time.Second,
		WaitDuration:  time.Second / 2,
	})

	for {
		<-d.Write("Input dominant functions (e.g. FeNi) or a Myers-Briggs type indicator, or type \"exit\" to close the program.\n").
			Write("-> ", time.Duration(0)).
			Do()

		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Input error:", err)
			exitCode = 1

			return
		}

		if input == "exit" {
			break
		}

		ego, err := personalityFromInput(input)
		if err != nil {
			log.Printf("%s\n\n", err)

			continue
		}

		unconscious := ego.Unconscious()
		subconscious := ego.Subconscious()
		superEgo := ego.SuperEgo()

		<-d.Write("Ego: %s (%s)\n", ego, formatFunctions(ego.Functions()), time.Second).Wait().
			Write("Unconscious: %s (%s)\n", unconscious, formatFunctions(unconscious.Functions())).Wait().
			Write("Subconscious: %s (%s)\n", subconscious, formatFunctions(subconscious.Functions())).Wait().
			Write("Super-ego: %s (%s)\n\n", superEgo, formatFunctions(superEgo.Functions())).Wait().
			Do()
	}
}

func personalityFromInput(input string) (*mbti.Personality, error) {
	if mbti.FunctionCountInString(input) == 2 {
		functions, _ := mbti.FunctionsFromString(input)

		return mbti.FromDominantFunctions(functions[0], functions[1])
	} else if mbti.IsIndicatorString(input) {
		return mbti.FromIndicator(input)
	}

	return nil, errors.New("invalid input")
}

func formatFunctions(functions []mbti.Function) string {
	representations := make([]string, 0, len(functions))

	for _, fn := range functions {
		representations = append(representations, fn.String())
	}

	return strings.Join(representations, " ")
}
