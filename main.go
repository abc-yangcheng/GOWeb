package main

import (
	"net/http"
	"time"
)

/**
1. 客户端向服务器发送请求
2. 多路复用器接收到请求，并将其重定向到正确的处理器
3. 处理器对请求进行处理
4. 在需要访问数据库的情况下，处理器会使用一个或多个数据结构，这些数据结构都是蜂聚数据库中的数据建模而来的
5. 当处理器调用与数据结构有关的函数或者方法时，这些数据结构背后的模型会与数据库进行连接，并执行相应的操作
6. 当请求处理完毕时，处理器会调用模板引擎，有时候还会向模板引擎传递一些通过模型获取到的数据
7. 模板引擎会对模板文件进行语法分析并创建相应的模板，而这些模板又会与处理器传递的数据合并生成最终的HTML
8. 生成的HTML会作为响应的一部分传至客户端
*/

func main() {
	p("ChitChat", version(), "started at", config.Address)

	// handle static assets
	// 首先创建一个多路复用器
	mux := http.NewServeMux()
	// 除了负责请求重定向到相应的处理器外，多路复用器还需要为静态文件提供服务
	// FileServer函数创建了一个能够为指定目录中的静态文件服务的处理器
	files := http.FileServer(http.Dir(config.Static))
	// 当服务器接收到一个以/static/ 开头的URL请求时，StripPrefix会移除URL中的/static/ 字符串
	// 然后在public目录中查找被请求的文件
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// all route patterns matched here
	// route handler functions defined in other files
	//

	// index
	mux.HandleFunc("/", index)
	// error
	mux.HandleFunc("/err", err)

	// defined in route_auth.go
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)

	// defined in route_thread.go
	mux.HandleFunc("/thread/new", newThread)
	mux.HandleFunc("/thread/create", createThread)
	mux.HandleFunc("/thread/post", postThread)
	mux.HandleFunc("/thread/read", readThread)

	// starting up the server
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
