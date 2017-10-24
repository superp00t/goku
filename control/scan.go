package control

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func Get(url string) (string, *http.Response, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", nil, err
	}

	return string(b), rsp, nil
}

func WorkerScan(ips []string, success chan string, wg *sync.WaitGroup) {
	for _, v := range ips {
		log.Println("Scanning", v)
		_, h, err := Get("http://" + v + ":8060/")
		if err != nil {
			continue
		}

		if strings.Contains(h.Header.Get("Server"), "Roku") {
			success <- v
			return
		}
	}

	wg.Done()
}

func Scan(cidr string, procs int) (string, error) {
	ips, err := Hosts(cidr)
	if err != nil {
		return "", err
	}

	str := make([][]string, procs)

	// 8 ips 7 procs
	for i := 0; i < len(ips); i += procs {
		for ie := 0; ie < procs; ie++ {
			ind := i + ie
			if ind > len(ips)-1 {
				break
			}
			ip := ips[ind]
			str[ie] = append(str[ie], ip)
		}
	}
	// d, _ := json.MarshalIndent(str, "", "\t")
	// log.Fatalf("%s", d)
	timeout := make(chan bool)

	wg := new(sync.WaitGroup)
	success := make(chan string)

	for _, v := range str {
		wg.Add(1)
		go WorkerScan(v, success, wg)
	}

	go func() {
		wg.Wait()
		timeout <- true
	}()

	select {
	case s := <-success:
		return s, nil
	case <-timeout:
		return "", fmt.Errorf("No Roku devices detected.")
	}
}

func parseInt(i string) int64 {
	ie, e := strconv.ParseInt(i, 10, 64)
	if e != nil {
		return -1
	}
	return ie
}

func Hosts(p string) ([]string, error) {
	ips := []string{}
	h := strings.Split(p, ".")

	if strings.Contains(p, "-") == false {
		return []string{p}, nil
	}

	for i, v := range h {
		t := strings.Split(v, "-")

		if len(t) == 2 {
			low := parseInt(t[0])
			high := parseInt(t[1])

			for x := low; x < high+1; x++ {
				inst := h
				inst[i] = fmt.Sprintf("%d", x)
				ips = append(ips, strings.Join(inst, "."))
			}

			if low < 0 || high < 0 {
				return nil, fmt.Errorf("invalid IP range")
			}
		}
	}

	return ips, nil
}
