package misc

import (
	"io"

	"github.com/hashicorp/go-multierror"
)

type MultiCloser struct {
	closers []io.Closer
}

func NewMultiCloser(closers []io.Closer) *MultiCloser {
	return &MultiCloser{
		closers: closers,
	}
}

func (m *MultiCloser) Close() error {

	var err error

	for _, c := range m.closers {
		if e := c.Close(); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return err
}

func MultiReaderCreate(inReader io.Reader, outReadersCount int, errch chan error) []io.Reader {

	// Compose readers for return
	readers, writer, closer := getReaders(outReadersCount)

	// Start writes
	go func() {

		// Copy input reader to all readers
		_, err := io.Copy(writer, inReader)

		closer.Close()

		// Send error
		errch <- err
	}()

	return readers
}

func getReaders(count int) ([]io.Reader, io.Writer, io.Closer) {

	readers := make([]io.Reader, 0, count)
	writers := make([]io.Writer, 0, count)
	closers := make([]io.Closer, 0, count)

	for i := 0; i < count; i++ {

		r, w := io.Pipe()

		readers = append(readers, r)
		writers = append(writers, w)
		closers = append(closers, w)
	}

	return readers, io.MultiWriter(writers...), NewMultiCloser(closers)
}
