package output

import (
	"fmt"
	"os"
	"time"

	"github.com/tupyy/vwap/internal/entity"
)

type Writer struct {
	dest *os.File
}

func NewStdOutputWriter() *Writer {
	return newWriter(os.Stdout)
}

func NewFileWriter(f *os.File) *Writer {
	return newWriter(f)
}

func newWriter(dest *os.File) *Writer {
	return &Writer{dest}
}

func (o *Writer) Write(r entity.AverageResult) error {
	msg := fmt.Sprintf("[%s], ProductID: %s, Average: %f, Total data points: %d\n", r.Timestamp.Format(time.RFC1123Z), r.ProductID, r.Average, r.TotalPoints)
	fmt.Fprint(o.dest, msg)

	return nil
}
