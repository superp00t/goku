package control

import (
	"fmt"
	"os"

	"github.com/pkg/term"
)

var (
	t *term.Term
)

const (
	SIGQUIT uint8 = 3
	ENTER   uint8 = 13
)

func SendChar(ip, key string) {
	fmt.Fprint(t, "Sending key ", key)
	fmt.Fprintf(t, "\r\n")
	go func() {
		Get("http://" + ip + ":8060/keypress/Lit_" + key)
	}()
}

func Shell(ip string) {
	t, _ = term.Open("/dev/tty")
	term.RawMode(t)
	for {
		b := make([]byte, 4)
		os.Stdin.Read(b)
		key := b[0]

		// directional key
		if b[1] == 91 {
			switch b[3] {
			case 65:
				SendChar(ip, "Up")
			case 66:
				SendChar(ip, "Down")
			case 68:
				SendChar(ip, "Left")
			case 67:
				SendChar(ip, "Right")
			}
			continue
		}

		switch key {
		case ENTER:
			SendChar(ip, "Enter")
		case SIGQUIT:
			t.Restore()
			t.Close()
			return
		default:
			k := fmt.Sprintf("%c", key)
			SendChar(ip, k)
		}
	}
}
