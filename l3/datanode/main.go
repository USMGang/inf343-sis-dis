package main

import (
	"fmt"
	datanode "l3/datanode_backend"
	g "l3/globals"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main(){
    id, err := strconv.Atoi(os.Args[1])
    g.FailOnError(err, "Error al obtener el id del datanode")

    err = godotenv.Load()
    g.FailOnError(err, "Error al cargar el archivo .env")
    host := os.Getenv(fmt.Sprintf("DATANODE%d_HOST", id))
    port := os.Getenv(fmt.Sprintf("DATANODE%d_PORT", id))
    ip_conn := fmt.Sprintf("%s:%s", host, port)

    // ================== Inicializar el servidor ==================

    lis, err := net.Listen("tcp", ip_conn)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establece el listener en: %s", ip_conn))

    s := datanode.Server{}
    s.Id = id

    grpcServer := grpc.NewServer()

    fmt.Printf("Servidor gRPC del Datanode%d, iniciado en: %s\n", id, ip_conn)

    datanode.RegisterDatanodeServiceServer(grpcServer, &s)

    err = grpcServer.Serve(lis)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establecer el servidor gRPC en: %s", ip_conn))

}
