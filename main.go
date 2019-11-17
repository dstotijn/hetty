package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/google/martian"
	"github.com/google/martian/mitm"
)

func main() {
	p := martian.NewProxy()
	defer p.Close()

	tlsc, err := tls.LoadX509KeyPair("/Users/dstotijn/.ssh/gurp_cert.pem", "/Users/dstotijn/.ssh/gurp_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	priv := tlsc.PrivateKey

	x509c, err := x509.ParseCertificate(tlsc.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}

	mc, err := mitm.NewConfig(x509c, priv)
	if err != nil {
		log.Fatal(err)
	}
	mc.SetValidity(time.Hour)
	mc.SetOrganization("Gurp, Inc.")

	p.SetMITM(mc)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := p.Serve(l); err != nil {
			log.Println(err)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)
	<-sigc
}
