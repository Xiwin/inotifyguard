package inotify

import (
	"syscall"
	"unsafe"
)

// readInotifyInstance reads and processes events from the inotify instance.
func (instance *InotifyInstance) readInotifyInstance() {
	var buffer [syscall.SizeofInotifyEvent * 1024]byte
	for {
		select {
		case exit := <-instance.Done:
			if exit {
				close(instance.Events)
				close(instance.Error)
				return
			}
		default:
			n, err := syscall.Read(instance.InotifyFd, buffer[:])
			if err != nil {
				if err == syscall.EAGAIN {
					continue
				}
				instance.Error <- err
				return
			}

			// 一次读取的内容可能会有多个事件，采用循环的方式解析并推送事件到事件通知管道中
			var offset uint32

			// 读取的长度 - InotifyEvent的默认长度,能够获取当前当前是否读取了一个完整的事件通知，只要有一个就处理
			for offset <= uint32(n-syscall.SizeofInotifyEvent) {

				// 将当前 offset 所指定的 byte 类型指针，转为 inotifyEvent 指针类型
				event := (*syscall.InotifyEvent)(unsafe.Pointer(&buffer[offset]))

				var name string
				nameLength := event.Len
				if nameLength > 1 {

					nameBytes := buffer[offset+syscall.SizeofInotifyEvent : offset+syscall.SizeofInotifyEvent+event.Len] // 获取 name 字段的长度
					name = string(nameBytes[:len(nameBytes)-1])                                                          // name 字段的长度，因为 len 包含了 null 字符，所以减去 1，string 需要指定拷贝的字符数量，它会自动拷贝 null 字符，所以无需指定
				}

				path := instance.GetPath(int(event.Wd))

				// 将 buffer 中的事件转换为事件管道的类型
				inotifyEvent := InotifyEvent{
					Wd:     int(event.Wd),
					Mask:   event.Mask,
					Cookie: event.Cookie,
					Name:   name,
					Path:   path,
				}

				instance.Events <- inotifyEvent

				offset += syscall.SizeofInotifyEvent + event.Len
			}
		}
	}
}

func (instance *InotifyInstance) GetPath(wd int) string {
	for k, v := range instance.PathList {
		if v == wd {
			return k
		}
	}
	return ""
}
