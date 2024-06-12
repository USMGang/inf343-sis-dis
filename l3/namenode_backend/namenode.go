package namenode_backend

import (
	"fmt"
	datanode "l3/datanode_backend"
	g "l3/globals"
	"math/rand"
	"os"
	"sync"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataNode struct {
	Id   int
	Ip   string
	Port string
	// TODO: Si no funca cambiarlo a copia completa y shao
	DataNodeConn *grpc.ClientConn
	DataNode     datanode.DatanodeServiceClient
}

type MercenaryNode struct {
	MercenaryId int
	IdNode      int
}

type Server struct {
	UnimplementedNamenodeServiceServer
	mu             sync.Mutex
	DataNodes      []DataNode
	MercenaryNodes []MercenaryNode
}

func (s *Server) setClient(id int, ip string, port string) {

	ip_conn := fmt.Sprintf("%s:%s", ip, port)
	conn, err := grpc.NewClient(ip_conn, grpc.WithTransportCredentials(insecure.NewCredentials()))

	g.FailOnError(err, "Error, no se pudo establecer comunicación con el director")

    s.DataNodes[id].DataNodeConn = conn
    s.DataNodes[id].DataNode = datanode.NewDatanodeServiceClient(conn)

    fmt.Println("Conexión establecida con el datanode", ip, port)
}

func (s *Server) InitServer(n_mercenaries int) {

	s.mu.Lock()
	file, err := os.Create("txt/namenode.txt")
	g.FailOnError(err, "Fallo al crear el archivo")
	file.Close()
	s.mu.Unlock()

	ips := make([]string, 3)
	ports := make([]string, 3)

	for i := 0; i < 3; i++ {
		switch i+1 {
		case 1:
			ips[i] = "10.35.169.82"
			ports[i] = "8071"
		case 2:
			ips[i] = "10.35.169.80"
			ports[i] = "8072"
		case 3:
			ips[i] = "10.35.169.79"
			ports[i] = "8073"
		default:
			ips[i] = "10.35.169.82"
			ports[i] = "8071"
		}
	}

	s.DataNodes = make([]DataNode, 3)

	for i := 0; i < 3; i++ {
		s.DataNodes[i] = DataNode{}
		s.DataNodes[i].Id = i+1
		s.DataNodes[i].Ip = ips[i]
		s.DataNodes[i].Port = ports[i]

		s.setClient(i, s.DataNodes[i].Ip, s.DataNodes[i].Port)
	}

	s.MercenaryNodes = make([]MercenaryNode, n_mercenaries)

	for i := 0; i < n_mercenaries; i++ {
		rand := rand.Intn(3)
		s.MercenaryNodes[i] = MercenaryNode{i + 1, rand}
	}
}

func (s *Server) SaveStep(ctx context.Context, request *SaveStepRequest) (*SaveStepResponse, error) {
	nameNode := s.DataNodes[s.MercenaryNodes[request.Id-1].IdNode]

	s.mu.Lock()
	file, err := os.OpenFile("txt/namenode.txt", os.O_APPEND|os.O_WRONLY, 0644)
	g.FailOnError(err, "Fallo al abrir el archivo")
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s Piso_%d %s\n", g.GetName(int(request.Id)), request.Floor, nameNode.Ip))
	g.FailOnError(err, "Fallo al escribir en el archivo")
	s.mu.Unlock()

    _, err = nameNode.DataNode.SaveStep(ctx, &datanode.SaveStepRequest{Id: request.Id, Floor: request.Floor, Steps: request.Steps})
	g.FailOnError(err, "Fallo al guardar el mensaje al datanode")

	return &SaveStepResponse{}, nil
}

func (s *Server) GetIdStepts(ctx context.Context, request *GetIdSteptsRequest) (*GetIdSteptsResponse, error) {
	nameNode := s.DataNodes[s.MercenaryNodes[request.Id-1].IdNode]

	response, err := nameNode.DataNode.GetIdStepts(ctx, &datanode.GetIdSteptsRequest{Id: request.Id, Floor: request.Floor})
	g.FailOnError(err, "Fallo al recibir el mensaje al datanode")

	return &GetIdSteptsResponse{Steps: response.Steps}, nil
}
