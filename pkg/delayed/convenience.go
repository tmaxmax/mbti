package delayed

import "time"

// Write creates a Delayed utility with a write operation queued.
func Write(format string, args ...interface{}) *Delayed {
	return New().Write(format, args...)
}

// Wait creates a Delayed utility with a wait operation queued.
// See Delayed.Write.
func Wait(duration time.Duration) *Delayed {
	return New().Wait(duration)
}

// DoWrite executes a write operation using a newly created Delayed utility.
// See Delayed.Write.
func DoWrite(format string, args ...interface{}) <-chan error {
	return Write(format, args...).Do()
}

// DoWait executes a single wait operation.
func DoWait(duration time.Duration) <-chan error {
	e := make(chan error)
	op := waitOperation{Duration: duration}

	go func() {
		e <- op.Run(nil)
	}()

	return e
}
