package binarywriter

import (
	"encoding/binary"
)

type BinaryWriter struct {
	// LargeID is a flag to indicate if IDs are 32bit or 64bit
	LargeID bool

	buf []byte
}

func NewBinaryWriter(b []byte) *BinaryWriter {
	return &BinaryWriter{
		buf: b,
	}
}

func (bw *BinaryWriter) StringNul(s string) string {
	copy(bw.buf[:len(s)], s)
	bw.buf = bw.buf[len(s):]
	bw.U8(0)
	return s
}

func (bw *BinaryWriter) U8(u8 uint8) uint8 {
	bw.buf[0] = u8
	bw.buf = bw.buf[1:]
	return u8
}

func (bw *BinaryWriter) U8AsInt(u8 uint8) int {
	return int(bw.U8(u8))
}

func (bw *BinaryWriter) U16(u16 uint16) uint16 {
	binary.BigEndian.PutUint16(bw.buf[0:2], u16)
	bw.buf = bw.buf[2:]
	return u16
}

func (bw *BinaryWriter) U16AsInt(u16 uint16) int {
	return int(bw.U16(u16))
}

func (bw *BinaryWriter) U32(u32 uint32) uint32 {
	binary.BigEndian.PutUint32(bw.buf[0:4], u32)
	bw.buf = bw.buf[4:]
	return u32
}

func (bw *BinaryWriter) U32AsInt(u32 uint32) int {
	return int(bw.U32(u32))
}

func (bw *BinaryWriter) U64(u64 uint64) uint64 {
	binary.BigEndian.PutUint64(bw.buf[0:8], u64)
	bw.buf = bw.buf[8:]
	return u64
}

func (bw *BinaryWriter) ID(id uint64) uint64 {
	if bw.LargeID {
		return uint64(bw.U32(uint32(id)))
	} else {
		return bw.U64(id)
	}
}

func (bw *BinaryWriter) Bytes(data []byte) {
	copy(bw.buf[:len(data)], data)
	bw.buf = bw.buf[len(data):]
}
