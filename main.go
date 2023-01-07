package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Writer int
type Reader struct {
	aaa *strings.Reader
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.aaa.Read(p)
}

func (*Writer) Write(p []byte) (n int, err error) {
	fmt.Println(len(p))
	return len(p), nil
}
func main() {
	r := &Reader{
		aaa: strings.NewReader("dasdlkjsakljdhaslkdhasldhaslkhdlsajkhdkl"),
	}

	w := new(Writer)
	buf := make([]byte, 3)

	// buf is used here...
	if _, err := io.CopyBuffer(w, r, buf); err != nil {
		log.Fatal(err)
	}
}
