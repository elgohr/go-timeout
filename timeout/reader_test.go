package timeout_test

import (
	"github.com/elgohr/go-timeout/timeout"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestTimeoutReader(t *testing.T) {
	for _, tt := range []struct {
		when   string
		given  func() (io.Reader, time.Duration)
		expect func(t *testing.T, content []byte, err error)
	}{
		{
			when: "reader is able to read everything",
			given: func() (io.Reader, time.Duration) {
				return strings.NewReader("content"), time.Second
			},
			expect: func(t *testing.T, content []byte, err error) {
				require.NoError(t, err)
				require.Equal(t, "content", string(content))
			},
		},
		{
			when: "reader times out before reader completes",
			given: func() (io.Reader, time.Duration) {
				return delayedReader{delay: 5 * time.Second, Reader: strings.NewReader("content")}, time.Millisecond
			},
			expect: func(t *testing.T, content []byte, err error) {
				require.EqualError(t, err, "timeout")
				require.Equal(t, "", string(content))
			},
		},
	} {
		t.Run(tt.when, func(t *testing.T) {
			r := timeout.Reader(tt.given())
			c, err := ioutil.ReadAll(r)
			tt.expect(t, c, err)
		})
	}
}

type delayedReader struct {
	delay time.Duration
	io.Reader
}

func (d delayedReader) Read(p []byte) (n int, err error) {
	time.Sleep(d.delay)
	return d.Reader.Read(p)
}
