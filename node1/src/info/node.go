package info

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

var Nodo1Host string = "localhost:8001"
var Nodo2Host string = "localhost:8002"
var Nodo3Host string = "localhost:8003"

type nodeService struct {
	IncomingInfoChan chan IncomingInfo
	OutgoingInfoChan chan OutgoingInfo
	Tickets          []string
}

type NodeService interface {
	AddToIncomingInfoChan(IncomingInfo)
	AddToOutgoingInfoChan(OutgoingInfo)
	DecodeFromNode3(string)
	HandleConn(net.Conn)
	PrepareOutgoingInfo(OutgoingInfo)
	SendToNode2(OutgoingInfo)
	GetIncomingInfoFromChan() IncomingInfo
	Start()
}

func NewNodeService() NodeService {
	nodeService := &nodeService{}
	nodeService.IncomingInfoChan = make(chan IncomingInfo, 10)
	nodeService.OutgoingInfoChan = make(chan OutgoingInfo, 10)
	go nodeService.Start()
	return nodeService
}

func (ns *nodeService) Start() {
	log.Println("Tcp server started on port 8001...")

	listener, err := net.Listen("tcp", Nodo1Host)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	rand.Seed(time.Now().Unix())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go ns.HandleConn(conn)
	}
}

func (ns *nodeService) HandleConn(conn net.Conn) {
	defer conn.Close()

	log.Printf("Serving %s\n", conn.RemoteAddr().String())

	netData, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}

	temp := strings.TrimSpace(string(netData))
	log.Println("Data received:", temp)

	if strings.HasPrefix(temp, "NODE3:") {
		ns.DecodeFromNode3(temp)
	}

	for i, v := range ns.Tickets {
		if v == temp {
			outgoingInfo := <-ns.OutgoingInfoChan
			ns.SendToNode2(outgoingInfo)
			ns.Tickets = FastRemove(ns.Tickets, i)
			break
		}
	}
}

func (ns *nodeService) PrepareOutgoingInfo(outgoingInfo OutgoingInfo) {
	conn, _ := net.Dial("tcp", Nodo1Host)
	defer conn.Close()
	fmt.Fprintln(conn, outgoingInfo.Ticket)

	ns.Tickets = append(ns.Tickets, outgoingInfo.Ticket)

	go ns.AddToOutgoingInfoChan(outgoingInfo)
}

func (ns *nodeService) AddToOutgoingInfoChan(outgoingInfo OutgoingInfo) {
	ns.OutgoingInfoChan <- outgoingInfo
}

func (ns *nodeService) AddToIncomingInfoChan(incomingInfo IncomingInfo) {
	ns.IncomingInfoChan <- incomingInfo
}

func (ns *nodeService) GetIncomingInfoFromChan() IncomingInfo {
	return <-ns.IncomingInfoChan
}

func (ns *nodeService) SendToNode2(info OutgoingInfo) {
	conn, _ := net.Dial("tcp", Nodo2Host)
	defer conn.Close()

	bytesMsg, _ := json.Marshal(info)
	fmt.Fprintln(conn, string(bytesMsg))
}

func (ns *nodeService) DecodeFromNode3(jsonPayload string) {
	jsonPayload = strings.TrimPrefix(jsonPayload, "NODE3:")

	var info IncomingInfo
	json.Unmarshal([]byte(jsonPayload), &info)

	go ns.AddToIncomingInfoChan(info)
}
