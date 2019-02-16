package main

import (
	"bytes"
	"fmt"
	tftp "igneous/tftp/cmd"
	"log"
	"net"
	"unicode"
)

var (
	host           string
	storeMap       map[string]string
	fileStore      map[string][]byte
	bufferSize     int
	bufferDataSize int
)

func main() {
	// TODO implement the in-memory tftp server

	fmt.Println("Hello TFTP")

	host = "127.0.0.1:69"
	bufferSize = 516     // Default tftp spec buffer size??? Could be trouble
	bufferDataSize = 512 // Typical size of the data packets
	fileStore = make(map[string][]byte)
	storeMap = make(map[string]string)

	conn, err := net.ListenPacket("udp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close() // Some Go magic but I am not sure it ever executes
	fmt.Printf("Started on %v\n", conn.LocalAddr())

	for {
		buf := make([]byte, bufferSize) //TODO
		bufLength, sourceAddr, err := conn.ReadFrom(buf)

		if err != nil {
			fmt.Printf("Read Error: %s\n", err)
			continue
		}

		go process(buf, conn, sourceAddr, bufLength)
	}

}

func process(buf []byte, conn net.PacketConn, srcAddr net.Addr, bufLength int) {

	p, err := tftp.ParsePacket(buf)
	if err != nil {
		fmt.Printf("Halt - Parse Error: %s", err)
	}

	switch v := p.(type) {
	// Process the header
	case *tftp.PacketRequest:
		pkRequest := p.(*tftp.PacketRequest)

		// Process a GET
		if pkRequest.Op == 1 {
			processGet(conn, pkRequest, srcAddr)
		}

		// Process a PUT ack
		if pkRequest.Op == 2 {
			//Map the source Addr + port to a specfic file name (poormans threading)
			storeMap[srcAddr.String()] = pkRequest.Filename

			// Reset the fileStore for a new file
			fileStore[pkRequest.Filename] = make([]byte, 0, 0)

			// return a Ack to start the data transfer
			res := tftp.PacketAck{BlockNum: 0}
			conn.WriteTo(res.Serialize(), srcAddr)
		}

	// Process the main PUT data packets
	case *tftp.PacketData:
		pkData := p.(*tftp.PacketData)

		// Unpack the file name using the upload address + port
		fileName := storeMap[srcAddr.String()]

		//Store the date with control chars removed - TODO: verify this is correct behavior
		fileStore[fileName] = append(fileStore[fileName], bytes.TrimFunc(pkData.Data, unicode.IsControl)...)

		// Return Ack with block number to start next block
		res := tftp.PacketAck{BlockNum: pkData.BlockNum}
		conn.WriteTo(res.Serialize(), srcAddr)

		// Requirement to print the entire file at the end of the upload.
		// Any bufer less than 512 should indicate the upload is complete.
		if bufLength < bufferDataSize {
			// fmt.Println(len(fileStore[fileName]))
			printBuffer(fileStore[fileName])
		}

	// Process Packet errors  TODO: Needs testcase
	case *tftp.PacketError:
		pkData := p.(*tftp.PacketError)
		fmt.Printf("Transfer Error: %v\n", string(pkData.Msg))
		res := tftp.PacketAck{BlockNum: 1}
		conn.WriteTo(res.Serialize(), srcAddr)

	// Nothing to do here but acknolege the response
	case *tftp.PacketAck:
		pkData := p.(*tftp.PacketAck)
		fmt.Printf("PacketAck: %v\n", string(pkData.BlockNum))

	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

/*
 * helper function to manage downloading of files larger than bufferDataSize
 *
 */
func processGet(conn net.PacketConn, pkRequest *tftp.PacketRequest, srcAddr net.Addr) {
	content := fileStore[pkRequest.Filename]
	fmt.Println(len(content))

	// Create a packet of 512 or the actual size of remainning buffer
	bufLengthRemain := len(content)
	var blockNum uint16
	var x int

	// Segment the file into 512k sections and send
	// last segement is empty to signal download complete
	for {
		blockNum++
		// fmt.Printf("bufLengthRemain: %v  x:%v block:%v\n", bufLengthRemain, x, blockNum)
		if bufLengthRemain > bufferDataSize {
			res := tftp.PacketData{blockNum, content[x : x+bufferDataSize]}
			bufLengthRemain = bufLengthRemain - bufferDataSize
			x = x + bufferDataSize
			conn.WriteTo(res.Serialize(), srcAddr)
		} else {
			res := tftp.PacketData{blockNum, content[x : x+bufLengthRemain]}
			conn.WriteTo(res.Serialize(), srcAddr)

			//send empty packet to end connection
			res = tftp.PacketData{blockNum + 1, make([]byte, 0)}
			conn.WriteTo(res.Serialize(), srcAddr)
			break
		}
	}
}

/*
 * Print the buffer to stdout with help to make readable
 *
 */
func printBuffer(buf []byte) {
	fmt.Println("------Start----------------")
	fmt.Println(string(buf))
	fmt.Println("------End----------------")
}
