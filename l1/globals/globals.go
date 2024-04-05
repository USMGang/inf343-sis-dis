package globals

type Message struct {
    Tresure int `json:"tresure"`
    Captain string `json:"captain"`
    Planet string `json:"planet"`
}

var N_leters int = 6
var Max_tresures int = 15
var Base_leter byte = byte('A')
var Task_counter int = 5
