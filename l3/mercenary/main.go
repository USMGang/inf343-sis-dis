package main

import (
	"context"
	"fmt"
	f "l3/floors"
	"l3/ui"
	"log"
	"math/rand"
	"os"
	"strconv"
    g "l3/globals"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
    WIDTH = 150
    N_NOTIFICATIONS = 25
)

var (
    INTERFACE_PROMPT = "Seleccione una de las opciones: "
    INTERFACE_OPTIONS = []string{ "Escribir resupesta", "Mercenarios", "Historial" }
)

// TODO: Agregar que al implementar la interfaz, al momento de usar getUserChoice, se espere la opcion y luego vuelva la interfaz
// TODO: VAmos a crear un canal que sea que si se muestra la interfaz o se muestra el prompt que se necesita en el momento
// TODO: O la otra, es que el el switch, hay una opcion que es, ingresar respuesta, y que la pregunta este como notificacion como [Sistema]

func getBotChoice(n_options int) (int) {
    return rand.Intn(n_options) + 1
}

func isReady(id int32, is_ready bool, c f.FloorsServiceClient) bool {
    ready := f.ReadyRequest{ Id: id, IsReady: true }
    response, err := c.MercenaryReady(context.Background(), &ready)
    g.FailOnError(err, "Error, no se pudo recibir el mensaje")

    return(response.Continue)
}

