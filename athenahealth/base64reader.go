package athenahealth

import (
	"encoding/base64"
	"io"
)

type base64Reader struct {
	reader io.Reader
}

func newBase64Reader(reader io.Reader) *base64Reader {
	return &base64Reader{
		reader: reader,
	}
}

func (br *base64Reader) Read(p []byte) (int, error) {
	buffer := make([]byte, base64.StdEncoding.DecodedLen(len(p)))

	read, err := br.reader.Read(buffer)
	if err != nil {
		return 0, err
	}

	encodedLen := base64.StdEncoding.EncodedLen(read)
	if encodedLen > len(p) {
		return 0, io.ErrShortBuffer
	}

	encBuf := make([]byte, encodedLen)
	base64.StdEncoding.Encode(encBuf, buffer[:read])

	copy(p, encBuf[:encodedLen])

	return encodedLen, nil
}
