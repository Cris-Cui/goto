package main

import "sync"

// URLStore URLStore类型 是一个结构体
type URLStore struct {
	urls map[string]string // 短网址到长网址的映射, key 是 短网址, value 是 长网址
	mu   sync.RWMutex      // 读写锁
}

// NewURLStore URLStore工厂函数
func NewURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
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

/*
 * map    			章节8
 * Mutex		    章节9
 * struct和方法	    章节10
 */
