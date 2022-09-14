package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"syscall"

	"github.com/SandyWalsh/runner/sprinter"
	gen "github.com/SandyWalsh/runner/sprinter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func TLSServer(caCert, pemPath, keyPath string) (*grpc.Server, error) {
	pemClientCA, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	serverCert, err := tls.LoadX509KeyPair(pemPath, keyPath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	cred := credentials.NewTLS(config)

	s := grpc.NewServer(
		grpc.Creds(cred),
	)
	return s, nil
}

type Server struct {
	rlib sprinter.Runner
	gen.UnimplementedRunnerServer
}

var _ gen.RunnerServer = (*Server)(nil)

func extractUID(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if ok {
		tlsInfo := pr.AuthInfo.(credentials.TLSInfo)
		for _, x := range tlsInfo.State.VerifiedChains {
			for _, y := range x {
				for _, n := range y.Subject.Names {
					if n.Type.String() == "0.9.2342.19200300.100.1.1" {
						if uid, ok := n.Value.(string); ok {
							return uid, nil
						}
					}
				}
			}
		}
	}
	return "", errors.New("no UID field")
}

func (s *Server) Run(ctx context.Context, r *gen.LaunchRequest) (*gen.Process, error) {
	uid, err := extractUID(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("UID=", uid)
	p, err := s.rlib.Run(ctx, uid, r.Cmd, r.Arg...)
	if err != nil {
		return nil, err
	}
	return &gen.Process{JobId: string(p)}, nil
}

func sprinterStatusToGRPC(status sprinter.Status) gen.Status_StatusEnum {
	// Icky ... a simple little data structure could map this much cleaner.
	switch status {
	case sprinter.Aborted:
		return gen.Status_ABORTED
	case sprinter.Running:
		return gen.Status_RUNNING
	case sprinter.Completed:
		return gen.Status_COMPLETED
	}
	return gen.Status_ERROR
}

func (s *Server) GetStatus(ctx context.Context, r *gen.GeneralRequest) (*gen.Status, error) {
	status, ec, err := s.rlib.GetStatus(ctx, sprinter.Process(r.JobId))
	if err != nil {
		return nil, err
	}

	return &gen.Status{Status: sprinterStatusToGRPC(status), ExitCode: int32(ec)}, nil
}

func (s *Server) StreamOutput(r *gen.GeneralRequest, ss gen.Runner_StreamOutputServer) error {
	var sb sprinter.SafeBuffer

	doneWriting := make(chan bool)

	go func() {
		reader, data := sb.NewReader()
		for {
			select {
			case <-ss.Context().Done(): // caller disconnected
				goto done
			case x := <-data:
				if x == 0 {
					goto done
				}
				b := make([]byte, x)
				n, rerr := reader.Read(b)
				if n == 0 {
					if rerr == io.EOF {
						goto done
					}
					if rerr != nil {
						log.Println("error reading from main stream:", rerr)
						goto done
					}
				}
				if n > 0 {
					ss.Send(&gen.Output{Data: b[:n]})
				}
			}
		}
	done:
		doneWriting <- true
	}()

	// blocks until streaming ends
	err := s.rlib.StreamOutput(ss.Context(), sprinter.Process(r.JobId), &sb)
	if err != nil {
		log.Fatal("stdout returning error:", err)
	}

	// make sure we're done writing all the data before closing the stream
	select {
	case <-ss.Context().Done():
		break
	case <-doneWriting:
		break
	}

	return err
}

func (s *Server) Abort(ctx context.Context, r *gen.GeneralRequest) (*gen.Status, error) {
	st, err := s.rlib.Abort(ctx, sprinter.Process(r.JobId))
	if err != nil {
		return nil, err
	}
	return &gen.Status{Status: sprinterStatusToGRPC(st)}, nil
}

func main() {
	tls, err := TLSServer("../cert/ca.cert", "../cert/server.cert", "../cert/server.key")
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", "127.0.0.1:9990")
	if err != nil {
		log.Fatal(err)
	}

	// Note: these should come from env vars.
	// we'll use the "nobody" account. May vary from distro to distro.
	//var uid uint32 = 65534
	//var gid uint32 = 65534

	// These authz rules could come from a local file.
	// NOTE: for now we're just limiting on cpu, but the same mechanism could be used for any
	// resource controller.
	// cpu.weight isn't a great parameter since it's based on the number of running processes in the group (always 1).
	// It's here for illustrative purposes.

	// NOTE: we should query available resources before setting these limits (such as memory.current, /proc/partitions for disk io)
	sta := "259:0" // NOTE: Your disk will vary
	cgs := map[string]sprinter.ControlGroup{
		"good": {
			Limits: []sprinter.Limit{
				{Var: "cpuset.cpus", Value: "1"},
				{Var: "cpu.max", Value: "100000 1000000"},
				{Var: "cpu.weight", Value: "50"},
				{Var: "memory.max", Value: "250M"},
				{Var: "io.max", Value: fmt.Sprintf("%s rbps=2097152 wiops=120", sta)},
			},
			SysProcAttr: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
				//Credential: &syscall.Credential{Uid: uid, Gid: gid},
			}},
		"better": {
			Limits: []sprinter.Limit{
				{Var: "cpuset.cpus", Value: "1"},
				{Var: "cpu.max", Value: "250000 1000000"},
				{Var: "cpu.weight", Value: "100"},
				{Var: "memory.max", Value: "500M"},
				{Var: "io.max", Value: fmt.Sprintf("%s 16:32 rbps=4194304 wiops=240", sta)},
			},
			SysProcAttr: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
				//Credential: &syscall.Credential{Uid: uid, Gid: gid},
			}},
		"best": {
			Limits: []sprinter.Limit{
				{Var: "cpuset.cpus", Value: "1"},
				{Var: "cpu.max", Value: "500000 1000000"},
				{Var: "cpu.weight", Value: "150"},
				{Var: "memory.max", Value: "1G"},
				{Var: "io.max", Value: fmt.Sprintf("%s 32:64 rbps=8388608 wiops=360", sta)},
			},
			SysProcAttr: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
				//Credential: &syscall.Credential{Uid: uid, Gid: gid},
			}}}
	authz := sprinter.AuthZRules{
		ControlGroups:  cgs,
		ClientToCGroup: map[string]string{"cert-A": "good", "cert-B": "better", "cert-C": "best"},
	}

	lr, err := sprinter.NewRunner(authz, &sprinter.LinuxDriver{})
	if err != nil {
		log.Fatalln("unable to initialize sprinter library:", err)
	}
	s := &Server{rlib: lr}
	gen.RegisterRunnerServer(tls, s)
	log.Println("starting Runner ...")
	log.Fatal(tls.Serve(lis))
}
