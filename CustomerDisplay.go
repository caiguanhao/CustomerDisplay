package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"go.bug.st/serial.v1"
	"go.bug.st/serial.v1/enumerator"
)

func openPort(name string) {
	mode := &serial.Mode{
		BaudRate: 2400,
	}
	port, err := serial.Open(name, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	if listen {
		buff := make([]byte, 20)
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
				break
			}
			if n == 0 {
				log.Println("EOF")
				break
			}
			re := regexp.MustCompile("[0-9.]+")
			fmt.Println(re.FindString(string(buff[:n])))
		}
		return
	}

	for {
		rand.Seed(time.Now().UnixNano())
		num := fmt.Sprintf("%.2f", (rand.Float64()*899)+100)
		_, err = port.Write([]byte("\x18\x1b\x41" + num + "\x0d"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(num)
		time.Sleep(1 * time.Second)
	}
}

var listen bool

func main() {
	var name = flag.String("port", "", "port name")
	flag.BoolVar(&listen, "listen", false, "listen")
	flag.Parse()
	if len(*name) > 0 {
		openPort(*name)
		return
	}
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
		return
	}
	for _, port := range ports {
		if port.IsUSB && port.VID == "067B" && port.PID == "2303" {
			log.Println("Opening port", port.Name)
			openPort(port.Name)
		}
	}
}
