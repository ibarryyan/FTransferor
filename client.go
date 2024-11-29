package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

var (
	server, file, action string
)

func ClientCmd() *cobra.Command {
	command := &cobra.Command{
		Use: "cli",
		Run: func(cmd *cobra.Command, args []string) {
			startClient(server, file)
		},
	}

	command.Flags().StringVarP(&server, "server", "", "", "path to serve")
	command.Flags().StringVarP(&file, "file", "", "", "path to serve")
	command.Flags().StringVarP(&action, "action", "", "", "path to serve")
	return command
}

// Client: 发送文件
func startClient(serverAddr string, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	// 获取文件信息
	fileInfo, err := f.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	// 发送文件元信息（文件名|文件大小）
	meta := fmt.Sprintf("%s|%d\n", fileInfo.Name(), fileInfo.Size())
	if _, err = conn.Write([]byte(meta)); err != nil {
		fmt.Println("Error sending metadata:", err)
		return
	}

	// 创建进度条
	bar := pb.Start64(fileInfo.Size())
	bar.Set(pb.Bytes, true)
	defer bar.Finish()

	// 发送文件数据
	proxyWriter := bar.NewProxyWriter(conn)
	if _, err = io.Copy(proxyWriter, f); err != nil {
		fmt.Println("Error sending file:", err)
		return
	}

	fmt.Println("File sent successfully!")
}
