package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const crlf = "\r\n"

var (
	hostname = flag.String("host", "192.168.178.31", "address of the light")
	commands map[string]string
)

func init() {
	commands = map[string]string{
		"on":     `{"id": 1, "method": "set_power",  "params":["on", "smooth", 3000]}`,
		"off":    `{"id": 1, "method": "set_power",  "params":["off", "smooth", 3000]}`,
		"toggle": `{"id": 1, "method": "toggle",     "params":[]}`,

		"warm": colourTemperatureCommand(2700),
		"cold": colourTemperatureCommand(6500),

		"high": brightnessCommand(100),
		"low":  brightnessCommand(1),

		"get": `{"id": 1, "method": "get_prop",  "params":["power", "bright", "ct", "flowing","name"]}`,

		"flow": `{"id": 1, "method":"stop_cf", "params":[]}`, // start_cf's look like flow:count:action:args

		"blink": `{"id":1,"method":"start_cf","params":[1, 0, "50,2,1000,1"]}`,
	}
}

func brightnessCommand(brightness int) string {
	return fmt.Sprintf(`{"id": 1, "method": "set_bright", "params":[%v, "sudden", 0]}`, brightness)
}

func colourTemperatureCommand(temperature int) string {
	return fmt.Sprintf(`{"id": 1, "method": "set_ct_abx", "params":[%v, "sudden", 0]}`, temperature)
}

func printHelp() {
	fmt.Print(`

help page goes here


`)

	for k, v := range commands {
		fmt.Printf("    %10v: %v\n", k, v)
	}
	fmt.Println()
}

func main() {
	flag.Parse()

	hostAndPort := *hostname + ":55443"

	conn, err := net.Dial("tcp", hostAndPort)
	if err != nil {
		log.Fatal(err)
	}

	// Copy connection output to stdout.
	go func() {
		io.Copy(os.Stdout, conn)
	}()

	// Sending too soon after connection open seems to confuse my lamp.
	time.Sleep(300 * time.Millisecond)

	args := flag.Args()

	if len(args) == 0 {
		printHelp()
		return
	}

	for _, arg := range args {
		if arg == "," {
			time.Sleep(time.Second)
			continue
		}

		if arg == "listen" {
			fmt.Println("listening, type ^C to exit")
			<-make(chan bool) // block forver
		}

		cmd := commands[arg]

		if cmd == "" && strings.HasPrefix(arg, "flow:") {
			options := strings.Split(arg, ":")
			if len(options) != 4 {
				log.Fatal("flow with args must be flow:count:action:args")
			}
			cmd = fmt.Sprintf(`{"id":1,"method":"start_cf","params":[%v, %v, "%v"]}`, options[1], options[2], options[3])
		}

		if cmd == "" && strings.HasPrefix(arg, "rest:") {
			if d, err := time.ParseDuration(arg[5:]); err != nil {
				log.Fatalf("can't parse duration from %v : %v", arg, err)
			} else {
				fmt.Printf("resting for %v\n", d)
				time.Sleep(d)
				continue
			}
		}

		if cmd == "" {
			if numeric, err := strconv.Atoi(arg); err == nil {
				if numeric == 0 {
					cmd = commands["off"]
				} else if numeric > 1000 {
					cmd = colourTemperatureCommand(numeric)
				} else {
					cmd = brightnessCommand(numeric)
				}
			}
		}

		if cmd == "" {
			fmt.Printf("unknown option: %v\n", arg)
			printHelp()
			return
		}

		fmt.Printf("%v <- %v\n", hostAndPort, cmd)
		conn.Write([]byte(cmd + crlf))

		// Talking too fast seems to confuse it and 1/second is noted in the spec docs as a limit.
		time.Sleep(200 * time.Millisecond)
	}

}
