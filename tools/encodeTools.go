// 該檔案提供編碼轉換相關的工具函式
package tools

import (
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// Big5Decoder 會回傳一個 Big5 解碼器
func Big5Decoder() *encoding.Decoder {
	return traditionalchinese.Big5.NewDecoder()
}

// Big5ToUtf8S 會將 Big5 編碼的字串轉換為 UTF-8 編碼
func Big5ToUtf8S(s string) (string, error) {
	big5ToUtf8 := Big5Decoder()
	utf8, _, err := transform.String(big5ToUtf8, s)
	return utf8, err
}

// Big5ToUtf8B 會將 Big5 編碼的 byte slice 轉換為 UTF-8 編碼
func Big5ToUtf8B(bs []byte) ([]byte, error) {
	big5ToUtf8 := Big5Decoder()
	utf8, _, err := transform.Bytes(big5ToUtf8, bs)
	return utf8, err
}

// NewBig5ToUtf8Writer 會回傳一個 io.Writer，
// 任何寫入該 writer 的 Big5 編碼資料，都會被轉換為 UTF-8 後，再寫入傳入的 writer
func NewBig5ToUtf8Writer(w io.Writer) io.Writer {
	big5ToUtf8 := Big5Decoder()
	return transform.NewWriter(w, big5ToUtf8)
}

// NewBig5ToUtf8Reader 會回傳一個 io.Reader，
// 從該 reader 讀取資料時，會將 Big5 編碼的資料轉換為 UTF-8
func NewBig5ToUtf8Reader(r io.Reader) io.Reader {
	big5ToUtf8 := Big5Decoder()
	return transform.NewReader(r, big5ToUtf8)
}
