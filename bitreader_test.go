package bitio

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestBitRead(t *testing.T) {
	tests := []struct {
		in  []byte
		out []int
	}{
		{
			[]byte{0b00000000},
			[]int{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			[]byte{0b10101010},
			[]int{1, 0, 1, 0, 1, 0, 1, 0},
		},
		{
			[]byte{0b10101010, 0b10101010},
			[]int{1, 0, 1, 0, 1, 0, 1, 0,
				1, 0, 1, 0, 1, 0, 1, 0},
		},
		{
			[]byte{0b11111111, 0b00000000, 0b10011001},
			[]int{1, 1, 1, 1, 1, 1, 1, 1,
				0, 0, 0, 0, 0, 0, 0, 0,
				1, 0, 0, 1, 1, 0, 0, 1},
		},
		{
			[]byte{},
			[]int{},
		},
	}

testLoop:
	for idx, test := range tests {
		r := NewBitReader(bytes.NewBuffer(test.in))
		for i := 0; i < len(test.in)*8; i++ {
			bit, err := r.ReadBit()
			if err != nil {
				t.Errorf("[%d] %dth bit unexpected error: %v", idx, i, err)
				continue testLoop
			}
			if bit != test.out[i] {
				t.Errorf("[%d] %dth bit expected %d, but %d", idx, i, test.out[i], bit)
			}
		}
		bit, err := r.ReadBit()
		if err == nil {
			t.Errorf("[%d] %dth bit not io.EOF: bit = %d, err = %v",
				idx, len(test.in)*8, bit, err)
		}
		if err != io.EOF {
			t.Errorf("[%d] %dth bit unexpected error: bit = %d, err = %v",
				idx, len(test.in)*8, bit, err)
		}
	}
}

func TestReadBits(t *testing.T) {
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
			inType{[]byte{0b00000000}, 0},
			outType{nil, 0, nil},
		},
		{
			inType{[]byte{0b10101000}, -1},
			outType{nil, 0, nil},
		},
		{
			inType{[]byte{}, 1},
			outType{nil, 0, io.EOF},
		},
		{
			inType{nil, 1},
			outType{nil, 0, io.EOF},
		},
		{
			inType{[]byte{0b10010010}, 8},
			outType{[]byte{0b10010010}, 8, nil},
		},
		{
			inType{[]byte{0b10010010}, 4},
			outType{[]byte{0b10010000}, 4, nil},
		},
		{
			inType{[]byte{0b10010010, 0b01010101}, 12},
			outType{[]byte{0b10010010, 0b01010000}, 12, nil},
		},
		{
			inType{[]byte{0b10010010, 0b01010101}, 16},
			outType{[]byte{0b10010010, 0b01010101}, 16, nil},
		},
		{
			inType{[]byte{0b10010010, 0b01010101}, 20},
			outType{[]byte{0b10010010, 0b01010101}, 16, io.EOF},
		},
	}

	for idx, test := range tests {
		r := NewBitReader(bytes.NewReader(test.in.b))
		bits, l, err := r.ReadBits(test.in.n)

		if err != test.out.err {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out.err, err)
			continue
		}

		if !reflect.DeepEqual(test.out.b, bits) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out.b, bits)
		}

		if test.out.l != l {
			t.Errorf("[%d] expected %v, but got %v", idx, test.out.l, l)
		}
	}
}
