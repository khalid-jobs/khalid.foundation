package main

import (
	"context"
	"encoding/pem"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	//libp2ptls "github.com/libp2p/go-libp2p-tls"
	libp2ptlsca "khalid.fondation/libp2pdemo/go-libp2p-tls-ca"
	"log"
	"os"
)

func main() {
	// openssl genrsa -out rsa_private.key 2048
	certBytes, err := ioutil.ReadFile("D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\client\\client.key")
	if err != nil {
		log.Println("unable to read client.pem, error: ", err)
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalECDSAPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	os.Setenv("CERTFILE", "D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\client\\client.crt")
	os.Setenv("KEYFILE", "D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\client\\client.key")
	os.Setenv("CAFILE", "D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\ca-ca-server.crt")

	// 这里使用libp2p.NoSecurity，如果对端不使用这个或者Security里面的协议不一样的化，两者就无法建立连接
	// 因为没有对应的协议做Secrutiy
	node, err := libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.Ping(false),
		//libp2p.NoSecurity)
		//libp2p.Security(libp2ptls.ID, libp2ptls.New))
		libp2p.Security(libp2ptlsca.ID, libp2ptlsca.New))
	if err != nil {
		panic(err)
	}

	//if len(os.Args) <= 1 {
	//	panic("Please provide the peer addr")
	//}
	//addr, err := multiaddr.NewMultiaddr(os.Args[1])
	//addrStr := "/ip4/192.168.0.38/tcp/10001/p2p/QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
	addrStr := "/ip4/127.0.0.1/tcp/10001/p2p/QmZaB9gz9Vhuc3Gc1mz1tAXUCqkfsKZXieQiEiUk57xQiF"
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	peer, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := node.Connect(ctx, *peer); err != nil {
		panic(err)
	}

	s, err := node.NewStream(ctx, peer.ID, "/echo/1.0.0")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Sender syaning hello")
	_, err = s.Write([]byte("Hello, world\n"))

	if err != nil {
		log.Println(err)
		return
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("read reply: %q\n", out)

	if err := node.Close(); err != nil {
		panic(err)
	}
}
