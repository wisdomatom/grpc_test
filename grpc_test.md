#### gRPC with TSL certificate

- 证书认证

```shell script
# 生成服务端私钥 长度2048
openssl genrsa -out server.key 2048
# 生成服务端证书 有效期10年
openssl req -new -x509 -days 3650 \
-subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io" \
-key server.key -out server.crt

# 生成客户端私钥
openssl genrsa -out client.key 2048

# 生成客户端证书 有效期10年
openssl req -new -x509 -days 3650 \
-subj "/C=GB/L=China/O=grpc-client/CN=client.grpc.io" \
-key client.key -out client.crt
```

```go
package main
// server with crt
import (
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "log"
)
func main(){
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err!=nil{
        log.Fatal(err)
    }
    server := grpc.NewServer(grpc.Creds(creds))
    
    fmt.Println("server", server)
    // ...
}
```

```go
package main
import (

"google.golang.org/grpc"
"google.golang.org/grpc/credentials"
"log"
)
// client with crt

func main()  {
    creds, err := credentials.NewClientTLSFromFile("server.crt", "server.grpc.io")
    if err!=nil{
        log.Fatal(err)
    }
    conn, err := grpc.Dial("localhost:5656", grpc.WithTransportCredentials(creds))
    if err != nil{log.Fatal(err)}
    defer conn.Close()
}
```

- 根证书验证

```shell script
$ openssl genrsa -out ca.key 2048
$ openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=gobook/CN=github.com" \
    -key ca.key -out ca.crt

# 重新对服务器端证书进行签名
$ openssl req -new \
    -subj "/C=GB/L=China/O=server/CN=server.io" \
    -key server.key \
    -out server.csr
$ openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in server.csr \
    -out server.crt

# 对客户端证书进行签名
$ openssl req -new \
    -subj "/C=GB/L=China/O=client/CN=client.io" \
    -key client.key \
    -out client.csr
$ openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in client.csr \
    -out client.crt
```

```go
package main
// 服务端ca证书验证
import (

"crypto/tls"
"crypto/x509"
"google.golang.org/grpc"
"google.golang.org/grpc/credentials"
"io/ioutil"
"log"
)

func main()  {
    certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err!=nil{
        log.Fatal(err)
    }
    certpool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("ca.crt")
    if err!=nil{log.Fatal(err)}
    if ok:=certpool.AppendCertsFromPEM(ca);!ok{log.Fatal(err)}
    creds := credentials.NewTLS(&tls.Config{
        Certificates:[]tls.Certificate{certificate},
        ClientAuth:tls.RequireAndVerifyClientCert,
        ClientCAs:certpool,
    })
    server := grpc.NewServer(grpc.Creds(creds))
    log.Println("server", server)
}
```

```go
package main
import (

"crypto/tls"
"crypto/x509"
"google.golang.org/grpc"
"google.golang.org/grpc/credentials"
"io/ioutil"
"log"
)
func main()  {
	certificate, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err != nil {
        log.Fatal(err)
    }

    certPool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        log.Fatal(err)
    }
    if ok := certPool.AppendCertsFromPEM(ca); !ok {
        log.Fatal("failed to append ca certs")
    }

    creds := credentials.NewTLS(&tls.Config{
        Certificates:       []tls.Certificate{certificate},
        ServerName:         "server.io", // NOTE: this is required!
        RootCAs:            certPool,
    })

    conn, err := grpc.Dial(
        "localhost:5000", grpc.WithTransportCredentials(creds),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
}
```