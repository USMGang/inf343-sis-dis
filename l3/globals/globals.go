package globals

import "log"

var (
    WIDTH = 150
    N_NOTIFICATIONS = 15
    REWARD_BONUS = 100000000

    VOID_PROMPT = " "
    VOID_OPTIONS = []string{}

    INTERFACE_PROMPT = "Seleccione una de las opciones: "
    INTERFACE_OPTIONS = []string{ "Escribir resupesta", "Mercenarios", "Historial" }

    DIRECTOR_PROMPT = "Elige una opción: "
    DIRECTOR_OPTIONS = []string{ "Continuar Mision", "Mercenarios", "Historial",  "Salir" }

    MERCENARIES_NAMES = []string{"Reverend David", "Ana Larive", "D.A.R", `Donovan "Bulldog" Neal`, "Mr.Foster", "Hayato Tanaka", "Rae Higgens" }

)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetName(id int) string {
    return MERCENARIES_NAMES[id-1]
}
