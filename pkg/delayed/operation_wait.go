package delayed

import (
	"time"
)

type waitOperation struct {
	Duration time.Duration
}


func (w *waitOperation) Run(cancel chan struct{}) error {
	select {
	case <-cancel:
		return errCanceled
	case <-time.After(w.Duration):
		return nil
	}
}