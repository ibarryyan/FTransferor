package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

const (
	DefaultPathDir    = "tmp"
	DefaultServerPort = 8088
	DefaultWebPort    = 8089
)

var (
	port, webport int
	path, secret  string
)

func ServerCmd() *cobra.Command {
	command := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

			initCommand()

			go func() {
				runHttpServer(webport)
			}()

			go func() {
				runTcpServer(port, path)
			}()

			<-quit
		},
	}

	command.Flags().StringVarP(&path, "path", "", "", "File save path")
	command.Flags().StringVarP(&secret, "secret", "", "", "Http transmission key")
	command.Flags().IntVarP(&port, "port", "", 0, "Tcp server port")
	command.Flags().IntVarP(&webport, "webport", "", 0, "Http server port")
	return command
}

func initCommand() {
	FetchDeviceInfo()
	if path == "" {
		path = DefaultPathDir
	}
	if port == 0 {
		port = DefaultServerPort
	}
	if webport == 0 {
		webport = DefaultWebPort
	}
	_ = os.Mkdir(fmt.Sprintf("%s/", path), os.ModePerm)
}

func runTcpServer(port int, savePath string) {
	fn := "runTcpServer"
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("%s starting server err:%s \n", fn, err)
		return
	}

	defer func() {
		_ = listener.Close()
	}()

	fmt.Printf("%s is listening on port %d \n", fn, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s connection error:%s \n", fn, err)
			continue
		}
		go handleConnection(conn, savePath)
	}
}

func handleConnection(conn net.Conn, savePath string) {
	fn := "handleConnection"
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)

	meta, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("%s reading file metadata:%s \n", fn, err)
		return
	}

	meta = strings.TrimSpace(meta)
	parts := strings.Split(meta, "|")
	if len(parts) != 2 {
		fmt.Printf("%s invalid metadata received", fn)
		return
	}

	fileName := parts[0]
	fileSize := 0
	if _, err = fmt.Sscanf(parts[1], "%d", &fileSize); err != nil {
		fmt.Printf("%s parsing file size err:%s \n", fn, err)
		return
	}

	fullPath := filepath.Join(savePath, fileName)
	if err = os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		fmt.Printf("%s creating directories err:%s \n", fn, err)
		return
	}

	f, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("%s creating file:%s \n", fn, err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	bar := pb.Start64(int64(fileSize))
	bar.Set(pb.Bytes, true)

	defer bar.Finish()

	proxyReader := bar.NewProxyReader(reader)
	if _, err = io.Copy(f, proxyReader); err != nil {
		fmt.Printf("%s saving file:%s", fn, err)
		return
	}
	fmt.Printf("File received:%s", fullPath)
}
