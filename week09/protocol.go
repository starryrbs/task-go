package week09

import (
	"encoding/binary"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_bodyOffset   = _seqOffset + _seqSize
)

type Proto struct {
	packLen   uint32
	headerLen uint32
	Ver       uint16
	Op        uint32
	Seq       uint32
	Body      []byte
}

func decode(buf []byte) *Proto {
	packLen := binary.BigEndian.Uint32(buf[_packOffset:])
	headerLen := binary.BigEndian.Uint32(buf[_headerOffset:])
	ver := binary.BigEndian.Uint16(buf[_verOffset:])
	op := binary.BigEndian.Uint32(buf[_opOffset:])
	seq := binary.BigEndian.Uint32(buf[_seqOffset:])

	return &Proto{
		packLen:   packLen,
		headerLen: headerLen,
		Ver:       ver,
		Op:        op,
		Seq:       seq,
		Body:      buf[_bodyOffset:],
	}
}

func encode(data string) []byte {
	packLen := len(data) + _rawHeaderSize
	version := 1
	op := 1
	seq := 1
	buf := make([]byte, packLen)

	binary.BigEndian.PutUint32(buf[_packOffset:], uint32(packLen))
	binary.BigEndian.PutUint16(buf[_headerOffset:], uint16(_rawHeaderSize))
	binary.BigEndian.PutUint16(buf[_verOffset:], uint16(version))
	binary.BigEndian.PutUint32(buf[_opOffset:], uint32(op))
	binary.BigEndian.PutUint32(buf[_seqOffset:], uint32(seq))

	byteBody := []byte(data)
	copy(buf[_bodyOffset:], byteBody)
	return buf
}
