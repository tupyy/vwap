package output

import (
	"fmt"
	"os"
	"time"

	"github.com/tupyy/vwap/internal/entity"
)

type OutputWriter struct {
	dest *os.File
}

func NewStdOutputWriter() *OutputWriter {
	return newWriter(os.Stdout)
}

func NewFileOutputWriter(f *os.File) *OutputWriter {
	return newWriter(f)
}

func newWriter(dest *os.File) *OutputWriter {
	return &OutputWriter{dest}
}

func (o *OutputWriter) Write(r entity.AverageResult) error {
	msg := fmt.Sprintf("[%s], ProductID: %s, Average: %f, Total data points: %d\n", r.Timestamp.Format(time.RFC1123Z), r.ProductID, r.Average, r.TotalPoints)
	fmt.Fprint(o.dest, msg)

	return nil
}
