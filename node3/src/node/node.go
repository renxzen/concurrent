package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"nodo3/src/info"
	"strings"
	"time"
)

var Nodo1Host string = "localhost:8001"
var Nodo2Host string = "localhost:8002"
var Nodo3Host string = "localhost:8003"

type nodeService struct{}

type NodeService interface {
	Start()
	HandleConn(net.Conn)
	SendToNode1(info.IncomingInfo)
	DecodeStuff(string) info.IncomingInfo
	DecodeVaccine(int) string
	DecodeGender(int) string
}

func NewNodeService() NodeService {
	nodeService := &nodeService{}
	nodeService.Start()
	return nodeService
}

func (ns *nodeService) Start() {
	log.Println("Tcp server started on port 8003...")

	listener, err := net.Listen("tcp", Nodo3Host)
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
	incomingInfo := ns.DecodeStuff(temp)
	ns.SendToNode1(incomingInfo)

}

func (ns *nodeService) SendToNode1(incomingInfo info.IncomingInfo) {
	conn, _ := net.Dial("tcp", Nodo1Host)
	defer conn.Close()

	outgoingInfo := info.OutgoingInfo{
		Age:           incomingInfo.Age,
		Gender:        ns.DecodeGender(incomingInfo.Gender),
		FirstVaccine:  ns.DecodeVaccine(incomingInfo.FirstVaccine),
		SecondVaccine: ns.DecodeVaccine(incomingInfo.SecondVaccine),
		Prediction:    incomingInfo.Prediction,
		Ticket:        incomingInfo.Ticket,
	}

	bytesMsg, _ := json.Marshal(outgoingInfo)
	fmt.Fprintln(conn, "NODE3:"+string(bytesMsg))
}

func (ns *nodeService) DecodeStuff(jsonPayload string) info.IncomingInfo {
	jsonPayload = strings.TrimPrefix(jsonPayload, "COMMON:")
	var info info.IncomingInfo
	json.Unmarshal([]byte(jsonPayload), &info)
	return info
}

func (ns *nodeService) DecodeVaccine(vaccine int) string {
	vaccines := []string{"Ninguna", "Sinopharm", "Pfizer", "AstraZeneca"}

	return vaccines[vaccine]
}

func (ns *nodeService) DecodeGender(gender int) string {
	genders := []string{"Female", "Male"}

	return genders[gender]
}
