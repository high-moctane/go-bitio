package bitio

import (
	"io"
	"math/bits"
)

type BitReader struct {
	r    io.ByteReader
	buf  byte
	mask uint8
}

func NewBitReader(br io.ByteReader) *BitReader {
	return &BitReader{
		r:    br,
		mask: 0b00000001,
	}
}

func (br *BitReader) ReadBit() (bit int, err error) {
	br.mask = bits.RotateLeft8(br.mask, -1)

	if br.mask == 0b10000000 {
		br.buf, err = br.r.ReadByte()
		if err != nil {
			return
		}
	}

	return bits.OnesCount8(br.buf & br.mask), nil
}
