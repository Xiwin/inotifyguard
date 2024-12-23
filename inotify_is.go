package inotify

import "syscall"

// IsModify，当文件的内容被修改时，会出现该事件通告
func (event InotifyEvent) IsModify() (filename string, ok bool) {
	if event.Mask == syscall.IN_MODIFY {
		return event.Path, true
	}
	return "", false
}

// IsWrite, 当文件被使用 read 系统调用以读写的当时打开并关闭后，会出现该事件通告
func (event InotifyEvent) IsWrite() (filename string, ok bool) {
	if event.Mask == syscall.IN_CLOSE_WRITE {
		return event.Path, true
	}

	return "", false
}

// IsRead，当文件被使用 read系统调用以只读的方式打开并关闭后，会出现该事件通告
func (event InotifyEvent) IsRead() (filename string, ok bool) {
	if event.Mask == syscall.IN_CLOSE_NOWRITE {
		return event.Path, true
	}

	return "", false
}

// IsCreate，当监视的目录中有新增文件或目录时，会出现该事件通告，如果是该函数会返回新增内容和被监控目标目录的路径
func (event InotifyEvent) IsCreate() (desFilename string, newFilename string, ok bool) {
	if event.Mask == syscall.IN_CREATE {
		return event.Path, event.Name, true
	}

	return "", "", false
}
