package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/google/gopacket"
)

var GeeZipLayerType = gopacket.RegisterLayerType(
	12345,
	gopacket.LayerTypeMetadata{Name: "GeeZipLayer", Decoder: gopacket.DecodeFunc(decodeGeeZipLayer)},
)

// Implement my layer
type GeeZipLayer struct {
	TriggerFlag []byte
	CBAddr      net.IP
	CBPort      uint16
	CBString    string
	payload     []byte
}

func (m GeeZipLayer) LayerType() gopacket.LayerType {
	return GeeZipLayerType
}

func (m GeeZipLayer) LayerContents() []byte {
	return []byte(m.CBString)
}

func (m GeeZipLayer) LayerPayload() []byte {
	return m.payload
}

// Now implement a decoder... this one strips off the first 4 bytes of the
// packet.
func decodeGeeZipLayer(data []byte, p gopacket.PacketBuilder) error {
	var layer GeeZipLayer
	layer.TriggerFlag = data[0:8]
	layer.CBAddr = net.IPv4(data[8], data[9], data[10], data[11])
	layer.CBPort = binary.BigEndian.Uint16(data[12:14])
	layer.CBString = layer.CBAddr.String() + ":" + fmt.Sprintf("%d", layer.CBPort)

	p.AddLayer(layer)
	return p.NextDecoder(gopacket.LayerTypePayload)
}
