package main

import (
	"bytes"
	"fmt"
	"os"

	"github.org/tiedotguy/scrubeer/internal/args"
	"github.org/tiedotguy/scrubeer/internal/binaryreader"
	"github.org/tiedotguy/scrubeer/internal/binarywriter"
	"github.org/tiedotguy/scrubeer/internal/io"
)

func main() {
	opts := args.ParseArgs[Opts](os.Args[1:])

	fIn := io.Open(opts.InputFile)
	mmIn := io.MapRO(fIn)
	fOut := io.Create(opts.OutputFile, int64(len(mmIn)))
	mmOut := io.MapRW(fOut)

	br := binaryreader.NewBinaryReader(mmIn)
	bw := binarywriter.NewBinaryWriter(mmOut)

	s := bw.StringNul(br.StringNul())
	if s != "JAVA PROFILE 1.0.1" && s != "JAVA PROFILE 1.0.2" {
		panic("unrecognized header: " + s)
	}

	idSize := bw.U32(br.U32())
	switch idSize {
	case 4:
		br.LargeID = false
	case 8:
		br.LargeID = true
	default:
		panic(fmt.Sprintf("invalid id size: %d", idSize))
	}
	bw.U64(br.U64()) // timestamp

	brIterator := br.Clone()
	for !brIterator.EOF() {
		tagType := bw.U8(brIterator.U8())
		bw.U32(brIterator.U32()) // Timestamp
		tagLength := bw.U32AsInt(brIterator.U32())
		brTag := brIterator.Slice(tagLength)
		switch tagType {
		case TagString: // This can be copied
			bw.Bytes(brTag.Bytes(brTag.Length()))
		case TagLoadClass: // This can be copied
			bw.Bytes(brTag.Bytes(brTag.Length()))
		case TagStackFrame: // This can be copied
			bw.Bytes(brTag.Bytes(brTag.Length()))
		case TagStackTrace: // This can be copied
			bw.Bytes(brTag.Bytes(brTag.Length()))
		case TagHeapDump, TagHeapDumpSegment:
			for !brTag.EOF() {
				objType := bw.U8(brTag.U8())
				switch objType {
				case SubTagRootJNIGlobal: // 1 / 0x01
					bw.ID(brTag.ID()) // Object ID
					bw.ID(brTag.ID()) // JNI global ref ID
				case SubTagRootJNILocal: // 2 / 0x02
					bw.ID(brTag.ID())   // Object ID
					bw.U32(brTag.U32()) // Thread serial
					bw.U32(brTag.U32()) // Thread frame index
				case SubTagRootJavaFrame: // 3 / 0x03
					bw.ID(brTag.ID())   // Object ID
					bw.U32(brTag.U32()) // Thread serial
					bw.U32(brTag.U32()) // Thread frame index
				case SubTagRootStickyClass: // 5 / 0x05
					bw.ID(brTag.ID()) // Object ID
				case SubTagRootMonitorUsed: // 7 / 0x07
					bw.ID(brTag.ID()) // Thread Object ID
				case SubTagRootThreadObject: // 8 / 0x08
					bw.ID(brTag.ID())   // Thread Object ID
					bw.U32(brTag.U32()) // Thread serial
					bw.U32(brTag.U32()) // Stack serial
				case SubTagClassDump: // 32 / 0x20
					bw.ID(brTag.ID())   // Object ID
					bw.U32(brTag.U32()) // Stack serial
					bw.ID(brTag.ID())   // Super Class Object ID
					bw.ID(brTag.ID())   // Class Loader ObjectID
					bw.ID(brTag.ID())   // Signers Object ID
					bw.ID(brTag.ID())   // Protection Domain Object ID
					bw.ID(brTag.ID())   // reserved1
					bw.ID(brTag.ID())   // reserved2
					bw.U32(brTag.U32()) // Instance size

					// Skip constant pool
					cpSize := bw.U16AsInt(brTag.U16())
					for i := 0; i < cpSize; i++ {
						bw.U16(brTag.U16()) // Constant pool index
						cpType := bw.U8AsInt(brTag.U8())
						bw.Bytes(brTag.Bytes(cpType))
					}

					// Skip static fields
					sSize := bw.U16AsInt(brTag.U16())
					for i := 0; i < sSize; i++ {
						bw.ID(brTag.ID())
						fieldType := bw.U8AsInt(brTag.U8())
						bw.Bytes(brTag.Bytes(brTag.Size(fieldType)))
					}

					// Skip instance fields
					iCount := bw.U16AsInt(brTag.U16())
					for i := 0; i < iCount; i++ {
						bw.ID(brTag.ID())
						bw.U8(brTag.U8())
					}
				case SubTagInstance: // 33 / 0x21
					bw.ID(brTag.ID())   // Object ID
					bw.U32(brTag.U32()) // Stack serial
					bw.ID(brTag.ID())   // Class Object ID
					bytesInInstance := bw.U32AsInt(brTag.U32())
					bw.Bytes(brTag.Bytes(bytesInInstance)) // Values
				case SubTagInstanceArray: // 34 / 0x22
					bw.ID(brTag.ID())   // Array ID
					bw.U32(brTag.U32()) // Stack serial
					elementCount := bw.U32AsInt(brTag.U32())
					bw.ID(brTag.ID()) // Element class ID
					for i := 0; i < elementCount; i++ {
						bw.ID(brTag.ID())
					}
				case SubTagValueArray:
					bw.ID(brTag.ID())
					bw.U32(brTag.U32())
					elementCount := bw.U32AsInt(brTag.U32())
					elementType := bw.U8AsInt(brTag.U8())
					arrayByteSize := brTag.Size(elementType) * elementCount
					if _, ok := opts.keep[elementType]; ok {
						bw.Bytes(brTag.Bytes(arrayByteSize))
					} else {
						var replaceChar byte
						if elementType != binaryreader.OHT_BOOL && elementType != binaryreader.OHT_FLOAT && elementType != binaryreader.OHT_DOUBLE {
							replaceChar = 'x'
						}
						brTag.Bytes(arrayByteSize)
						bw.Bytes(bytes.Repeat([]byte{replaceChar}, arrayByteSize))
					}
				default:
					panic(fmt.Sprintf("unrecognized object type: %d / 0x%x", objType, objType))
				}
			}
		case TagHeapDumpEnd: // This can be copied
			bw.Bytes(brTag.Bytes(brTag.Length()))
		default:
			panic(fmt.Sprintf("unrecognized tag: %d", tagType))
		}
	}

	io.Unmap(mmIn)
	_ = fIn.Close()
	io.Unmap(mmOut)
	_ = fOut.Close()
}

const (
	TagString          = 1
	TagLoadClass       = 2
	TagStackFrame      = 4
	TagStackTrace      = 5
	TagHeapDump        = 0xc
	TagHeapDumpSegment = 0x1c
	TagHeapDumpEnd     = 0x2c
)

const (
	SubTagRootJNIGlobal    = 0x01
	SubTagRootJNILocal     = 0x02
	SubTagRootJavaFrame    = 0x03
	SubTagRootStickyClass  = 0x05
	SubTagRootMonitorUsed  = 0x07
	SubTagRootThreadObject = 0x08
	SubTagClassDump        = 0x20
	SubTagInstance         = 0x21
	SubTagInstanceArray    = 0x22
	SubTagValueArray       = 0x23
)
