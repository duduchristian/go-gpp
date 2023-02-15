package util

import (
	"encoding/base64"
	"fmt"
)

var (
	fibEncodeNumOutOfRangeErr = fmt.Errorf("the number to be encoded is out of range")
)

func (bs *BitStream) enlarge(n uint16) {
	final := (bs.p + n + 7) / 8
	for final > bs.Len() {
		bs.b = append(bs.b, 0x00)
	}
}

func (bs *BitStream) appendNBits(b []byte, n uint16) {
	bs.enlarge(n)
	var (
		i         uint16
		byteIndex uint16
		offset    uint16
	)
	for i = 0; i < n; i++ {
		byteIndex = bs.p / 8
		offset = bs.p % 8
		bs.b[byteIndex] |= (b[i/8] << (i % 8)) & 0x80 >> offset
		bs.p++
	}
}

func (bs *BitStream) Base64Encode() []byte {
	encoded := make([]byte, base64.RawURLEncoding.EncodedLen(len(bs.b)))
	base64.RawURLEncoding.Encode(encoded, bs.b)
	return encoded
}

func (bs *BitStream) Reset() {
	bs.b = nil
	bs.p = 0
}

func (bs *BitStream) WriteByte1(b byte) {
	bs.appendNBits([]byte{b << 7}, 1)
}

func (bs *BitStream) WriteByte2(b byte) {
	bs.appendNBits([]byte{b << 6}, 2)
}

func (bs *BitStream) WriteByte4(b byte) {
	bs.appendNBits([]byte{b << 4}, 4)
}

func (bs *BitStream) WriteByte6(b byte) {
	bs.appendNBits([]byte{b << 2}, 6)
}

func (bs *BitStream) WriteByte8(b byte) {
	bs.appendNBits([]byte{b}, 8)
}

func (bs *BitStream) WriteUInt12(b uint16) {
	first := byte(b >> 4)
	second := byte(b << 4)
	bs.appendNBits([]byte{first, second}, 12)
}

func (bs *BitStream) WriteUInt16(b uint16) {
	first := byte(b >> 8)
	second := byte(b)
	bs.appendNBits([]byte{first, second}, 16)
}

func (bs *BitStream) WriteTwoBitField(bList []byte) {
	for _, b := range bList {
		bs.WriteByte2(b)
	}
}

func (bs *BitStream) WriteFibonacciInt(num uint16) error {
	// The num should be [1,6765). Actually once the num is larger than or equal to 987,
	// the efficiency of Fibonacci Encoding would be no better than 'WriteUint16'.
	if num <= 0 || num >= fibonacci(fibLen) {
		return fibEncodeNumOutOfRangeErr
	}
	// Binary Search to find the largest fibonacci number less than or equal to num.
	lo, hi := 2, fibLen-1
	for lo < hi {
		mid := (lo + hi) / 2
		if num >= fibonacci(mid) && num < fibonacci(mid+1) {
			lo = mid
			break
		}
		if num < fibonacci(mid) {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	// Calculate the fibonacci-encoded sequence.
	var fibEncoded uint32 = 1
	offset := 1
	for i := lo; i >= 2; i-- {
		if num >= fibLookup[i] {
			num -= fibLookup[i]
			fibEncoded |= 1 << offset
		}
		offset++
	}
	encodedLength := lo
	fibEncoded <<= 32 - encodedLength

	bs.appendNBits([]byte{
		byte(fibEncoded >> 24),
		byte(fibEncoded >> 16),
		byte(fibEncoded >> 8),
		byte(fibEncoded),
	}, uint16(encodedLength))
	return nil
}

func (bs *BitStream) WriteIntRange(intRange *IntRange) error {
	var err error
	bs.WriteUInt12(intRange.Size)
	// Assume that the ranges are ordered.
	for _, r := range intRange.Range {
		if r.StartID == r.EndID {
			bs.WriteByte1(0)
			err = bs.WriteFibonacciInt(r.StartID)
			if err != nil {
				return fmt.Errorf("write int range error: %v", err)
			}
		} else {
			bs.WriteByte1(1)
			err = bs.WriteFibonacciInt(r.StartID)
			if err != nil {
				return fmt.Errorf("write int range error: %v", err)
			}
			err = bs.WriteFibonacciInt(r.EndID)
			if err != nil {
				return fmt.Errorf("write int range error: %v", err)
			}
		}
	}
	return nil
}