package main

import (
	"fmt"
	datanode "l3/datanode_backend"
	g "l3/globals"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

func main(){
    id, err := strconv.Atoi(os.Args[1])
    g.FailOnError(err, "Error al obtener el id del datanode")

    var host, port string
    switch id {
    case 1:
        host = "10.35.169.82"
        port = "8071"
    case 2:
        host = "10.35.169.80"
        port = "8072"
    case 3:
        host = "10.35.169.79"
        port = "8073"
    default:
        host = "10.35.169.82"
        port = "8071"
    }

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
