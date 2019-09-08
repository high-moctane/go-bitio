package bitio

import (
	"fmt"
	"io"
	"math/bits"
)

type NotBitError struct {
	n int
}

func newNotBitError(n int) *NotBitError {
	return &NotBitError{n: n}
}

func (e *NotBitError) Error() string {
	return fmt.Sprintf("%d is not a bit", e.n)
}

type BitWriter struct {
	w    io.ByteWriter
	buf  byte
	mask uint8
}

func NewBitWriter(w io.ByteWriter) *BitWriter {
	return &BitWriter{
		w:    w,
		mask: 0b10000000,
	}
}

func (bw *BitWriter) WriteBit(bit int) error {
	if bit != 0 && bit != 1 {
		return newNotBitError(bit)
	}

	if bit == 1 {
		bw.buf |= bw.mask
	}

	if bw.mask == 0b00000001 {
		if err := bw.w.WriteByte(bw.buf); err != nil {
			return err
		}
		bw.buf = 0b00000000
	}

	bw.mask = bits.RotateLeft8(bw.mask, -1)

	return nil
}

func (bw *BitWriter) Flush() error {
	if bw.mask == 0b10000000 {
		return nil
	}

	if err := bw.w.WriteByte(bw.buf); err != nil {
		return err
	}

	bw.buf = 0b00000000
	bw.mask = 0b10000000
	return nil
}

func (bw *BitWriter) FlushWithOnes() error {
	for bw.mask != 0b10000000 {
		if err := bw.WriteBit(1); err != nil {
			return err
		}
	}

	return nil
}
