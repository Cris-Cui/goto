package main

import (
	"flag"
	"fmt"
	"net/http"
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
)

var store *URLStore = nil

func main() {
	flag.Parse()
	store = NewURLStore(*dataFile)
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	http.ListenAndServe(*listenAddr, nil)
}

// Redirect 重定向Http Request处理函数
func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	fmt.Println("重定向的 key 为: " + key)
	url := store.Get(key)
	fmt.Println(key + "对应的value为: " + url)
	if url == "" {
		http.NotFound(w, r)
		return
	}
	// NOTION: 应该保证URL为绝对URL, 以http:// 或https:// 开头
	http.Redirect(w, r, url, http.StatusFound)
}

// Add 映射短URL处理函数
func Add(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		w.Header().Set("Content-Type", "text/html") // 设置Content-Type为HTML
		fmt.Fprint(w, AddForm)
		return
	}
	key := store.Put(url)
	fmt.Fprintf(w, "http://%s/%s", *hostname, key)
}
