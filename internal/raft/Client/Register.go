package raft

import (
	"bytes"
	"cache/config"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type ServiceRegistryInf interface {
	InitServiceRegister(*config.Config)
	DeregisterService(*config.Config)
	FetchJWTToken() (string, error)
}

type ServiceRegistry struct {
	client *Client
	Secret string
}

type NODETYPE int

const (
	NEWJOIN NODETYPE = iota
	LEAVE
)

type ClusterMember struct {
	NodeType  NODETYPE
	NodeID    string
	NodeAddr  string
	NodePort  string
	ClusterID string
	GrpcPort  string
}

func (sr *ServiceRegistry) jnitServiceRegister(config *config.Config) {
	conn, err := net.Dial("tcp", config.ServiceDiscoveryAddr)
	if err != nil {
		log.Fatal("connect service register fail", err)
	}
	msg := bytes.NewBuffer(nil)
	member := &ClusterMember{
		NodeType:  NEWJOIN,
		NodeID:    "2",
		NodeAddr:  config.Host,
		NodePort:  config.Port,
		ClusterID: config.ClusterUUID,
		GrpcPort:  sr.client.NodeDetails.GrpcPort,
	}
	encodedJson, _ := json.Marshal(member)
	err = binary.Write(msg, binary.BigEndian, uint8(0))
	if err != nil {
		log.Fatal("marshal service register fail", err)
	}
	msg.Write(encodedJson)
	write, err := conn.Write(msg.Bytes())
	//defer conn.Close()
	if err != nil {
		fmt.Println("err", err)
	}
	go readResponse(sr, conn)
	fmt.Println("written bytes :", write)
}

func (sr *ServiceRegistry) DeregisterService(config *config.Config) {
	log.Printf("Node  %s is deregistered from the cluster %s", config.Host+":"+config.Port, config.ServiceDiscoveryAddr)
	conn, err := net.Dial("tcp", config.ServiceDiscoveryAddr)
	if err != nil {
		log.Fatal("connect service register fail", err)
	}
	msg := bytes.NewBuffer(nil)
	member := &ClusterMember{
		NodeType:  LEAVE,
		NodeID:    "2",
		NodeAddr:  config.Host,
		NodePort:  config.Port,
		ClusterID: config.ClusterUUID,
		GrpcPort:  sr.client.NodeDetails.GrpcPort,
	}
	encodedJson, _ := json.Marshal(member)
	err = binary.Write(msg, binary.BigEndian, uint8(1))
	if err != nil {
		log.Fatal("marshal service register fail", err)
	}
	msg.Write(encodedJson)
	write, err := conn.Write(msg.Bytes())
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("written bytes :", write)
}

/*
reads the reposne from the
*/
func readResponse(service *ServiceRegistry, conn net.Conn) {
	var buffer bytes.Buffer
	var resp *NodeRegisterationResponse = &NodeRegisterationResponse{}
	//closing the channel when to have graceful shutdown
	defer close(service.client.MemberChannel)
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			// Handle error or potential connection close
		}
		buffer.Write(data[:n])
		if bytes.IndexByte(buffer.Bytes(), '\n') > -1 { // Check for newline delimiter (replace if needed)
			break
		}
		err = json.Unmarshal(data[:n], resp)
		if err != nil {
			log.Println("unmarshal fail", err)
			return
		}
		service.Secret = resp.Secret
		for _, peer := range resp.Members {
			fmt.Println("grpc port ", peer.GrpcPort)
			fmt.Printf("Secret %s & Total Cluster Members %v\n", service.Secret, len(resp.Members))
			service.client.MemberChannel <- resp.Members
		}

	}
}

func (sr *ServiceRegistry) RegisterRaftClient(config *config.Config) {
	sr.jnitServiceRegister(config)
}
