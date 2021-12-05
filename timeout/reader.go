package timeout

import (
	"context"
	"errors"
	"io"
	"time"
)

func Reader(reader io.Reader, timeout time.Duration) io.Reader {
	return &timeoutReader{
		reader:  reader,
		timeout: timeout,
	}
}

type timeoutReader struct {
	reader  io.Reader
	timeout time.Duration
}

func (t *timeoutReader) Read(p []byte) (n int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	outChan := make(chan out)
	go func(p []byte) {
		n, err := t.reader.Read(p)
		outChan <- out{n: n, err: err}
		close(outChan)
	}(p)

	select {
	case <-ctx.Done():
		return 0, errors.New("timeout")
	case o := <-outChan:
		return o.n, o.err
	}
}

type out struct {
	n   int
	err error
}
