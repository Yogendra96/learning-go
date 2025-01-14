package minidb

import "encoding/binary"

const entryHeaderSize = 10

const (
	PUT uint16 = iota
	DEL
)

// Entry 写入文件的记录
type Entry struct {
	Key       []byte
	Value     []byte
	KeySize   uint32
	ValueSize uint32
	Mark      uint16
}

func NewEntry(key, value []byte, mark uint16) *Entry {
	return &Entry{
		Key:       key,
		Value:     value,
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Mark:      mark,
	}
}

func (e *Entry) GetSize() int64 {
	return int64(entryHeaderSize + e.KeySize + e.ValueSize)
}

// Encode 编码 Entry，返回字节数组 在写入数据的时候，需要Encode
func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, e.GetSize())
	// 先处理各个部分size的值，0:10 存储size信息
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	binary.BigEndian.PutUint16(buf[8:10], e.Mark)
	// 再处理key value的值
	copy(buf[entryHeaderSize:entryHeaderSize+e.KeySize], e.Key)
	copy(buf[entryHeaderSize+e.KeySize:], e.Value)
	return buf, nil
}

// 解码 buf 字节数组，返回 Entry,在读取数据的时候，需要Decode
func Decode(buf []byte) (*Entry, error) {
	ks := binary.BigEndian.Uint32(buf[0:4])
	vs := binary.BigEndian.Uint32(buf[4:8])
	mark := binary.BigEndian.Uint16(buf[8:10])
	return &Entry{KeySize: ks, ValueSize: vs, Mark: mark}, nil
}
