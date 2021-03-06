package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	timeout time.Duration
	host    string
	port    int
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	// Определяем параметры вызова
	flag.Parse()

	host = flag.Arg(0)
	if host == "" {
		host = "localhost"
	}

	if flag.Arg(1) == "" {
		port = 4242
	} else {
		port, _ = strconv.Atoi(flag.Arg(1))
		if port <= 0 {
			fmt.Println("bad port number")
			os.Exit(1)
		}
	}

	sendBuff := &bytes.Buffer{}
	address := fmt.Sprintf("%s:%d", host, port)
	net.JoinHostPort(host, strconv.Itoa(port))
	client := NewTelnetClient(address, timeout, ioutil.NopCloser(sendBuff), os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	// Читаем строчки из stdin и отправляем на сервер
	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				return
			}
			line := scanner.Text()

			sendBuff.WriteString(line + "\n")
			err := client.Send()
			if err != nil {
				fmt.Fprintln(os.Stderr, "sending connection closed")
				return
			}
		}
	}()

	// Получаем строчки из сервера и печатаем их в stdout
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			err := client.Receive()
			if err != nil {
				fmt.Fprintln(os.Stderr, "receiving connection closed")
				return
			}
		}
	}()

	wg.Wait()

	fmt.Fprintln(os.Stderr, "communication finished")
}
