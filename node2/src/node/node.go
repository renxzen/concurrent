package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"nodo2/src/tree/model"
	"nodo2/src/tree/trees"
	"strconv"
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
	SendToNode3(model.OutgoingInfo)
	DecodeStuff(string) model.IncomingInfo
}

func NewNodeService() NodeService {
	nodeService := &nodeService{}
	nodeService.Start()
	return nodeService
}

func (ns *nodeService) Start() {
	log.Println("Tcp server started on port 8002...")

	listener, err := net.Listen("tcp", Nodo2Host)
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

	info := model.Info{
		Age:           incomingInfo.Age,
		Gender:        incomingInfo.Gender,
		FirstVaccine:  incomingInfo.FirstVaccine,
		SecondVaccine: incomingInfo.SecondVaccine,
	}

	prediction := trees.GetSinglePrediction(info)
	predictionInt, _ := strconv.Atoi(prediction)

	outgoingInfo := model.OutgoingInfo{
		Age:           info.Age,
		Gender:        info.Gender,
		FirstVaccine:  info.FirstVaccine,
		SecondVaccine: info.SecondVaccine,
		Prediction:    predictionInt,
		Ticket:        incomingInfo.Ticket,
	}

	ns.SendToNode3(outgoingInfo)
}

func (ns *nodeService) SendToNode3(outgoingInfo model.OutgoingInfo) {
	conn, _ := net.Dial("tcp", Nodo3Host)
	defer conn.Close()

	bytesMsg, _ := json.Marshal(outgoingInfo)
	fmt.Fprintln(conn, string(bytesMsg))
}

func (ns *nodeService) DecodeStuff(jsonPayload string) model.IncomingInfo {
	jsonPayload = strings.TrimPrefix(jsonPayload, "COMMON:")
	var incomingInfo model.IncomingInfo
	json.Unmarshal([]byte(jsonPayload), &incomingInfo)
	return incomingInfo
}
