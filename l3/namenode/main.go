package main

import (
	"fmt"
	g "l3/globals"
	namenode "l3/namenode_backend"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	n_mercenaries, _ := strconv.Atoi(os.Args[1])

	// ================== Inicializar el servidor ==================
	err := godotenv.Load()
	g.FailOnError(err, "Error al cargar el archivo .env")

	namenodeHost := os.Getenv("NAMENODE_HOST")
	namenodePort := os.Getenv("NAMENODE_PORT")

	ip_conn := fmt.Sprintf("%s:%s", namenodeHost, namenodePort)
	lis, err := net.Listen("tcp", ip_conn)
	g.FailOnError(err, fmt.Sprintf("Error, no se pudo establece el listener en: %s", ip_conn))

	fmt.Println("Servidor gRPC iniciado en: ", ip_conn)

	s := namenode.Server{}
	s.InitServer(n_mercenaries)
	grpcServer := grpc.NewServer()

	namenode.RegisterNamenodeServiceServer(grpcServer, &s)

	err = grpcServer.Serve(lis)
	g.FailOnError(err, fmt.Sprintf("Error, no se pudo establecer el servidor gRPC en: %s", ip_conn))

}
