package main

import (
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
	"test.com/gotest/user_proto/github.com/wisdomatom/prototest"
)

func main() {
	net.Listen("", "")
	userReq := prototest.UserRequest{
		Name: "jack",
		Age:  1660,
	}
	
	out, err := proto.Marshal(&userReq)
	if err != nil {
		log.Fatalln("failed to encode userRequest", err)
	}
	if err := ioutil.WriteFile("userReq", out, 0644); err != nil {
		log.Fatalln("failed to write userReq to file", err)
	}

	userReqRead := &prototest.UserRequest{}
	reqBin, err := ioutil.ReadFile("userReq")
	if err != nil {
		log.Fatalln("err read file", err)
	}
	err = proto.Unmarshal(reqBin, userReqRead)
	if err != nil {
		log.Fatalln("err unmashal", err)
	}
	log.Println("success unmarshal", userReqRead)
}
