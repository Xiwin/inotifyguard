## Overview

一个基于 Linux 的 inotify 机制的文件监控库，可以监控文件状态


## Feature

程序添加被监控文件所有的事件。需要用到什么事件的通告，在接受时只处理该通告即可。其它的事件通告可以忽略不计。在watch过程中，捕获的所有错误都会输出到 Error 管道中。事件管道是一个可缓冲管道，事件的通告使用的是单线通告，一次只通告一个事件的内容。

可以使用 `is` 系列函数，用以判断是否是指定事件，`is` 系列函数，会返回监视目标文件。

基本案例，可在 `test` 目录下找到

## 相关文档

inotify 机制说明文档:[https://www.man7.org/linux/man-pages/man7/inotify.7.html](https://www.man7.org/linux/man-pages/man7/inotify.7.html)

golang **syscall** 标准库文档: [https://pkg.go.dev/syscall](https://pkg.go.dev/syscall)