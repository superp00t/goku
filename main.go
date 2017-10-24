package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ogier/pflag"
	"github.com/superp00t/goku/control"
)

var (
	rng    = pflag.StringP("range", "r", "192.168.1.0-255", "port range to scan")
	ep     = pflag.StringP("addr", "a", "", "Roku IP address")
	proles = pflag.IntP("workers", "w", 6, "default number of scanner workers")

	handle = map[string]*handler{
		"scan": &handler{
			"scan your network for Roku devices",
			func() {
				ip, err := control.Scan(*rng, *proles)
				if err != nil {
					fail(err)
				}
				fmt.Printf("Roku found at %s!\n", ip)
				ioutil.WriteFile(filePoint(), []byte(ip), 0700)
				fmt.Printf("Address saved to %s.\n", filePoint())
			},
		},

		"shell": &handler{
			"use keyboard as Roku remote",
			func() {
				control.Shell(getIP())
			},
		},
	}
)

func filePoint() string {
	return os.Getenv("HOME") + "/.goku-endpoint"
}

func getIP() string {
	if *ep == "" {
		fmt.Printf("Loading saved IP from %s\n", filePoint())
		b, _ := ioutil.ReadFile(filePoint())
		return string(b)
	}

	return *ep
}

type handler struct {
	Usage string
	Fn    func()
}

func fail(e error) {
	fmt.Fprintf(os.Stderr, "%s\n	", e)
	os.Exit(-1)
}

func main() {
	pflag.Parse()

	if h := handle[pflag.Arg(0)]; h != nil {
		h.Fn()
	} else {
		fmt.Printf("goku usage:\n")
		for k, v := range handle {
			fmt.Printf("\tgoku %s\tusage: %s\n", k, v.Usage)
		}
	}
}
