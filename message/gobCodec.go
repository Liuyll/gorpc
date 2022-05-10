package message

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	dec  *gob.Decoder
	enc  *gob.Encoder
	buf  *bufio.Writer
	mu *sync.Mutex
}

func NewGobCodec(conn io.ReadWriteCloser) GobCodec {
	buf := bufio.NewWriter(conn)
	return GobCodec{
		conn,
		gob.NewDecoder(conn),
		gob.NewEncoder(buf),
		buf,
		new(sync.Mutex),
	}
}

func (c GobCodec) Write(h interface{}, body interface{}) {
	defer func() {
		c.mu.Unlock()
		if err := c.buf.Flush(); err != nil {
			fmt.Println("flush error:", err)
		}
	}()

	c.mu.Lock()
	if err := c.enc.Encode(h); err != nil {
		fmt.Println("encode header err:", err)
	}
	if body != nil {
		if err := c.enc.Encode(body); err != nil {
			fmt.Println("encode body err: \n", err, fmt.Sprintf("debug stack: %s", debug.Stack()))
		}
	}
}

func (c GobCodec) ReadHeader(h interface{}) error {
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
