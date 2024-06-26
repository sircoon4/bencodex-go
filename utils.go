package bencodex

import (
	"io"
)

func writeToWriter(w io.Writer, b []byte) error {
	l := len(b)
	for t := 0; t < l; {
		n, err := w.Write(b[t:])
		if err != nil {
			return err
		}
		t += n
	}
	return nil
}

func readFromReader(r io.Reader, b *[]byte) error {
	n, err := r.Read(*b)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	} else if n < 1024 {
		*b = (*b)[:n]
	} else {
		for {
			ab := make([]byte, 64)
			n, err = r.Read(ab)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			*b = append(*b, ab[:n]...)
		}
	}
	return nil
}
