package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 读取index.html文件
		content, err := os.ReadFile("plugin_usage/index.html")
		if err != nil {
			http.Error(w, "文件未找到", http.StatusNotFound)
			return
		}
		// 设置响应的内容类型为HTML
		w.Header().Set("Content-Type", "text/html")
		// 写入HTML内容到响应
		w.Write(content)
	})

	// 设置服务器地址和端口
	addr := ":8080"
	log.Printf("服务器正在运行，访问 http://localhost%s\n", addr)

	// 启动HTTP服务
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
