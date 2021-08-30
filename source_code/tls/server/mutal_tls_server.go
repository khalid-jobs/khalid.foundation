package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	//cert, err := tls.LoadX509KeyPair("./server/kubeedge/server.crt", "./server/kubeedge/server.key")
	//cert, err := tls.LoadX509KeyPair("./server/server.crt", "./server/server.key")
	cert, err := tls.LoadX509KeyPair("./server/ecdsa/server.crt", "./server/ecdsa/server.key")
	if err != nil {
		log.Println(err)
		return
	}

	//caCertBytes, err := ioutil.ReadFile("./ca-kubeedge.crt")
	caCertBytes, err := ioutil.ReadFile("./ca.crt")
	if err != nil {
		panic("unable to read client.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(caCertBytes)
	if !ok {
		panic("failed to parse root certificate")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: clientCertPool,
	}

	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println("err")
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn1(conn)
	}
}

func handleConn1(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
