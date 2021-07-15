package delayed

import (
	"io"
)

type writeOperation struct {
	Text string
	Writer io.StringWriter
}

func (p *writeOperation) Run(_ chan struct{}) error {
	_, err := p.Writer.WriteString(p.Text)

	return err
}
