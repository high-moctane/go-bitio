# go-bitio

[![TravisCI](https://travis-ci.org/high-moctane/go-bitio.svg?branch=master)](https://travis-ci.org/high-moctane/go-bitio)
[![GolangCI](https://golangci.com/badges/github.com/high-moctane/go-bitio.svg)](https://golangci.com/r/github.com/high-moctane/go-bitio)
[![CodeCov](https://codecov.io/gh/high-moctane/go-bitio/branch/master/graph/badge.svg)](https://codecov.io/gh/high-moctane/go-bitio)
[![GoDoc](https://godoc.org/github.com/high-moctane/go-bitio?status.svg)](https://godoc.org/github.com/high-moctane/go-bitio)

Package `bitio` implements a simple bitwise reader and writer.

The client can use `BitReader.ReadBit()` to read a bit from an `io.ByteReader`,
and also use `BitWriter.WriteBit()` to write a bit to an `io.ByteWriter`.
After write with `BitWriter.WriteBit()`, you must call `BitWriter.Flush()` or
`BitWriter.FlushWithOnes()`.

## Examples

### BitReader

```go
r := NewBitReader(bytes.NewBuffer([]byte{0b10110100}))
bit, err := r.ReadBit()	// 1, nil
bit, err := r.ReadBit()	// 0, nil
bit, err := r.ReadBit()	// 1, nil
bit, err := r.ReadBit()	// 1, nil
bit, err := r.ReadBit()	// 0, nil
bit, err := r.ReadBit()	// 1, nil
bit, err := r.ReadBit()	// 0, nil
bit, err := r.ReadBit()	// 0, nil
bit, err := r.ReadBit()	// 0, io.EOF
```

### BitWriter

```go
buf := new(bytes.Buffer)
w := NewBitWriter(buf)
defer w.Flush()

err := w.WriteBit(1)
err := w.WriteBit(0)
err := w.WriteBit(1)
err := w.WriteBit(1)
err := w.WriteBit(0)
err := w.WriteBit(1)

byte, err := buf.ReadByte()	// 0b10110100, nil
```
