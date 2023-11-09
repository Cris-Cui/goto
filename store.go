package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

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
	s := &URLStore{urls: make(map[string]string), save: make(chan record, saveQueueLength)}

	// 从磁盘中加载数据到map中
	if err := s.load(filename); err != nil {
		log.Println("Error loading URLStore: ", err)
	}
	fmt.Println(s)
	go s.saveLoop(filename)
	return s
}

// Get 重定向读类型请求的URLStore指针变量的方法
func (s *URLStore) Get(key string) string {
	s.mu.RLock()         // 上读锁
	defer s.mu.RUnlock() // 函数结束时释放读锁
	return s.urls[key]   // 返回value string类型
}

// Set 处理写请求的URLStore指针变量的**方法**
func (s *URLStore) Set(key, url string) bool {
	s.mu.Lock()                             // 上写锁
	defer s.mu.Unlock()                     // 函数结束后释放写锁
	if _, present := s.urls[key]; present { // 逗号ok模式,
		return false // key存在, 返回false
	}
	s.urls[key] = url
	return true
}

// Count 计算map中键值对的数量的URLStore指针变量的**方法**
func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

// Put 将长网址映射到短网址并set到map的URLStore指针变量的方法
func (s *URLStore) Put(url string) string {
	for { // for死循环一直尝试keygen
		key := genKey(s.Count()) // generate the short URL
		if ok := s.Set(key, url); ok {
			// 先做持久化(将key-value放到channel通道中), 再返回key
			s.save <- record{key, url}
			return key
		}
	}
	panic("shouldn't get here")
}

// saveLoop 将给定的 key 和 url 作为一个 gob 编码的 record 写入到磁盘
func (s *URLStore) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal("URLStore: ", err)
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	for {
		r := <-s.save
		if err := e.Encode(r); err != nil {
			log.Println("URLStore: ", err)
		}
	}
}

// load 在程序启动后, 需要将磁盘上的数据读到URLStore中
func (s *URLStore) load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening URLStore: ", err)
		return err
	}
	defer f.Close()
	d := gob.NewDecoder(f) // 解码器
	for err == nil {
		var r record
		if err = d.Decode(&r); err == nil {
			s.Set(r.Key, r.URL)
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
