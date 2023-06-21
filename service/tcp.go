package service

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"log"
	"net"
	"strconv"
	"time"
)

var socket net.Conn
var dataPool []byte
var heartBeatTicker *time.Ticker = nil

func StartSocket(ip string) {
	var err interface{}
	socket, err = net.Dial("tcp", ip+":6000")
	if err != nil {
		println(err)
		time.Sleep(time.Second * 2)
		StartSocket(ip)
		return
	}
	log.Println("Connected TCP:", ip)
	startHeartBeat()
	go func() {
		for {
			if socket == nil {
				return
			}
			ReceiveSocketMessage()
		}
	}()
}

// CloseSocket 关闭socket
func CloseSocket() {
	socket.Close()
}

// SendSocketMessage 发送socket消息
func SendSocketMessage(data string, command int) {
	if socket == nil {
		return
	}
	stopHeartBeat()
	tmpToSend := make([]byte, 0)
	tmpToSend = append(tmpToSend, 0x4E)
	tmpToSend = append(tmpToSend, 0x66)
	dataBytes := []byte(data)
	dataLen := len(dataBytes)
	_dd := int(dataLen / 256)
	_dd1 := int(dataLen % 256)
	tmpToSend = append(tmpToSend, byte(_dd))
	tmpToSend = append(tmpToSend, byte(_dd1))
	_co1 := int(command / 256)
	_co2 := int(command % 256)
	tmpToSend = append(tmpToSend, byte(_co1))
	tmpToSend = append(tmpToSend, byte(_co2))
	tmpToSend = append(tmpToSend, dataBytes...)
	_Crc := crc32.ChecksumIEEE(tmpToSend[2:])
	dataCrc := bytes.NewBuffer([]byte{})
	binary.Write(dataCrc, binary.BigEndian, _Crc)
	_dataCrc := dataCrc.Bytes()
	tmpToSend = append(tmpToSend, _dataCrc...)
	socket.Write(tmpToSend)
	startHeartBeat()
	return
}

func ReceiveSocketMessage() {
	stopHeartBeat()
	data := make([]byte, 1024)
	info, _ := socket.Read(data)
	dataPool = append(dataPool, data[:info]...)
	handleDataPool()
	startHeartBeat()
}

func handleDataPool() {
	if len(dataPool) < 10 {
		return
	}
	index := -1
	for i := 0; i < len(dataPool); i++ {
		if dataPool[i] == 0x4E {
			index = i
			break
		}
	}
	if index == -1 {
		dataPool = []byte{}
		return
	}
	if dataPool[index+1] != 0x66 {
		dataPool = dataPool[index+1:]
		return
	}
	dataLength := int(dataPool[index+2])*256 + int(dataPool[index+3])
	println(dataLength)
	if len(dataPool) < 10+dataLength {
		return
	}
	commandInt := int(dataPool[index+4])*256 + int(dataPool[index+5])
	dataBytes := dataPool[index+6 : index+6+dataLength]
	message := Message{commandInt, string(dataBytes)}
	dataPool = dataPool[index+10+dataLength:]
	TCP2WS(message)
	handleDataPool()
}

func startHeartBeat() {
	if heartBeatTicker != nil {
		return
	}
	go sendHeartBeat()
}

func sendHeartBeat() {
	heartBeatTicker = time.NewTicker(2 * time.Second)
	for range heartBeatTicker.C {
		timestamp := time.Now().Unix()
		SendSocketMessage("{\"time\":"+strconv.FormatInt(timestamp, 10)+"}", 0x7266)
		if socket == nil {
			return
		}
	}
}

func stopHeartBeat() {
	if heartBeatTicker != nil {
		heartBeatTicker.Stop()
		heartBeatTicker = nil
	}
}
