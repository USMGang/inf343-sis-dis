package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	dosh "l3/doshbank_backend"
	f "l3/floors"
	u "l3/ui"
    g "l3/globals"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
    WIDTH = 150
    N_NOTIFICATIONS = 15
)

var (
    DIRECTOR_PROMPT = "Elige una opción: "
    DIRECTOR_OPTIONS = []string{ "Continuar Mision", "Mercenarios", "Historial",  "Salir" }
)

func main(){
    N_MERCENARIES, _ := strconv.Atoi(os.Args[1])

    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("Fallo al escuchar el puerto 8080: %v", err)
    }

    conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
    g.FailOnError(err, "Error, no se pudo establecer comunicación con el director")

    c := dosh.NewDoshBankClient(conn)

    // Inicializar el server 
    s := f.Server{}

    s.NMercenaries = N_MERCENARIES
    s.CurrentMercenaries = 0
    s.Cond = sync.NewCond(&s.Mutex)
    s.Wait = make(chan bool, s.NMercenaries)
    s.Quit = make(chan struct{})
    quit := make(chan bool)

    s.Ui = u.NewUI(WIDTH, N_NOTIFICATIONS)
    s.Ui.ChangeOptions(DIRECTOR_PROMPT, DIRECTOR_OPTIONS)
    s.Ui.InitInterfaceChoice()

    s.Dosh = dosh.DoshBank{}
    s.Dosh.InitDoshBank()
    defer s.Dosh.Conn.Close()
    defer s.Dosh.Ch.Close()

    // Interfaz del director
    go func(){
        s.Ui.ShowNotifications()
        s.Ui.AddNotification("[Director] Iniciando el la mision...")

        for {
            // Obtener imput del usuario
            var choice = s.Ui.GetInterfaceChoice()

            switch choice {

            // Continuar mision
            case 1:
                if (s.CurrentMercenaries == s.NMercenaries) {
                    s.Cond.Broadcast()
                    s.CurrentMercenaries = 0
                    s.Ui.AddNotification("[Director] Esperando los resultados del piso...")

                    // TODO: Agregar un mutex para que se vuelvan a contar los mercenarios para que cuando la mision termine se ejecute el separador
                    // TODO: Debo incluir esto si quiero que al final de la mison, si se cumple, avise que la mision ha sido completada
                    // s.Ui.AddSeparator()
                } else {
                    s.Ui.AddNotification(fmt.Sprintf("[Director] Esperando a los mercenarios (%d/%d)", s.CurrentMercenaries, s.NMercenaries))
                }

                if (s.NMercenaries == 0) {
                    s.Ui.AddNotification("[Director] La mision ha fallado...")
                }

            // Mercenarios
            case 2:
                request_reward := dosh.GetCurrentRewardRequest{}
                response_reward, err := c.GetCurrentReward(context.Background(), &request_reward)
                g.FailOnError(err, "Error, no se pudieron recibir las opciones")

                s.Ui.AddNotification(fmt.Sprintf("[Director] Recompensa actual: %d", response_reward.Reward))

            // Historial
            case 3:
                s.Ui.ShowNotifications()

            // Salir
            case 4:
                close(quit)

            default:
                log.Fatalf("Error, opción no válida")
            }
        }
    }()

    grpcServer := grpc.NewServer()

    f.RegisterFloorsServiceServer(grpcServer, &s)

    if err := grpcServer.Serve(lis)
    err != nil {
        log.Fatalf("Fallo al ejecutar grcp en el puerto 8080: %v", err)
    }

    go func(){
        <-quit
        lis.Close()
        grpcServer.GracefulStop()
        os.Exit(0)
    }()
}
