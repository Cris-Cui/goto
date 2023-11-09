package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"sync"
)

type Store interface {
	Put(url, key *string) error
	Get(key, url *string) error
}

// ProxyStore RPC代理工厂
type ProxyStore struct {
	urls   *URLStore
	client *rpc.Client
}

// URLStore URLStore类型 是一个结构体
type URLStore struct {
	urls map[string]string // 短网址到长网址的映射, key 是 短网址, value 是 长网址
	mu   sync.RWMutex      // 读写锁
	save chan record       // record 类型channel
}

// record 持久化到文件中的kv记录
type record struct {
	Key, URL string
}

// save channel 缓冲区大小
const saveQueueLength = 1000

// NewURLStore URLStore工厂函数
func NewURLStore(filename string) *URLStore {
	s := &URLStore{urls: make(map[string]string)}

	if filename != "" {
		// 从磁盘中加载数据到map中
		s.save = make(chan record, saveQueueLength)
		if err := s.load(filename); err != nil {
			log.Println("Error loading URLStore: ", err)
		}
		fmt.Println(s)
		go s.saveLoop(filename)
	}
	return s
}

// NewProxyStore RPC客户端代理工厂构造函数
func NewProxyStore(addr string) *ProxyStore {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Println("Error constructing ProxyStore:", err)
	}
	return &ProxyStore{urls: NewURLStore(""), client: client}
}

// Get ProxyStore Get方法 传递Get请求给RPC服务端
func (s *ProxyStore) Get(key, url *string) error {
	// 先查看从服务器缓存中是否存在数据记录, 如果没有RPC调用去查询主服务器
	if err := s.urls.Get(key, url); err == nil { // 本地map中找到url
		return nil
	}
	if err := s.client.Call("Store.Get", key, url); err != nil {
		return err
	}
	s.urls.Set(key, url) // 将RPC调用返回的数据写到本地map中
	return nil
}

// Put ProxyStore Put方法 传递Put请求给RPC服务端
func (s *ProxyStore) Put(url, key *string) error {
	if err := s.client.Call("Store.Put", url, key); err != nil {
		return err
	}
	s.urls.Set(key, url)
	return nil
}

// Get URLStore 重定向URL处理器
func (s *URLStore) Get(key, url *string) error {
	s.mu.RLock()         // 上读锁
	defer s.mu.RUnlock() // 函数结束时释放读锁
	if u, ok := s.urls[*key]; ok {
		*url = u
		return nil
	}
	return errors.New("key not found")
}

// Set URLStore 处理写请求的URLStore指针变量的**方法**
func (s *URLStore) Set(key, url *string) error {
	s.mu.Lock()                              // 上写锁
	defer s.mu.Unlock()                      // 函数结束后释放写锁
	if _, present := s.urls[*key]; present { // 逗号ok模式,
		return errors.New("key already exists")
	}
	s.urls[*key] = *url
	return nil
}

// Count URLStore 计算map中键值对的数量的URLStore指针变量的**方法**
func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

// Put URLStore 长URL转短URL处理器
func (s *URLStore) Put(url, key *string) error {
	for { // for死循环一直尝试keygen
		*key = genKey(s.Count()) // generate the short URL
		if err := s.Set(key, url); err == nil {
			break
		}
	}
	if s.save != nil {
		s.save <- record{*key, *url}
	}
	return nil
}

// saveLoop URLStore 将给定的 key 和 url 作为一个 gob 编码的 record 写入到磁盘
func (s *URLStore) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal("URLStore: ", err)
	}
	defer f.Close()
	e := json.NewEncoder(f)
	for {
		r := <-s.save
		if err := e.Encode(r); err != nil {
			log.Println("URLStore: ", err)
		}
	}
}

// load URLStore 在程序启动后, 需要将磁盘上的数据读到URLStore中
func (s *URLStore) load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening URLStore: ", err)
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f) // 解码器
	for err == nil {
		var r record
		if err = d.Decode(&r); err == nil {
			s.Set(&r.Key, &r.URL)
		}
	}
	if err == io.EOF {
		return nil
	}
	// error occurred
	log.Println("Error decoding URLStore: ", err)
	return err
}

/*
 * map    			章节8
 * Mutex		    章节9
 * struct和方法	    章节10
 */
