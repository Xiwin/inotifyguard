package inotify

import (
	"log"
	"syscall"
)

// Create a new Inotify instance
func New(bufferSize int) (*InotifyInstance, error) {
	newInotify := &InotifyInstance{
		PathList: make(WDMap),
		Events:   make(chan InotifyEvent, bufferSize),
		Error:    make(chan error, 1),
		Done:     make(chan bool, 1),
	}

	var err error
	newInotify.InotifyFd, err = syscall.InotifyInit1(syscall.IN_NONBLOCK | syscall.O_CLOEXEC)
	if err != nil {
		return nil, err
	}

	go newInotify.readInotifyInstance()

	// Add inotify instance into global instance list
	instanceList = append(instanceList, newInotify)

	return newInotify, nil
}

// Add a specified watcher to inotify descriptor
func (instance *InotifyInstance) AddWatch(pathname string) error {
	wd, err := syscall.InotifyAddWatch(instance.InotifyFd, pathname, syscall.IN_ALL_EVENTS)
	if err != nil {
		return err
	}

	instance.Mutex.Lock()
	defer instance.Mutex.Unlock()
	instance.PathList[pathname] = wd

	return nil
}

// Remove a specified watcher from an inotify instance
func (instance *InotifyInstance) RmWatch(pathname string) error {
	instance.Mutex.Lock()
	defer instance.Mutex.Unlock()

	watch, ok := instance.PathList[pathname]
	if !ok {
		return nil // Path not being watched
	}

	_, err := syscall.InotifyRmWatch(instance.InotifyFd, uint32(watch))
	if err != nil {
		return err
	}

	delete(instance.PathList, pathname)

	return nil
}

// Remove any watchers that still exist, and close inotify's file descriptor
// Close removes any watchers that still exist and closes inotify's file descriptor.
func (instance *InotifyInstance) Close() {
	instance.Mutex.Lock()
	defer instance.Mutex.Unlock()

	for pathname := range instance.PathList {
		err := instance.RmWatch(pathname)
		if err != nil {
			log.Println("Error removing watch:", err)
		}
	}

	err := syscall.Close(instance.InotifyFd)
	if err != nil {
		log.Fatalln("Error closing inotify instance:", err)
	}

	instance.Done <- true
	close(instance.Events)
	close(instance.Error)
	close(instance.Done)
}
