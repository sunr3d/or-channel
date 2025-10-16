// Package orchan предоставляет функцию для объединения нескольких done-каналов
// в один done-канал, который закрывается при закрытии любого канала из списка.
package orchan

import "sync"

// Or принимает несколько done-каналов и возвращает done-канал,
// который закрывается, когда закрывается любой из входных каналов.
//
// Функция использует примитив синхронизации sync.Once для безопасного закрытия результирующего канала
// только один раз, даже если несколько каналов из списка закрываются одновременно.
//
// Если список каналов пуст, возвращает закрытый канал (предотвращает блокировку).
//
// Если в списке только один канал, возвращает его (не создаем лишних горутин).
//
// Если один из входных каналов уже закрыт, то done-канал закрывается немедленно.
//
// Пример использования:
//
//	done := orchan.Or(ch1, ch2, ch3)
//	<-done // блокируется до закрытия любого из ch1, ch2, ch3
func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		closed := make(chan interface{})
		close(closed)
		return closed
	case 1:
		return channels[0]
	default:
		res := make(chan interface{})
		var once sync.Once

		for _, channel := range channels {
			go func(c <-chan interface{}) {
				<-c
				once.Do(func() {
					close(res)
				})
			}(channel)
		}

		return res
	}
}
