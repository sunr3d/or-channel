package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sunr3d/or-channel"
)

func main() {
	n := rand.Intn(5) + 5 // 5-10 воркеров
	fmt.Printf("Запускаем %d воркеров...\n", n)

	var workers []<-chan interface{}
	for i := 1; i <= n; i++ {
		workers = append(workers, worker(i))
	}

	done := orchan.Or(workers...)
	<-done

	fmt.Println("Один из воркеров завершился!")
}

func worker(id int) <-chan interface{} {
	done := make(chan interface{})

	go func() {
		fmt.Printf("Воркер %d начал работу\n", id)
		time.Sleep(time.Duration(rand.Intn(5)+2) * time.Second) // работает 2-7 секунды
		fmt.Printf("Воркер %d завершился\n", id)
		close(done)
	}()

	return done
}