func main(){
    // ================== Inicialización ==================
    // Verificar si el mercenario sera controlado o un bot
    player, err := strconv.Atoi(os.Args[1])
    g.FailOnError(err, "Error, no se pudo recibir el tipo de jugador")

    conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    g.FailOnError(err, "Error, no se pudo establecer comunicación con el director")

    c := f.NewFloorsServiceClient(conn)

    // ================== Interfaz ==================
    // TODO: Implementar que la ui este en una goroutine para ver si podemos hacerla asincrona mientras se ejecuta el programa

    notifications := ui.NewUI(WIDTH, N_NOTIFICATIONS)
    notifications.InitUserChoice()
    notifications.InitInterfaceChoice()
    get_anwser := make(chan bool)
    get_choice := make(chan bool)

    if (player == 1){
        go func(){
            for {
                // Obtener imput del usuario
                <-get_choice
                notifications.ChangeOptions(INTERFACE_PROMPT, INTERFACE_OPTIONS)
                notifications.ShowNotifications()
                var option = notifications.GetInterfaceChoice()

                switch option {

                // Continuar mision
                case 1:
                    notifications.ChangeOptions("Escriba su respuesta: ", []string{})
                    get_anwser <- true

                // Mercenarios
                case 2:
                    notifications.ShowNotifications()

                // Historial
                case 3:
                    notifications.ShowNotifications()

                default:
                    log.Fatalf("Error, opción no válida")
                }
            }
        }()

        get_choice <- true
    } else {
        notifications.ShowNotifications()
    }


    // ================== Preparacion ==================
    // === Obtener id ===
    request_id := f.Start{ Id: -1 }
    response_id, err := c.StartMission(context.Background(), &request_id)
    g.FailOnError(err, "Error, no se pudo recibir la id")

    id := response_id.Id
    notifications.AddNotification(fmt.Sprintf("ID de la mision: %d", id))

    // ================== Piso 1 ==================
    // === Obtener las armas disponibles ===
    request_weapons := f.Floor1WeaponsRequest{}
    response_weapons, err := c.Floor1(context.Background(), &request_weapons)
    g.FailOnError(err, "Error, no se pudieron recibir las opciones")


    // TODO: Hacer que en vez de preguntarle al server, solo printee las opciones disponibles
    weapons := response_weapons.Options
    notifications.AddNotification("Armas disponibles:")
    for i := 0; i < len(weapons); i++ {
        notifications.AddNotification(fmt.Sprintf("- [%d] %s", i+1, weapons[i]))
    }

    // === Seleccionar el arma y obtener el resultado ===
    var choice int
    if (player == 1){
        <-get_anwser
        notifications.AddNotification("Seleccione el arma que desea usar...")
        choice = notifications.GetUserChoice()
        get_choice <- true
    } else {
        choice = getBotChoice(len(weapons))
    }

    notifications.AddNotification(fmt.Sprintf("Arma seleccionada: %s", weapons[choice-1]))

    request_results := f.Floor1ResultsRequest{ Id: id, SelectedWeapon: int32(choice), RandNumber: int32(rand.Intn(101))}
    response_results, err := c.Floor1Results(context.Background(), &request_results)
    g.FailOnError(err, "Error, no se pudo recibir el resultado")

    notifications.AddNotification(response_results.Message)

    // === Verificar si el mercenario murio ===
    if response_results.IsDead { return }

    // === Confirmar que se esta listo para continuar ===
    
    if (player == 1) {
        notifications.AddNotification("Mande 1 si está listo para el siguiente... ")
        <-get_anwser
        notifications.ShowNotifications()
        choice = notifications.GetUserChoice()
        get_choice <- true
        if (choice == 1){
            isReady(int32(id), true, c)
        }
    } else {
        isReady(int32(id), true, c)
    }

    // ================== Piso 2 ==================
    // // === Elegir el camino ===
    // notifications.AddNotification("Entrando al piso 2")
    // notifications.AddNotification("Recorriendo el piso te encuentras con 2 caminos:")
    // notifications.AddNotification("- [1] Camino a la izquierda")
    // notifications.AddNotification("- [2] Camino a la derecha")
    //
    // if (player == 1){
    //     <-get_anwser
    //     notifications.AddNotification("Elige tu camino...")
    //     choice = notifications.GetUserChoice()
    //     get_choice <- true
    // } else {
    //     choice = getBotChoice(2)
    // }
    //
    // // === Ver si el camino es correcto ===
    // notifications.AddNotification(fmt.Sprintf("Camino seleccionado: %d", choice))
    // path_request := f.Floor2PathRequest{ Id: id, SelectedPath: int32(choice) }
    // path_response, err := c.Floor2(context.Background(), &path_request)
    // g.FailOnError(err, "Error, no se pudo recibir el camino")
    //
    // notifications.AddNotification(path_response.Message)
    //
    // // === Verificar si el mercenario fue traicionado ===
    // if (path_response.IsOut) { return }
    //
    // if (player == 1) {
    //     notifications.AddNotification("Mande 1 si está listo para el siguiente... ")
    //     <-get_anwser
    //     notifications.ShowNotifications()
    //     choice = notifications.GetUserChoice()
    //     get_choice <- true
    //
    //     if (choice == 1){
    //         isReady(int32(id), true, c)
    //     }
    // } else {
    //     isReady(int32(id), true, c)
    // }
    //
    // // ================== Piso 3 ==================
    // // === Realizar las 5 rondas ===
    // notifications.AddNotification("Entrando al piso 3")
    // n_good_tries := 0
    // var rand_number int
    // for i := 1; i<=5; i++ {
    //     rand_number = rand.Intn(16)+1
    //     notifications.AddNotification(fmt.Sprintf("Intento %d: %d", i, rand_number))
    //     tries_request := f.Floor3Try{ Id: id, NTries: int32(i), NGoodTries: int32(n_good_tries), RandNumber: int32(rand_number) }
    //     tries_response, err := c.Floor3(context.Background(), &tries_request)
    //     g.FailOnError(err, "Error, no se pudo realizar la ronda")
    //
    //     n_good_tries = int(tries_response.NGoodTries)
    // }
    //
    // // === Revisar el resultado de las rondas ===
    // floor3_request := f.Floor3ResultsRequest{ Id: id, NGoodTries: int32(n_good_tries) }
    // floor3_response, err := c.Floor3Results(context.Background(), &floor3_request)
    // g.FailOnError(err, "Error, no se pudo recibir el resultado")
    //
    // notifications.AddNotification(floor3_response.Message)
    //
    // // === Verificar si el mercenario murio ===
    // if floor3_response.IsDead { return }

}
