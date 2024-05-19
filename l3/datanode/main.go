package main

import (
	"fmt"
	g "l3/globals"
    datanode "l3/datanode_backend"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

func main(){
    id, _ := strconv.Atoi(os.Args[1])
    ip_conn := fmt.Sprintf("%s:%s", os.Args[2], "8070")

    // ================== Inicializar el servidor ==================
    lis, err := net.Listen("tcp", ip_conn)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establece el listener en: %s", ip_conn))

    s := datanode.Server{}
    s.Id = id

    grpcServer := grpc.NewServer()

    fmt.Printf("Servidor gRPC del Datanode%d, iniciado en: %s", id, ip_conn)

    datanode.RegisterDatanodeServiceServer(grpcServer, &s)

    err = grpcServer.Serve(lis)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establecer el servidor gRPC en: %s", ip_conn))

}
