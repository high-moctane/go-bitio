package bitio

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestWriteBit(t *testing.T) {
	type outType struct {
		b   []byte
		err error
	}

	tests := []struct {
		in  []int
		out outType
	}{
		{
			[]int{1, 0, 1, 0, 1, 0, 1, 0},
			outType{[]byte{0b10101010}, nil},
		},
		{
			[]int{1, 0, 1, 0, 1, 0, 1, 0,
				0, 0, 0, 0, 1, 1, 1, 1},
			outType{[]byte{0b10101010, 0b00001111}, nil},
		},
		{
			[]int{0, 0, 1, 0, 1, 0, 1, 1,
				0, 0, 0, 1, 1, 0, 1, 1,
				1, 0, 1, 1, 1, 1, 0, 1},
			outType{[]byte{0b00101011, 0b00011011, 0b10111101}, nil},
		},
		{
			[]int{2},
			outType{[]byte{}, &NotBitError{n: 2}},
		},
	}

testLoop:
	for idx, test := range tests {
		buf := new(bytes.Buffer)
		w := NewBitWriter(buf)

		for i, bit := range test.in {
			if err := w.WriteBit(bit); !reflect.DeepEqual(test.out.err, err) {
				t.Errorf("[%d] %dth bit unexpected error: %v", idx, i, err)
				continue testLoop
			}
		}

		got := []byte{}
		for {
			b, err := buf.ReadByte()
			if err != nil {
				break
			}
			got = append(got, b)
		}

		if !reflect.DeepEqual(test.out.b, got) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out.b, got)
		}
	}
}

func TestFlush(t *testing.T) {
	tests := []struct {
		in  []int
		out []byte
	}{
		{[]int{1, 1, 1, 1}, []byte{0b11110000}},
		{[]int{1, 1, 1, 1, 1, 1, 1}, []byte{0b11111110}},
		{[]int{1, 1, 1, 1, 1, 1, 1, 1}, []byte{0b11111111}},
		{
			[]int{1, 1, 1, 1, 1, 1, 1, 1,
				1},
			[]byte{0b11111111, 0b10000000},
		},
		{
			[]int{1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1},
			[]byte{0b11111111, 0b11110000},
		},
		{
			[]int{1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1},
			[]byte{0b11111111, 0b11111111},
		},
		{[]int{}, []byte{}},
	}

testLoop:
	for idx, test := range tests {
		buf := new(bytes.Buffer)
		w := NewBitWriter(buf)

		for _, bit := range test.in {
			if err := w.WriteBit(bit); err != nil {
				t.Errorf("[%d] caught unexpected error: %v", idx, err)
				continue testLoop
			}
		}

		// Flush is idempotent
		for i := 0; i < 10; i++ {
			if err := w.Flush(); err != nil {
				t.Errorf("[%d] caught unexpected error: %v", idx, err)
				continue testLoop
			}
		}

		got := []byte{}
		for {
			b, err := buf.ReadByte()
			if err != nil {
				break
			}
			got = append(got, b)
		}

		if !reflect.DeepEqual(test.out, got) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out, got)
		}
	}
}

func TestFlushWithOnes(t *testing.T) {
	tests := []struct {
		in  []int
		out []byte
	}{
		{[]int{1, 1, 1, 1}, []byte{0b11111111}},
		{[]int{1, 1, 1, 1, 1, 1, 1}, []byte{0b11111111}},
		{[]int{1, 1, 1, 1, 1, 1, 1, 1}, []byte{0b11111111}},
		{
			[]int{1, 1, 1, 1, 1, 1, 1, 1,
				1},
			[]byte{0b11111111, 0b11111111},
		},
		{
			[]int{1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1},
			[]byte{0b11111111, 0b11111111},
		},
		{
			[]int{1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1},
			[]byte{0b11111111, 0b11111111},
		},
		{[]int{}, []byte{}},
	}

testLoop:
	for idx, test := range tests {
		buf := new(bytes.Buffer)
		w := NewBitWriter(buf)

		for _, bit := range test.in {
			if err := w.WriteBit(bit); err != nil {
				t.Errorf("[%d] caught unexpected error: %v", idx, err)
				continue testLoop
			}
		}

		// FlushWithOnes is idempotent
		for i := 0; i < 10; i++ {
			if err := w.FlushWithOnes(); err != nil {
				t.Errorf("[%d] caught unexpected error: %v", idx, err)
				continue testLoop
			}
		}

		got := []byte{}
		for {
			b, err := buf.ReadByte()
			if err != nil {
				break
			}
			got = append(got, b)
		}

		if !reflect.DeepEqual(test.out, got) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out, got)
		}
	}
}

func TestWriteBits(t *testing.T) {
	type inType struct {
		b []byte
		n int
	}
	type outType struct {
		b   []byte
		l   int
		err error
	}

	tests := []struct {
		in  inType
		out outType
	}{
		{
			inType{nil, 0},
			outType{nil, 0, nil},
		},
		{
			inType{nil, 1},
			outType{nil, 0, io.EOF},
		},
		{
			inType{[]byte{}, 0},
			outType{nil, 0, nil},
		},
		{
			inType{[]byte{}, 1},
			outType{nil, 0, io.EOF},
		},
		{
			inType{[]byte{0b01010101}, -1},
			outType{nil, 0, nil},
		},
		{
			inType{[]byte{0b10101000}, 3},
			outType{[]byte{0b10100000}, 3, nil},
		},
		{
			inType{[]byte{0b10101001}, 8},
			outType{[]byte{0b10101001}, 8, nil},
		},
		{
			inType{[]byte{0b10101001, 0b11110011}, 8},
			outType{[]byte{0b10101001}, 8, nil},
		},
		{
			inType{[]byte{0b10101001, 0b11110011}, 12},
			outType{[]byte{0b10101001, 0b11110000}, 12, nil},
		},
		{
			inType{[]byte{0b10101001, 0b11110011, 0b10101001}, 20},
			outType{[]byte{0b10101001, 0b11110011, 0b10100000}, 20, nil},
		},
	}

	for idx, test := range tests {
		buf := new(bytes.Buffer)
		w := NewBitWriter(buf)

		l, err := w.WriteBits(test.in.b, test.in.n)
		if err != test.out.err {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out.err, err)
			continue
		}

		if test.out.l != l {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out.l, l)
		}

		if err := w.Flush(); err != nil {
			t.Errorf("[%d] unexpected error: %v", idx, err)
		}

		bits := buf.Bytes()
		if !reflect.DeepEqual(test.out.b, bits) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out.b, bits)
		}
	}
}
