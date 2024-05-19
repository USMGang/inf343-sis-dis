package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
    g "l3/globals"
)

type UI struct {
	Width             int
	maxNotifications  int
	notifications     []string
	visibleNotifications []string
    prompt            string
	options           []string
	topLeft           rune
	topRight          rune
	bottomLeft        rune
	bottomRight       rune
	horizontal        rune
	vertical          rune
    separator         string
    reader            *bufio.Reader
    ChoiceChan        chan bool
    AnswerChan        chan bool
    Choice            int
    Answer            int
    mu                sync.Mutex
    muPrint           sync.Mutex
}

// ========= Constructor ==========
func NewUI(maxNotifications int) UI {
    separator := fmt.Sprintf("%s%s%s", 
        strings.Repeat(string(" "), 2), 
        strings.Repeat(string("="), g.WIDTH-2-2-2), 
        strings.Repeat(string(" "), 2), 
    )

	return UI{
		Width:             150,
        prompt:            "",
		options:           []string{},
		maxNotifications:  maxNotifications,
		topLeft:           '┏',
		topRight:          '┓',
		bottomLeft:        '┗',
		bottomRight:       '┛',
		horizontal:        '━',
		vertical:          '┃',
        separator:         separator,
        reader:            bufio.NewReader(os.Stdin),
        AnswerChan:        make(chan bool),
        ChoiceChan:        make(chan bool),
        Choice:            0,
	}
}

// ========== Notificaciones ==========
func (u *UI) AddNotification(notification string) {
    u.mu.Lock()
	u.notifications = append(u.notifications, notification)
	u.visibleNotifications = append(u.visibleNotifications, notification)

	if len(u.visibleNotifications) > u.maxNotifications {
		u.visibleNotifications = u.visibleNotifications[1:]
	}

	u.ShowNotifications()
    u.mu.Unlock()
}

func (u *UI) AddSeparator(){
    u.AddNotification(u.separator)
}

// ========== Interfaz ==========
func (u *UI) InitInterfaceChoice() {
    go func(){
        for {
            <-u.ChoiceChan
            input, err := u.reader.ReadString('\n')
            if err != nil {
                log.Fatalf("Error, no se pudo leer la entrada: %s", err)
            }

            input = strings.TrimSpace(input)
            aux, _ := strconv.Atoi(input)
            u.Choice = aux
            u.ChoiceChan <- true
        }
    }()
}

func (u *UI) InitUserChoice() {
    go func(){
        for {
            <-u.AnswerChan
            input, err := u.reader.ReadString('\n')
            if err != nil {
                log.Fatalf("Error, no se pudo leer la entrada: %s", err)
            }

            input = strings.TrimSpace(input)
            aux, _ := strconv.Atoi(input)
            u.Answer = aux
            u.AnswerChan <- true
        }
    }()
}

func (u *UI) GetInterfaceChoice() int {
    u.ChoiceChan <- true
    <-u.ChoiceChan
    return u.Choice
}

func (u *UI) GetUserChoice() int {
    u.AnswerChan <- true
    <-u.AnswerChan
    return u.Answer
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (u *UI) PrintNotifications(height int, notifications []string) {
	clearScreen()

	// Printear la parte superior
	fmt.Printf("%s%s%s\n", string(u.topLeft), strings.Repeat(string(u.horizontal), u.Width-2), string(u.topRight))

	// Printear las notificaciones
	nNotifications := len(notifications)
	for i := 0; i < nNotifications; i++ {
        nMessage := utf8.RuneCountInString(notifications[i])

		fmt.Print(string(u.vertical))
		if nMessage != 0 {
			fmt.Print(notifications[i])
			fmt.Printf("%s", strings.Repeat(" ", u.Width-2-nMessage))
		} else {
			fmt.Printf("%s", strings.Repeat(" ", u.Width-2))
		}
		fmt.Println(string(u.vertical))
	}

	for i := 0; i < height-2-nNotifications; i++ {
		fmt.Printf("%s%s%s\n", string(u.vertical), strings.Repeat(" ", u.Width-2), string(u.vertical))
	}

	// Printear la parte inferior
	fmt.Printf("%s%s%s\n", string(u.bottomLeft), strings.Repeat(string(u.horizontal), u.Width-2), string(u.bottomRight))
}

func (u *UI) ChangeOptions(prompt string, options []string) {
    u.prompt = prompt
	u.options = options
}

func (u *UI) ShowButtons() {
	for i := 1; i <= len(u.options); i++ {
		fmt.Printf("[%d] %s\n", i, u.options[i-1])
	}
}

func (u *UI) ShowNotifications() {
    u.muPrint.Lock()
    u.PrintNotifications(u.maxNotifications, u.visibleNotifications)
	u.ShowButtons()
    fmt.Printf("\n%s", u.prompt)
    u.muPrint.Unlock()
}

func (u *UI) ShowAllNotifications() {
    u.muPrint.Lock()
	u.PrintNotifications(len(u.notifications), u.notifications)
    u.muPrint.Unlock()
}
