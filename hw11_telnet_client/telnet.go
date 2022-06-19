package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

const buffSize = 1024

type TelnetClient interface {
	Connect() error
	//io.Closer
	Close() error
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	conn net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	dialer := &net.Dialer{
		Timeout: t.timeout,
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt)
	go func() {
		<-termSignal
		fmt.Fprintln(os.Stderr, "received SIGINT")
		t.conn.Close()
	}()

	conn, err := dialer.DialContext(ctx, "tcp", t.address)
	if err != nil {
		cancel()
		return errors.New("can not establish connection")
	}
	t.conn = conn

	return nil
}

func (t *telnetClient) Close() error {
	return t.conn.Close()
}

func (t *telnetClient) Send() error {
	data := make([]byte, buffSize)
	n, _ := t.in.Read(data)
	data = data[:n]

	_, err := t.conn.Write(data)
	return err
}

func (t *telnetClient) Receive() error {
	data := make([]byte, buffSize)
	n, err := t.conn.Read(data)
	if err != nil {
		return err
	}
	data = data[:n]

	_, err = t.out.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// P.S. Author's solution takes no more than 50 lines.
