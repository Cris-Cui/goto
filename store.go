package main

import (
	"encoding/gob"
	"log"
	"os"
	"sync"
)

// URLStore URLStore类型 是一个结构体
type URLStore struct {
	urls map[string]string // 短网址到长网址的映射, key 是 短网址, value 是 长网址
	mu   sync.RWMutex      // 读写锁
	file *os.File          // 文件指针, kv键值对的持久化文件
}

// record 持久化到文件中的kv记录
type record struct {
	Key, URL string
}

// NewURLStore URLStore工厂函数
func NewURLStore(filename string) *URLStore {
	s := &URLStore{urls: make(map[string]string)}
	// 追加模式可写打开文件
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal("URLStore: ", err)
	}
	s.file = f
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
			return key
		}
	}
	// shouldn't get here
	return ""
}

// save 将给定的 key 和 url 作为一个 gob 编码的 record 写入到磁盘
func (s *URLStore) save(key, url string) error {
	e := gob.NewEncoder(s.file) // e为编码器
	return e.Encode(record{key, url})
}

/*
 * map    			章节8
 * Mutex		    章节9
 * struct和方法	    章节10
 */
