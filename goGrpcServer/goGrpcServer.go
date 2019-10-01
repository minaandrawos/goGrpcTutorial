package main

import (
	"context"
	"encoding/csv"
	"flag"
	"goGrpcTutorial/grpcTest"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

func main() {
	var (
		output        string
		serverAddress string
	)
	flag.StringVar(&output, "o", "output.csv", "output directory")
	flag.StringVar(&serverAddress, "a", "127.0.0.1:9000", "destination server address")
	flag.Parse()
	sendChan := make(chan *grpcTest.TestMessage)
	l, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	testServer := &TestGrpcServer{sendChan}
	grpcTest.RegisterTestServiceServer(grpcServer, testServer)
	go grpcServer.Serve(l)
	for message := range sendChan {
		writeValuesTofile(message)
	}
}

//TestGrpcServer Test GRPC server struct
type TestGrpcServer struct {
	output chan<- *grpcTest.TestMessage
}

//SendMessage receives GRPC message
func (s *TestGrpcServer) SendMessage(ctx context.Context, in *grpcTest.TestMessage) (*grpcTest.Reply, error) {
	s.output <- in
	return &grpcTest.Reply{}, nil
}

func writeValuesTofile(datatowrite *grpcTest.TestMessage) {

	//Retreive client information from the protobuf message
	ClientName := datatowrite.GetClientName()
	ClientDescription := datatowrite.GetDescription()
	ClientID := strconv.Itoa(int(datatowrite.GetClientId()))

	// retrieve the message items list
	items := datatowrite.GetMessageitems()
	log.Println("Writing value to CSV file")
	//Open file for writes, if the file does not exist then create it
	file, err := os.OpenFile("CSVValues.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	//make sure the file gets closed once the function exists
	defer file.Close()
	//Go through the list of message items, insert them into a string array then write them to the CSV file.
	writer := csv.NewWriter(file)
	for _, item := range items {
		record := []string{ClientID, ClientName, ClientDescription, strconv.Itoa(int(item.GetId())), item.GetItemName(), strconv.Itoa(int(item.GetItemValue())), strconv.Itoa(int(item.GetItemType()))}
		writer.Write(record)
	}
	//flush data to the CSV file
	writer.Flush()
	log.Println("Finished Writing values to CSV file")
}
