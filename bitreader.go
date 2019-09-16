package bitio

import (
	"bytes"
	"io"
	"math/bits"
)

// BitReader implements a bitwise reader.
type BitReader struct {
	r    io.ByteReader
	buf  byte  // internal buffer
	mask uint8 // read a bit which mask covers
}

// NewBitReader returns a new BitReader. The reader changes the state of br
// internally.
func NewBitReader(r io.ByteReader) *BitReader {
	return &BitReader{
		r:    r,
		mask: 0b00000001,
	}
}

// ReadBit reads the next bit and returns it.
// At the EOF, err will be io.EOF
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

// ReadBits reads next n bits and returns them.
// It returns the number of bits successfully read.
// When didn't read any bit, bits will be nil.
func (br *BitReader) ReadBits(n int) (bits []byte, l int, err error) {
	buf := new(bytes.Buffer)
	w := NewBitWriter(buf)

	for ; l < n; l++ {
		var bit int
		bit, err = br.ReadBit()
		if err != nil {
			w.Flush()
			bits = buf.Bytes()
			return
		}
		w.WriteBit(bit)
	}

	w.Flush()
	bits = buf.Bytes()
	return
}
