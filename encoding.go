package gomcstat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func ReadVarInt(r io.ByteReader) (int, error) {
	var num uint32
	var shift uint
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		num |= uint32(b&0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
		if shift >= 32 {
			return 0, errors.New("VarInt is too big")
		}
	}
	return int(int32(num)), nil
}

func WriteVarInt(w io.ByteWriter, value int) error {
	num := uint32(value)
	for {
		b := byte(num & 0x7F)
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		if err := w.WriteByte(b); err != nil {
			return err
		}
		if num == 0 {
			break
		}
	}
	return nil
}

func ReadString(r io.Reader) (string, error) {
	br, ok := r.(io.ByteReader)
	if !ok {
		return "", errors.New("reader must implement io.ByteReader")
	}
	length, err := ReadVarInt(br)
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func WriteString(w io.Writer, value string) error {
	bw, ok := w.(io.ByteWriter)
	if !ok {
		return errors.New("writer must implement io.ByteWriter")
	}
	b := []byte(value)
	if err := WriteVarInt(bw, len(b)); err != nil {
		return err
	}
	_, err := w.Write(b)
	return err
}

type PacketBuffer struct {
	bytes.Buffer
}

func (b *PacketBuffer) WriteVarInt(value int) error {
	return WriteVarInt(b, value)
}

func (b *PacketBuffer) WriteString(value string) error {
	return WriteString(b, value)
}

func (b *PacketBuffer) WriteUShort(value uint16) error {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, value)
	_, err := b.Write(buf)
	return err
}
