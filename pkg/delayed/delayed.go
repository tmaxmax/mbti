/*
Package delayed provides a utility that is used to
print text to console with a 'typewriter' effect -
letters are printed sequentially in a delayed manner.

It uses an asynchronous API based on channels so the
caller goroutine isn't blocked.
*/
package delayed

import (
	"errors"
	"github.com/rivo/uniseg"
	"io"
	"os"
	"sync"
	"time"
)

// Properties is used to customize the behavior of the delayed utility.
type Properties struct {
	// The writer the Write operations write to. Defaults to os.Stdout.
	Writer io.StringWriter
	// The duration the Wait operations delay the execution.
	WaitDuration time.Duration
	// The duration it takes for a Write operation to execute.
	PrintDuration time.Duration
	// If true, all delays are ignored and the operations are executed instantly.
	IgnoreDelays bool
}


type delayed struct {
	properties Properties
	operations    []operation

	mu sync.Mutex
}

var defaultProperties = Properties{
	Writer: os.Stdout,
}

// New creates a delayed utility. Customize it using
// the Properties struct. Note that the underlying writer
// defaults to os.Stdout and if nil is given as a writer
// it is set back to os.Stdout.
//
//   d := New()
//
//   <-d.Write("hello", 200).
//     Wait(500).
//     Write("world!\n"). // the last explicit delay is used for subsequent operations
//     Do()
//
// It is safe for concurrent use (no more than one goroutine can access it) and
// it can be used for multiple executions.
func New(properties ...Properties) *delayed {
	props := defaultProperties
	if len(properties) > 0 {
		props = properties[0]
	}

	if props.Writer == nil {
		props.Writer = defaultProperties.Writer
	}

	return &delayed{properties: props}
}

func (d *delayed) pushWaitOperation(duration time.Duration) {
	if !d.properties.IgnoreDelays {
		d.operations = append(d.operations, &waitOperation{Duration: duration})
	}
}

func (d *delayed) pushPrintOperation(text string) {
	d.operations = append(d.operations, &writeOperation{Text: text, Writer: d.properties.Writer})
}

func getDuration(input []time.Duration, defaultDuration time.Duration) time.Duration {
	if len(input) > 0 {
		return input[0]
	}

	return defaultDuration
}

// Wait appends a wait operation for execution.
//
// The Wait operation is putting the executing goroutine to sleep for the given duration.
func (d *delayed) Wait(waitDuration ...time.Duration) *delayed {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.properties.WaitDuration = getDuration(waitDuration, d.properties.WaitDuration)
	if d.properties.WaitDuration == 0 {
		return d
	}

	d.pushWaitOperation(d.properties.WaitDuration)

	return d
}

// Write appends a print operation for execution.
//
// The Write operations is writing to the given writer each grapheme of the
// text with a delay between each other. printDuration is the duration of the
// whole print operation - the delay between each grapheme is the quotient of
// the division of the total duration with the grapheme count of the text.
func (d *delayed) Write(text string, printDuration ...time.Duration) *delayed {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.properties.PrintDuration = getDuration(printDuration, d.properties.PrintDuration)
	if d.properties.PrintDuration == 0 || d.properties.IgnoreDelays {
		d.pushPrintOperation(text)

		return d
	}

	graphemesCount := uniseg.GraphemeClusterCount(text)
	delayBetweenLetters := d.properties.PrintDuration / time.Duration(graphemesCount)
	graphemes := uniseg.NewGraphemes(text)
	appendWaitOperations := false

	for graphemes.Next() {
		if appendWaitOperations {
			d.pushWaitOperation(delayBetweenLetters)
		} else {
			appendWaitOperations = true
		}

		d.pushPrintOperation(graphemes.Str())
	}

	return d
}

// Do executes all the queued operations in a separate goroutine.
//
// Use the returned channel to wait for the execution to finish and check
// for eventual write errors.
// Use the cancel channel to stop the execution before it finishes.
func (d *delayed) Do(cancel ...chan struct{}) chan error {
	errChan := make(chan error)
	var cancelChan chan struct{}
	if len(cancel) > 0 {
		cancelChan = cancel[0]
	}

	go func() {
		d.mu.Lock()
		defer d.mu.Unlock()

		for _, op := range d.operations {
			err := op.Run(cancelChan)
			if err != nil {
				if !errors.Is(err, errCanceled) {
					errChan <- err
				}

				return
			}
		}

		d.operations = nil

		errChan <- nil
	}()

	return errChan
}

// IgnoreDelays gets or sets Properties.IgnoreDelays.
func (d *delayed) IgnoreDelays(new ...bool) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	value := d.properties.IgnoreDelays

	if len(new) > 0 {
		d.properties.IgnoreDelays = new[0]
	}

	return value
}

// Writer gets or sets Properties.Writer.
func (d *delayed) Writer(new ...io.StringWriter) io.StringWriter {
	d.mu.Lock()
	defer d.mu.Unlock()

	value := d.properties.Writer

	if len(new) > 0 && new[0] != nil {
		d.properties.Writer = new[0]
	}

	return value
}

// WaitDuration gets or sets Properties.WaitDuration.
func (d *delayed) WaitDuration(new ...time.Duration) time.Duration {
	d.mu.Lock()
	defer d.mu.Unlock()

	value := d.properties.WaitDuration

	if len(new) > 0 {
		d.properties.WaitDuration = new[0]
	}

	return value
}

// PrintDuration gets or sets Properties.PrintDuration.
func (d *delayed) PrintDuration(new ...time.Duration) time.Duration {
	d.mu.Lock()
	defer d.mu.Unlock()

	value := d.properties.PrintDuration

	if len(new) > 0 {
		d.properties.PrintDuration = new[0]
	}

	return value
}