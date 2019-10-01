package main

import (
	"context"
	"encoding/csv"
	"flag"
	"goGrpcTutorial/grpcTest"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

type Headers []string

func (h Headers) getHeaderIndex(headername string) int {
	if len(headername) >= 2 {
		for index, s := range h {
			if s == headername {
				return index
			}
		}
	}
	return -1
}

const ClientName = "GoClient"
const ClientID = 2
const ClientDescription = "This is a Go Protobuf client!!"

func main() {
	var (
		srvAddress string
		inputFile  string
	)
	flag.StringVar(&srvAddress, "sa", "127.0.0.1:9000", "server destination address")
	flag.StringVar(&inputFile, "i", "csvv.csv", "input csv file to process")
	flag.Parse()

	conn, err := grpc.Dial(srvAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := grpcTest.NewTestServiceClient(conn)
	if msg, err := retrieveDataFromFile(inputFile); err == nil {
		client.SendMessage(context.Background(), msg)
	}

}

func retrieveDataFromFile(fname string) (*grpcTest.TestMessage, error) {
	file, err := os.Open(fname)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	csvreader := csv.NewReader(file)
	var hdrs Headers
	hdrs, err = csvreader.Read()
	if err != nil {
		return nil, err
	}

	ITEMIDINDEX := hdrs.getHeaderIndex("itemid")
	ITEMNAMEINDEX := hdrs.getHeaderIndex("itemname")
	ITEMVALUEINDEX := hdrs.getHeaderIndex("itemvalue")
	ITEMTYPEINDEX := hdrs.getHeaderIndex("itemType")

	message := new(grpcTest.TestMessage)
	message.ClientName = ClientName
	message.ClientId = ClientID
	message.Description = ClientDescription

	var loopErr error
	//loop through the records
	for {
		record, err := csvreader.Read()
		if err != nil {
			break
		}
		//Populate items
		testMessageItem := new(grpcTest.TestMessage_MsgItem)
		itemid, err := strconv.Atoi(record[ITEMIDINDEX])
		if err != nil {
			loopErr = err
			break
		}
		testMessageItem.Id = int32(itemid)
		testMessageItem.ItemName = record[ITEMNAMEINDEX]
		itemvalue, err := strconv.Atoi(record[ITEMVALUEINDEX])
		if err != nil {
			loopErr = err
			break
		}
		testMessageItem.ItemValue = int32(itemvalue)
		itemtype, err := strconv.Atoi(record[ITEMTYPEINDEX])
		if err != nil {
			loopErr = err
			break
		}
		iType := grpcTest.TestMessage_ItemType(itemtype)
		testMessageItem.ItemType = iType

		message.Messageitems = append(message.Messageitems, testMessageItem)

	}

	//fmt.Println(ProtoMessage.Messageitems)
	return message, loopErr
}
