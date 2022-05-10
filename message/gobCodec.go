package message

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	dec  *gob.Decoder
	enc  *gob.Encoder
	buf  *bufio.Writer
}

func NewGobCodec(conn io.ReadWriteCloser) GobCodec {
	buf := bufio.NewWriter(conn)
	return GobCodec{
		conn,
		gob.NewDecoder(conn),
		gob.NewEncoder(buf),
		buf,
	}
}

func (c GobCodec) Write(h *RPCHeader, body interface{}) {
	defer func() {
		if err := c.buf.Flush(); err != nil {
			fmt.Println("flush error:", err)
		}
	}()

	if err := c.enc.Encode(h); err != nil {
		fmt.Println("encode header err $err")
	}
	if err := c.enc.Encode(body); err != nil {
		fmt.Println("encode body err $err")
	}
}

func (c GobCodec) ReadHeader(h *RPCHeader) error {
	if err := c.dec.Decode(h); err != nil {
		return err
	}

	return nil
}

func (c GobCodec) ReadBody(body interface{}) error {
	if err := c.dec.Decode(body); err != nil {
		return err
	}
	return nil
}

func (c GobCodec) Close() {
	c.conn.Close()
}
