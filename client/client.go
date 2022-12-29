package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	gRPC "github.com/Draosakel/Mock2021/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "default", "Senders name")
var serverPort0 = flag.String("server0", "5000", "Tcp server")
var serverPort1 = flag.String("server1", "5001", "Tcp server")
var serverPort2 = flag.String("server2", "5001", "Tcp server")

var servers = []gRPC.TemplateClient{}
var serverConns []*grpc.ClientConn

var responseNumber = 0
var currentMessage string

func main() {
	//parse flag/arguments
	flag.Parse()
	f := setLog()
	defer f.Close()

	fmt.Println("--- CLIENT APP ---")

	//log to file instead of console
	//f := setLog()
	//defer f.Close()

	//connect to server and close the connection when program closes
	for i := 0; i < 3; i++ {
		serverName := "server" + strconv.Itoa(i) // name of the server
		serverPort := "500" + strconv.Itoa(i)    // port of the server port
		fmt.Println("Servername: " + serverName + " - Serverport: " + serverPort)
		ConnectToServer(serverName, serverPort)
	}

	//start the increment calls
	ParseInput()
}

// connect to server
func ConnectToServer(serverName string, serverPort string) {
	serverFlag := serverPort
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	fmt.Printf("client %s: Attempts to dial on port %s\n", *clientsName, serverFlag)
	conn, err := grpc.Dial(fmt.Sprintf(":%s", serverFlag), opts...)
	if err != nil {
		fmt.Printf("Fail to Dial : %v", err)
		return
	}

	servers = append(servers, gRPC.NewTemplateClient(conn))
	serverConns = append(serverConns, conn)
}

func ParseInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type 'Increment' to get value and increment by 1")
	fmt.Println("--------------------")
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input
		if input == "Increment" {
			incrementVal()
		}
	}
}

func incrementVal() {
	values := []int64{}
	//create amount type
	amount := &gRPC.Amount{
		ClientName: *clientsName,
		Value:      1, //cast from int to int32
	}
	log.Printf("Client called increment")
	for i := 0; i < 3; i++ {
		//Make gRPC call to server with amount, and recieve acknowlegdement back.
		ack, err := servers[i].Increment(context.Background(), amount)
		if err != nil || ack == nil {
			log.Printf("Port: 500%d server has crashed\n", i)
			fmt.Printf("Port: 500%d server has crashed\n", i)
			continue
		}
		values = append(values, ack.NewValue)
	}
	log.Printf("Returned increment value is %d\n", values[0])
	fmt.Printf("Returned increment value is %d\n", values[0])
}

// sets the logger to use a log.txt file instead of the console
func setLog() *os.File {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
