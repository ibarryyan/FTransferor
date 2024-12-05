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
				runServer(port, path)
			}()

			<-quit
		},
	}

	command.Flags().StringVarP(&path, "path", "", "", "path to serve")
	command.Flags().StringVarP(&secret, "secret", "", "", "path to serve")
	command.Flags().IntVarP(&port, "port", "", 0, "path to serve")
	command.Flags().IntVarP(&webport, "webport", "", 0, "path to serve")
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

func runServer(port int, savePath string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer func() {
		_ = listener.Close()
	}()

	fmt.Println("Server is listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn, savePath)
	}
}

func handleConnection(conn net.Conn, savePath string) {
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)

	// 读取文件元信息：文件名和文件大小
	meta, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file metadata:", err)
		return
	}
	meta = strings.TrimSpace(meta) // 清除换行符
	parts := strings.Split(meta, "|")
	if len(parts) != 2 {
		fmt.Println("Invalid metadata received")
		return
	}
	fileName := parts[0]
	fileSize := 0
	_, err = fmt.Sscanf(parts[1], "%d", &fileSize)
	if err != nil {
		fmt.Println("Error parsing file size:", err)
		return
	}

	// 确保保存路径存在
	fullPath := filepath.Join(savePath, fileName)

	if err = os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}

	// 创建文件
	f, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	// 创建进度条
	bar := pb.Start64(int64(fileSize))
	bar.Set(pb.Bytes, true)

	defer bar.Finish()

	// 读取数据并写入文件
	proxyReader := bar.NewProxyReader(reader)
	if _, err = io.Copy(f, proxyReader); err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
	fmt.Println("File received:", fullPath)
}
