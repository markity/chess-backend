package protocol

import (
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
)

type Protocol struct{}

func (p *Protocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if buffer.Length() < 4 {
		return nil, nil
	}

	packetLength := buffer.PeekUint32()
	if packetLength == 0 {
		return nil, []byte{}
	}

	if buffer.Length() < 4+int(packetLength) {
		return nil, []byte{}
	}

	buffer.Retrieve(4)
	packetBytes := make([]byte, packetLength)
	copy(packetBytes, buffer.Bytes())
	buffer.Retrieve(int(packetLength))

	return nil, packetBytes
}

func (p *Protocol) Packet(c *gev.Connection, data interface{}) []byte {
	return data.([]byte)
}
