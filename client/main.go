package main

import (
	"crypto/tls"
	"crypto/x509"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	prototest "test.com/gotest/user_proto"
	"time"
)

const (
	Address = "127.0.0.1:5656"
)


func main()  {
	//conn,err := net.Dial("tcp", Address)
	//cred, err := credentials.NewClientTLSFromFile("server.crt", "server.grpc.io")
	certificate,err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err!=nil{
		log.Fatalln("cert credential create err", err)
	}
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.crt")
	if err !=nil{
		log.Fatalln("ca load failed", err)
	}
	if ok:= certpool.AppendCertsFromPEM(ca);!ok{
		log.Fatalln("failed to append ca certs")
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates:                []tls.Certificate{certificate},
		RootCAs:                     certpool,
		ServerName:                  "server.io",
	})
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))


	if err!=nil{
		log.Fatalln("client dial error", err)
	}
	defer conn.Close()
	req := prototest.UserRequest{
		Name: "iris",
	}
	ctx := context.Background()
	client := prototest.NewUserClient(conn)
	for i:=0;i<66;i++{
		res, err := client.GetUser(ctx, &req)
		if err!=nil{
			//log.Fatalln("client request error", err)
			log.Println("client request error", err)
			time.Sleep(time.Second * 2)
			continue
		}
		log.Println("response>>>>>>", res.Status, res.Msg, res.Data)

		time.Sleep(time.Second * 2)
	}


}

