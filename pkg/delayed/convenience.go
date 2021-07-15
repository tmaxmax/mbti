package delayed

import "time"

// Write creates a delayed utility with a write operation queued.
func Write(text string, duration time.Duration) *delayed {
	return New().Write(text, duration)
}

// Wait creates a delayed utility with a wait operation queued.
func Wait(duration time.Duration) *delayed {
	return New().Wait(duration)
}

// DoWrite executes a write operation using a newly created delayed utility.
func DoWrite(text string, duration time.Duration) chan error {
	return Write(text, duration).Do()
}

// DoWait executes a single wait operation.
func DoWait(duration time.Duration) chan error {
	e := make(chan error)
	op := waitOperation{Duration: duration}

	go func() {
		e <- op.Run(nil)
	}()

	return e
}