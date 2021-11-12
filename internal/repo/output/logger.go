package output

import (
	"fmt"
	"os"
	"time"

	"github.com/tupyy/vwap/internal/entity"
)

type OutputWrite struct {
	dest *os.File
}

func NewStdOutputWriter() *OutputWrite {
	return newWriter(os.Stdout)
}

func NewFileOutputWriter(f *os.File) *OutputWrite {
	return newWriter(f)
}

func newWriter(dest *os.File) *OutputWrite {
	return &OutputWrite{dest}
}

func (o *OutputWrite) Write(r entity.AverageResult) error {
	msg := fmt.Sprintf("[%s], ProductID: %s, Average: %f, Total data points: %d", r.Timestamp.Format(time.RFC1123Z), r.ProductID, r.Average, r.TotalPoints)

	if _, err := o.dest.WriteString(msg); err != nil {
		return err
	}

	return nil
}
