package tools

import (
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

func Big5Decoder() *encoding.Decoder {
	return traditionalchinese.Big5.NewDecoder()
}

func Big5ToUtf8S(s string) (string, error) {
	big5ToUtf8 := Big5Decoder()
	utf8, _, err := transform.String(big5ToUtf8, s)
	return utf8, err
}

func Big5ToUtf8B(bs []byte) ([]byte, error) {
	big5ToUtf8 := Big5Decoder()
	utf8, _, err := transform.Bytes(big5ToUtf8, bs)
	return utf8, err
}

func NewBig5ToUtf8Writer(w io.Writer) io.Writer {
	big5ToUtf8 := Big5Decoder()
	return transform.NewWriter(w, big5ToUtf8)
}

func NewBig5ToUtf8Reader(r io.Reader) io.Reader {
	big5ToUtf8 := Big5Decoder()
	return transform.NewReader(r, big5ToUtf8)
}
