package test

import (
	"fmt"
	"inotify"
	"log"
	"os"
	"syscall"
	"testing"
)

func TestAddWatch(t *testing.T) {

	work_dir, _ := os.Getwd()

	file_name := "watched_file"

	file_path := fmt.Sprintf("%s/%s", work_dir, file_name)

	instance, err := inotify.New(1024)
	if err != nil {
		log.Fatalln(err)
	}

	defer instance.Close()

	file, _ := os.Create(file_path)
	file.Close()

	fmt.Println("开始监控")
	err = instance.AddWatch(file_path)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			select {
			case event := <-instance.Events:
				path, ok := event.IsWrite()
				if ok {
					fmt.Println(path)
				}
			case err := <-instance.Error:
				log.Fatalln(err)
				syscall.Exit(1)
			}
		}
	}()

	<-make(chan struct{})

}
