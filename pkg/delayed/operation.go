package delayed

import "errors"

// operation is a generic interface for
// tasks executed by the Delayed utility.
type operation interface {
	Run(cancel <-chan struct{}) error
}

var errCanceled = errors.New("operation canceled")
