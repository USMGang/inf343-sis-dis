package main

import (
	"fmt"
	d "l3/doshbank_backend"
	f "l3/floors"
	g "l3/globals"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FloorsServers struct {
	// ========== gRPC listener ==========
	lis  net.Listener
    grpcServer *grpc.Server

	// ========== RabbitMQ producer ==========
	DoshbankConn *grpc.ClientConn
	Doshbank     d.DoshBankClient
}

func (s *FloorsServers) setListener(ip string, port string, fs f.FloorsServiceServer) {
    var err error
    ip_conn := fmt.Sprintf("%s:%s", ip, port)
    s.lis, err = net.Listen("tcp", ip_conn)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establece el listener en: %s", ip_conn))

    s.grpcServer = grpc.NewServer()

    f.RegisterFloorsServiceServer(s.grpcServer, fs)

    err = s.grpcServer.Serve(s.lis)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establecer el servidor gRPC en: %s", ip_conn))
}

func (s *FloorsServers) initDoshbank(ip string, port string) {
	ip_conn := fmt.Sprintf("%s:%s", ip, port)
	conn, err := grpc.NewClient(ip_conn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	g.FailOnError(err, "Error, no se pudo establecer comunicación con el director")

	s.DoshbankConn = conn
	s.Doshbank = d.NewDoshBankClient(conn)
}
