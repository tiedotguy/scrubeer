package binaryreader

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BinaryReader struct {
	// LargeID is a flag to indicate if IDs are 32bit or 64bit
	LargeID bool
	Trace   bool

	buf []byte
}

func NewBinaryReader(b []byte) *BinaryReader {
	return &BinaryReader{
		buf: b,
	}
}

func (br *BinaryReader) StringNul() string {
	s, b, ok := bytes.Cut(br.buf, []byte{0})
	if !ok {
		panic("no nul found")
	}
	br.buf = b
	return string(s)
}

func (br *BinaryReader) U8() uint8 {
	u8 := br.buf[0]
	br.buf = br.buf[1:]
	if br.Trace {
		fmt.Printf("Read %02x\n", u8)
	}
	return u8
}

func (br *BinaryReader) U8AsInt() int {
	return int(br.U8())
}

func (br *BinaryReader) U16() uint16 {
	u16 := binary.BigEndian.Uint16(br.buf[0:2])
	br.buf = br.buf[2:]
	if br.Trace {
		fmt.Printf("Read %04x\n", u16)
	}
	return u16
}

func (br *BinaryReader) U32() uint32 {
	u32 := binary.BigEndian.Uint32(br.buf[0:4])
	br.buf = br.buf[4:]
	if br.Trace {
		fmt.Printf("Read %08x\n", u32)
	}
	return u32
}

func (br *BinaryReader) U64() uint64 {
	u64 := binary.BigEndian.Uint64(br.buf[0:8])
	br.buf = br.buf[8:]
	if br.Trace {
		fmt.Printf("Read %016x\n", u64)
	}
	return u64
}

func (br *BinaryReader) ID() uint64 {
	if br.LargeID {
		return br.U64()
	}
	return uint64(br.U32())
}

func (br *BinaryReader) Clone() *BinaryReader {
	brClone := *br
	return &brClone
}

func (br *BinaryReader) Dump(length int) {
	fmt.Printf("%02x\n", br.buf[:length])
}

func (br *BinaryReader) EOF() bool {
	return len(br.buf) == 0
}

func (br *BinaryReader) Bytes(count int) []byte {
	r := br.buf[:count]
	br.buf = br.buf[count:]
	return r
}

func (br *BinaryReader) Size(dataType int) int {
	switch dataType {
	case OHT_OBJECT: // ID
		if br.LargeID {
			return 8
		} else {
			return 4
		}
	case OHT_BOOL, OHT_BYTE:
		return 1
	case OHT_CHAR, OHT_SHORT:
		return 2
	case OHT_FLOAT, OHT_INT:
		return 4
	case OHT_DOUBLE, OHT_LONG:
		return 8
	}
	panic(fmt.Sprintf("unrecognized type: %d", dataType))
}

func (br *BinaryReader) Slice(length int) *BinaryReader {
	brSlice := &BinaryReader{
		LargeID: br.LargeID,
		buf:     br.buf[:length],
	}
	br.buf = br.buf[length:]
	return brSlice
}

func (br *BinaryReader) Length() int {
	return len(br.buf)
}

const (
	OHT_OBJECT  = 2
	OHT_BOOL    = 4
	OHT_CHAR    = 5
	OHT_FLOAT   = 6
	OHT_DOUBLE  = 7
	OHT_BYTE    = 8
	OHT_SHORT   = 9
	OHT_INT     = 10
	OHT_LONG    = 11
	OHT_INVALID = 255
)
