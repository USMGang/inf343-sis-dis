package globals

import "log"

const (
    WIDTH = 150
    N_NOTIFICATIONS = 15
    REWARD_BONUS = 100000000
)

var (
    VOID_PROMPT = " "
    VOID_OPTIONS = []string{}

    INTERFACE_PROMPT = "Seleccione una de las opciones: "
    INTERFACE_OPTIONS = []string{ "Escribir resupesta", "Mercenarios", "Historial" }
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// TODO: Agregar una funcion para obtener el nombre segun el id del wea
