package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	prototest "test.com/gotest/user_proto"
)

const (
	Address = "127.0.0.1:5656"
)

type userService struct {
}

func (u *userService) GetUser(ctx context.Context, req *prototest.UserRequest) (*prototest.UserResponse, error) {
	log.Println("receive client conn", req)
	user := map[string]string{
		"Name": "lily",
		"Age":  "26",
	}
	userBts, _ := json.Marshal(&user)

	resp := prototest.UserResponse{
		Status: "200",
		Msg:    "get user:" + req.Name,
		Data:   string(userBts),
	}
	return &resp, nil
}



func main() {
	userServer := userService{}

	certificate,err := tls.LoadX509KeyPair("server.crt", "server.key")

	if err!=nil{
		log.Fatalln("server load x509 error", err)
	}
	certpool := x509.NewCertPool()
	ca,err := ioutil.ReadFile("ca.crt")
	if err!=nil{
		log.Fatalln("ca file open error", err)
	}
	if ok:=certpool.AppendCertsFromPEM(ca);!ok{
		log.Fatalln("failed to append certs")
	}
	cred := credentials.NewTLS(&tls.Config{
		Certificates:                []tls.Certificate{certificate},
		ClientAuth:                  tls.RequireAndVerifyClientCert,
		ClientCAs:                   certpool,
	})

	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatal("server listen err", err)
	}
	//cred, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	//if err!=nil{
	//	log.Fatalln("server crt credential err", err)
	//}

	s := grpc.NewServer(grpc.Creds(cred))

	prototest.RegisterUserServer(s, &userServer)

	log.Println("start listening", Address)
	s.Serve(listen)
}
