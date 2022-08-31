package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const ntpEpochOffset = 2208988800

type ntpPacket struct {
	Settings           byte      // leap indicator - 2bits , version - 3bits, mode - 3 bits
	Stratum            byte      // stratum level of the local clock
	Poll               byte      // exponent of the maximum interval between successive messages
	Precision          byte      // exponent of the precision of the local clock.
	RootDelay          uint32    // total round trip delay time
	RootDispersion     uint32    // max error aloud from primary clock source
	ReferenceID        uint32    // reference clock id
	ReferenceTimestamp [2]uint32 // reference timestamp - seconds and fraction of a second
	OriginTimestamp    [2]uint32 // originate timestamp - seconds and fraction of a second
	ReceiveTimestamp   [2]uint32 // received timestamp - seconds and fraction of a second
	TransmitTimestamp  [2]uint32 // transmit timestamp - seconds and fraction of a second
}

func main() {
	conn, err := net.Dial("udp", "88.147.254.232:123")
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		log.Fatal("failed to set deadline: ", err)
	}

	req := &ntpPacket{}
	req.Settings = 0x1B // no leap year warning, set ntp version 3, set client mode 3
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		log.Fatalf("failed to send request: %v", err)
	}

	rsp := &ntpPacket{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		log.Fatalf("failed to read server response: %v", err)
	}

	sec := int64(rsp.TransmitTimestamp[0]) - ntpEpochOffset
	//nsec := (int64(rsp.TransmitTimestamp[1]) * 1e9) >> 32
	nsec := (int64(rsp.TransmitTimestamp[1]) * 1e9) / 0xffffffff
	fmt.Printf("%v\n", time.Unix(sec, nsec))
}
