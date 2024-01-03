package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Args struct {
	host, port string
	timeout    time.Duration
}

func getArgs() (*Args, error) {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout")

	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		return nil, errors.New("you must specify HOST and PORT")
	}

	return &Args{
		host:    args[0],
		port:    args[1],
		timeout: *timeout,
	}, nil
}

func read(conn net.Conn, cancel context.CancelFunc) {
	scan := bufio.NewScanner(conn)
	for {
		if !scan.Scan() {
			fmt.Printf("Read: connection closed\n")
			cancel()
			return
		}
		text := scan.Text()
		fmt.Printf("%s\n", text)
	}
}

func write(conn net.Conn, cancel context.CancelFunc) {
	scan := bufio.NewScanner(os.Stdin)
	for {
		if !scan.Scan() {
			fmt.Printf("Write: connection closed\n")
			cancel()
			return
		}
		str := scan.Text()

		_, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Printf("Write: can't send message to server\n")
			cancel()
			return
		}
	}
}

func telnet(args *Args) error {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		cancel()
	}()

	connAddress := fmt.Sprintf("%s:%s", args.host, args.port)

	conn, err := net.DialTimeout("tcp", connAddress, args.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	go read(conn, cancel)
	go write(conn, cancel)

	<-ctx.Done()
	fmt.Println("Finish connect")
	return nil
}

func main() {
	args, err := getArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = telnet(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
