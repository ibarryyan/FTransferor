package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

type Action string

const (
	ActionList Action = "list"
	ActionGet  Action = "get"
	Scheme            = "http://"
)

var (
	server, file, action, passwd string
)

func ClientCmd() *cobra.Command {
	command := &cobra.Command{
		Use: "cli",
		Run: func(cmd *cobra.Command, args []string) {
			if server == "" {
				fmt.Println("err")
				return
			}
			if action != "" {
				httpClient(Action(action), file)
			} else {
				startClient(server, file)
			}
		},
	}

	command.Flags().StringVarP(&passwd, "passwd", "", "", "Server password")
	command.Flags().StringVarP(&server, "server", "", "", "Server address and port")
	command.Flags().StringVarP(&file, "file", "", "", "Upload file name")
	command.Flags().StringVarP(&action, "action", "", "", "Client action")
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

func httpClient(action Action, f string) {
	httpServerAddr := fmt.Sprintf("%s%s", Scheme, server)

	client := http.Client{}
	switch action {
	case ActionGet:
		url := fmt.Sprintf("%s%s%s?%s=%s", httpServerAddr, PathDownload, f, QueryParamKey, passwd)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Error downloading file:", err)
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		// 创建文件以保存下载的内容
		file, err := os.Create(f)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer func() {
			_ = file.Close()
		}()

		// 将响应的内容写入文件
		if _, err = io.Copy(file, resp.Body); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("File downloaded successfully!")
	case ActionList:
		url := fmt.Sprintf("%s%s?%s=%s", httpServerAddr, PathList, QueryParamKey, passwd)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Error downloading file:", err)
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		// 读取响应的body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// 将body转换为map
		var data map[string]interface{}
		if err = json.Unmarshal(body, &data); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
		fmt.Println("File List:")
		if val, ok := data["data"]; ok && val != nil {
			list := val.([]interface{})
			for i, item := range list {
				fmt.Printf("Num:%d  File: %s\n", i+1, item)
			}
		}
	default:

	}
}
