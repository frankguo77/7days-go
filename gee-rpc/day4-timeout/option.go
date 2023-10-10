package geerpc

import (
	"geerpc/codec"
	"time"
)

const MagicNumber = 0x3bef5c
const ADDR = "127.0.0.1:20000"

type Option struct {
	MagicNumber    int        // MagicNumber marks this's a geerpc request
	CodecType      codec.Type // client may choose different Codec to encode body
	ConnectTimeout time.Duration
	HandleTimeout  time.Duration
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber, 
	CodecType: codec.GobType,
	ConnectTimeout: time.Second * 10,
}