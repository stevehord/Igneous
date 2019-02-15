package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	tftp "igneous.io/tftp/cmd"
)

var (
	host     string
	storeMap map[string][]byte
	// store          []byte
	bufferSize     int
	bufferDataSize int
	res            tftp.PacketData
)

func main() {
	// TODO implement the in-memory tftp server
	host = "127.0.0.1:69"
	bufferSize = 516 // Default tftp spec buffer size??? Could be trouble
	bufferDataSize = 512
	storeMap = make(map[string][]byte)
	fmt.Println("Hello TFTP")

	conn, err := net.ListenPacket("udp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close() // Some Go magic but I am not sure it ever executes
	fmt.Printf("Started on %v\n", conn.LocalAddr())

	for {
		buf := make([]byte, bufferSize) //TODO
		bufLength, sourceAddr, err := conn.ReadFrom(buf)

		// fmt.Printf("buffLength: %d\n", bufLength)
		// fmt.Printf("addr: %s\n", sourceAddr)
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
		fmt.Printf("Parse Error: %s", err)
	}

	switch v := p.(type) {
	// Process the header
	case *tftp.PacketRequest:
		pkRequest := p.(*tftp.PacketRequest)
		// fmt.Printf("Packet Operation %v\n", pkRequest.Op)
		// fmt.Printf("Packet Mode %v\n", pkRequest.Mode)
		// fmt.Printf("Packet filename %v\n", pkRequest.Filename)

		// Process a GET
		if pkRequest.Op == 1 {
			processGet(conn, pkRequest, srcAddr)
		}

		// Process a PUT ack
		if pkRequest.Op == 2 {
			res := tftp.PacketAck{BlockNum: 0}
			conn.WriteTo(res.Serialize(), srcAddr)
		}
	// Process the main Data
	case *tftp.PacketData:
		source := srcAddr.String()
		pkData := p.(*tftp.PacketData)
		// fmt.Printf("Packet Data Length %v\n", bufLength)
		storeMap[source] = append(storeMap[source], pkData.Data...)
		// fmt.Println(string(pkData.Data))
		res := tftp.PacketAck{BlockNum: pkData.BlockNum}
		conn.WriteTo(res.Serialize(), srcAddr)

		// Requirement to print the entire file at the end of the upload.
		// Any bufer less than 512 should indicate the upload is complete.
		if bufLength < bufferDataSize {
			printBuffer(storeMap[source])
			// storeMap[source] = make([]byte, 0) //Clear the buffer? TODO: Look for memory leak here.
			delete(storeMap, source)
		}
	// Process Packet errors  TODO: Needs testcase
	case *tftp.PacketError:
		pkData := p.(*tftp.PacketError)
		fmt.Printf("Transfer Error: %v\n", string(pkData.Msg))
		res := tftp.PacketAck{BlockNum: 1}
		conn.WriteTo(res.Serialize(), srcAddr)

	case *tftp.PacketAck:
		pkData := p.(*tftp.PacketAck)
		fmt.Printf("PacketAck: %v\n", string(pkData.BlockNum))

	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

func processGet(conn net.PacketConn, pkRequest *tftp.PacketRequest, srcAddr net.Addr) {
	content := readFileContent(pkRequest.Filename)
	fmt.Println(content[:10])

	// Create a packet of 512 or the actual size of remainning buffer
	bufLengthRemain := len(content)
	var blockNum uint16
	var x int

	for {
		blockNum++
		// fmt.Printf("bufLengthRemain: %v  x:%v block:%v\n", bufLengthRemain, x, blockNum)
		if bufLengthRemain > bufferDataSize {
			res = tftp.PacketData{blockNum, content[x : x+bufferDataSize]}
			bufLengthRemain = bufLengthRemain - bufferDataSize
			x = x + bufferDataSize
			conn.WriteTo(res.Serialize(), srcAddr)
		} else {
			res = tftp.PacketData{blockNum, content[x : x+bufLengthRemain]}
			conn.WriteTo(res.Serialize(), srcAddr)
			//send empty packet to close connection
			res = tftp.PacketData{blockNum + 1, make([]byte, 0)}
			conn.WriteTo(res.Serialize(), srcAddr)
			break
		}
	}
}

func printBuffer(buf []byte) {
	fmt.Println("----------------------")
	fmt.Println(string(buf))
	fmt.Println("----------------------")
}

func writeBufferToFileSystem(path string, buf []byte) {
	// open output file
	fo, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	// close fo on exit print exit if any
	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	n, err := fo.Write(buf)
	fmt.Printf("%v Bytes Written", n)
	if err != nil {
		fmt.Println(err)
	}
}

func readFileContent(filePath string) (fileContent []byte) {
	file, err := os.Open("testFiles/" + filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close() // golang magic

	b, err := ioutil.ReadAll(file)
	return b
}
