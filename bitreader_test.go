package bitio

import (
	"bytes"
	"io"
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

	for idx, test := range tests {
		r := NewBitReader(bytes.NewBuffer(test.in))
		for i := 0; i < len(test.in)*8; i++ {
			bit, err := r.ReadBit()
			if err != nil {
				t.Errorf("[%d] %dth bit unexpected error: %v", idx, i, err)
				continue
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
