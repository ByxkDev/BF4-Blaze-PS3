package blaze

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	TDF_INT32  = 0x00
	TDF_STRING = 0x01
	TDF_BLOB   = 0x02
)

type TDF struct {
	Tag string
	Type byte
	Value interface{}
}

func WriteTDF(buf *bytes.Buffer, tag string, value string,) {
	for i:=0;i<4;i++ {

		if i < len(tag) {
			buf.WriteByte(tag[i])
		}else{
			buf.WriteByte(0)
		}
	}

	buf.WriteByte(TDF_STRING) //type string
	binary.Write(buf, binary.BigEndian, uint32(len(value)+1),) //length
	buf.WriteString(value)
	buf.WriteByte(0)
}

func ReadTDF(data []byte) []TDF {
	var result []TDF
	offset := 0

	for offset+5 <= len(data) {
		tag := string(data[offset:offset+4])
		offset += 4
		t := data[offset]
		offset++

		switch t {
		case TDF_STRING:
			if offset+4 > len(data) {
				return result
			}

			length := binary.BigEndian.Uint32(data[offset:offset+4],)

			offset += 4
			if offset+int(length) > len(data) {
				return result
			}

			value := string(data[offset:offset+int(length)-1],)
			offset += int(length)
			result = append(result, TDF{Tag:tag, Type:t, Value:value,},)

		default:
			fmt.Printf("Unknown TDF type %02X\n", t,)
			return result
		}
	}

	return result
}