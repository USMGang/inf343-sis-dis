package main

import (
	g "l3/globals"
	"log"
	"math/rand"
)

// ================== Bots ==================
func getBotChoice(n_options int) (int) {
    return rand.Intn(n_options) + 1
}

// ================== Interfaz ==================
func ShowInterface(player int){
    notifications.InitUserChoice()
    notifications.InitInterfaceChoice()

    if (player == 1){
        go func(){
            for {
                // Obtener imput del usuario
                <-get_choice
                notifications.ChangeOptions(g.INTERFACE_PROMPT, g.INTERFACE_OPTIONS)
                notifications.ShowNotifications()
                var option = notifications.GetInterfaceChoice()

                switch option {

                // Continuar mision
                case 1:
                    notifications.ChangeOptions("Escriba su respuesta: ", []string{})
                    get_answer <- true

                // Ver Recompensa
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
}
