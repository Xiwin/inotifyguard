package inotify

import (
	"sync"
)

// global variable define
var (
	// Global list for Instance
	// uint32: Timestamp for creating the instance
	instanceList []*InotifyInstance
)

// custom type alias
type WDMap map[string]int

// struct define
type InotifyInstance struct {
	InotifyFd int
	Events    chan InotifyEvent // 事件通告缓冲管道
	Error     chan error        // 错误通告管道
	Done      chan bool         // 监控结束通告管道
	Mutex     sync.Mutex        // 单线程锁
	PathList  WDMap             // Watcher Descriptor 与 Path 的映射表
}

type InotifyEvent struct {
	Path   string
	Wd     int
	Mask   uint32
	Cookie uint32
	Name   string
}
