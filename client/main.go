package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	gen "github.com/SandyWalsh/runner/sprinter/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TLSConfig() *tls.Config {
	clientCert, err := tls.LoadX509KeyPair("../cert/client.cert", "../cert/client.key")
	if err != nil {
		log.Fatal("failed to load client certs:", err)
		os.Exit(1)
	}

	b, _ := ioutil.ReadFile("../cert/server.cert")
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		log.Fatal("failed to append server certificate")
		os.Exit(1)
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            cp,
	}
	return config
}

func main() {
	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	creds := credentials.NewTLS(TLSConfig())
	conn, err := grpc.DialContext(ctx, "localhost:9990", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gen.NewRunnerClient(conn)

	lreq := gen.LaunchRequest{Cmd: "ls", Arg: []string{"/", "-la"}}
	process, err := client.Run(ctx, &lreq)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("JobId=", process.JobId)
	}

	gr := gen.GeneralRequest{JobId: process.JobId}

	status, err := client.GetStatus(ctx, &gr)
	text := status.Status.String()
	log.Println("status:", text, "exit code:", status.ExitCode)

	/*
		status, err = client.Abort(ctx, &gr)
		text = status.Status.String()
		log.Println("abort status:", text, "exit code:", status.ExitCode)
	*/
	stream, err := client.StreamOutput(ctx, &gr)
	if err != nil {
		log.Fatal("fatal error from StreamOutput:", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("Error on Recv()", err)
			}
		} else {
			fmt.Printf(string(res.Data))
		}
	}

}
