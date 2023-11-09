package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/rpc"
)

const AddForm = `
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
`

var (
	listenAddr = flag.String("http", ":8080", "http listen address")
	dataFile   = flag.String("file", "store.json", "data store file name")
	hostname   = flag.String("host", "localhost:8080", "http host name")
	rpcEnabled = flag.Bool("rpc", false, "enable rpc server")
	masterAddr = flag.String("master", "", "RPC master address")
)

var store Store

func main() {
	flag.Parse()
	if *masterAddr != "" { // 主服务器地址不为空, 是一个从服务器
		store = NewProxyStore(*masterAddr)
	} else {
		store = NewURLStore(*dataFile)
	}
	if *rpcEnabled { // 启动了RPC服务
		rpc.RegisterName("Store", store)
		rpc.HandleHTTP()
	}
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	http.ListenAndServe(*listenAddr, nil)
}

// Redirect 重定向Http Request处理函数
func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	fmt.Println("重定向的 key 为: " + key)
	var url string
	if err := store.Get(&key, &url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500 状态码响应
		return
	}
	fmt.Println(key + "对应的value为: " + url)
	// NOTION: 应该保证URL为绝对URL, 以http:// 或https:// 开头
	http.Redirect(w, r, url, http.StatusFound) // 302 状态码响应
}

// Add 映射短URL处理函数
func Add(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		w.Header().Set("Content-Type", "text/html") // 设置Content-Type为HTML
		fmt.Fprint(w, AddForm)
		return
	}
	var key string
	if err := store.Put(&url, &key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500 状态码响应
		return
	}
	fmt.Fprintf(w, "http://%s/%s", *hostname, key)
}
