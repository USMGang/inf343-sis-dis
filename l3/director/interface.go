package main

import (
	"context"
	"fmt"
	dosh "l3/doshbank_backend"
	g "l3/globals"
	"log"
)

func showInterface(){
    go func(){
        s.Ui.ShowNotifications()
        s.Ui.AddNotification("[Director] Iniciando el la mision...")

        for {
            // Obtener imput del usuario
            choice = s.Ui.GetInterfaceChoice()

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
                response_reward, err := s.Dosh.GetCurrentReward(context.Background(), &request_reward)
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
}
