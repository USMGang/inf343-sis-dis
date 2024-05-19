package namenode_backend

import (
	"fmt"
	g "l3/globals"
	"math/rand"
	"net"
	"os"
	"sync"
    datanode "l3/datanode_backend"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type DataNode struct {
    Id int
    Ip string
    Port string
    // TODO: Si no funca cambiarlo a copia completa y shao
    s *datanode.Server
	lis  net.Listener
    grpcServer *grpc.Server
}

type MercenaryNode struct {
    MercenaryId int
    IdNode int
}

type Server struct {
    UnimplementedNamenodeServiceServer
    mu sync.Mutex
    DataNodes []DataNode
    MercenaryNodes []MercenaryNode
}


func (s *Server) setListener(id int, ip string, port string) {
    ss := s.DataNodes[id]
    ss.s = &datanode.Server{}

    var err error
    ip_conn := fmt.Sprintf("%s:%s", ip, port)
    ss.lis, err = net.Listen("tcp", ip_conn)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establece el listener en: %s", ip_conn))

    ss.grpcServer = grpc.NewServer()

    datanode.RegisterDatanodeServiceServer(ss.grpcServer, ss.s)

    err = ss.grpcServer.Serve(ss.lis)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establecer el servidor gRPC en: %s", ip_conn))
}

func (s *Server) InitServer(ips []string, n_mercenaries int) {

    s.mu.Lock()
    file, err := os.Create("txt/namenode.txt")
    g.FailOnError(err, "Fallo al crear el archivo")
    file.Close()
    s.mu.Unlock()

    ports := []string{ "8071", "8072", "8073" }
    s.DataNodes = make([]DataNode, 3)

    for i := 0; i < 3; i++ {
        s.DataNodes[i] = DataNode{}
        s.DataNodes[i].Id = i
        s.DataNodes[i].Ip = ips[i]
        s.DataNodes[i].Port = ports[i]

        s.setListener(i, s.DataNodes[i].Ip, s.DataNodes[i].Port)
    }

    s.MercenaryNodes = make([]MercenaryNode, n_mercenaries)

    for i := 0; i < n_mercenaries; i++ {
        rand := rand.Intn(3)
        s.MercenaryNodes[i] = MercenaryNode{ i+1, rand }
    }
}

func (s *Server) SaveStep(ctx context.Context, request *SaveStepRequest) (*SaveStepResponse, error) {
    nameNode := s.DataNodes[s.MercenaryNodes[request.Id].IdNode]

    s.mu.Lock()
    file, err := os.OpenFile("txt/namenode.txt", os.O_APPEND|os.O_WRONLY, 0644)
    g.FailOnError(err, "Fallo al abrir el archivo")
    defer file.Close()

    _, err = file.WriteString(fmt.Sprintf("%s Piso_%d %s\n", g.GetName(int(request.Id)), request.Floor, nameNode.Ip))
    g.FailOnError(err, "Fallo al escribir en el archivo")
    s.mu.Unlock()

    _, err = nameNode.s.SaveStep(ctx, &datanode.SaveStepRequest{ Id: request.Id, Floor: request.Floor })
    g.FailOnError(err, "Fallo al enviar el mensaje al datanode")

    return &SaveStepResponse{}, nil
}

func (s *Server) GetIdStepts(ctx context.Context, request *GetIdSteptsRequest) (*GetIdSteptsResponse, error) {
    nameNode := s.DataNodes[s.MercenaryNodes[request.Id].IdNode]

    response, err := nameNode.s.GetIdStepts(ctx, &datanode.GetIdSteptsRequest{ Id: request.Id })
    g.FailOnError(err, "Fallo al enviar el mensaje al datanode")

    return &GetIdSteptsResponse{ Steps: response.Steps }, nil
}


