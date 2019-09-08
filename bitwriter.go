package bitio

import (
	"fmt"
	"io"
	"math/bits"
)

// NotBitError will be used when calling WriteBit with an incorrect argument.
type NotBitError struct {
	n int
}

func newNotBitError(n int) *NotBitError {
	return &NotBitError{n: n}
}

// Error returns a NotBitError message.
func (e *NotBitError) Error() string {
	return fmt.Sprintf("%d is not a bit", e.n)
}

// BitWriter implements a bitwise writer. The client should call the Flush or
// FlushWithOnes method when all data have been written.
type BitWriter struct {
	w    io.ByteWriter
	buf  byte
	mask uint8
}

// NewBitWriter returns a new BifWriter. The writer changes the state of w
// internally.
func NewBitWriter(w io.ByteWriter) *BitWriter {
	return &BitWriter{
		w:    w,
		mask: 0b10000000,
	}
}

// WriteBit writes bit. The err will be not nil when io.ByteWriter returns
// non-nil error.
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

// Flush writes down an internal buffer with zeros. The client should call this
// method when the all data have been written.
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

// FlushWithOnes writes down an internal buffer with ones. The client should
//call this method when the all data have been written.
func (bw *BitWriter) FlushWithOnes() error {
	for bw.mask != 0b10000000 {
		if err := bw.WriteBit(1); err != nil {
			return err
		}
	}

	return nil
}
