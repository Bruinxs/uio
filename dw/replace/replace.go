package replace

import "io"

type Replacer interface {
	Replace(year int, month int, day int) (io.WriteCloser, error)
}
