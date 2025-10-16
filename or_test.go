package orchan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createChannel(duration time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		time.Sleep(duration)
		close(c)
	}()
	return c
}

// Тест 1: Пустой список каналов
func TestOr_NoChannels(t *testing.T) {
	done := Or()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Канал должен быть закрыт сразу")
	}
}

// Тест 2: Один канал
func TestOr_OneChannel(t *testing.T) {
	ch := createChannel(50 * time.Millisecond)
	done := Or(ch)

	start := time.Now()
	<-done
	elapsed := time.Since(start)

	assert.True(t, elapsed >= 40*time.Millisecond, "Канал закрылся слишком рано")
	assert.True(t, elapsed <= 100*time.Millisecond, "Канал закрылся слишком поздно")
}

// Тест 3: Два канала - один быстрый, один медленный
func TestOr_TwoChannels(t *testing.T) {
	fast := createChannel(50 * time.Millisecond)
	slow := createChannel(200 * time.Millisecond)

	done := Or(fast, slow)

	start := time.Now()
	<-done
	elapsed := time.Since(start)

	assert.True(t, elapsed >= 40*time.Millisecond, "Канал закрылся слишком рано")
	assert.True(t, elapsed <= 100*time.Millisecond, "Канал закрылся слишком поздно")
}

// Тест 4: Уже закрытый канал
func TestOr_ClosedChannel(t *testing.T) {
	closed := make(chan interface{})
	close(closed)

	done := Or(closed)

	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		assert.Fail(t, "Канал должен быть закрыт сразу")
	}
}

// Тест 5: Три канала с разными временами
func TestOr_ThreeChannels(t *testing.T) {
	ch1 := createChannel(100 * time.Millisecond)
	ch2 := createChannel(50 * time.Millisecond)
	ch3 := createChannel(150 * time.Millisecond)

	done := Or(ch1, ch2, ch3)

	start := time.Now()
	<-done
	elapsed := time.Since(start)

	assert.True(t, elapsed >= 40*time.Millisecond, "Канал закрылся слишком рано")
	assert.True(t, elapsed <= 80*time.Millisecond, "Канал закрылся слишком поздно")
}

// Тест 6: Смесь обычных и закрытых каналов
func TestOr_MixedChannels(t *testing.T) {
	normal1 := createChannel(100 * time.Millisecond)
	normal2 := createChannel(300 * time.Millisecond)
	closed := make(chan interface{})
	close(closed)

	done := Or(normal1, normal2, closed)

	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		assert.Fail(t, "Канал должен быть закрыт сразу")
	}
}

// Тест 7: Проверка, что функция не падает при закрытии каналов одновременно
func TestOr_RaceCondition(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})
	ch4 := make(chan interface{})

	done := Or(ch1, ch2, ch3, ch4)

	close(ch1)
	close(ch2)
	close(ch3)
	close(ch4)

	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		assert.Fail(t, "Канал должен быть закрыт сразу")
	}
}

// Тест 8: Проверка времени с assert
func TestOr_Timing(t *testing.T) {
	ch := createChannel(100 * time.Millisecond)
	done := Or(ch)

	start := time.Now()
	<-done
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, elapsed, 90*time.Millisecond, "Слишком быстро")
	assert.LessOrEqual(t, elapsed, 150*time.Millisecond, "Слишком медленно")
}

// Тест 9: Проверка, что канал действительно закрыт
func TestOr_ChannelIsClosed(t *testing.T) {
	ch := createChannel(50 * time.Millisecond)
	done := Or(ch)

	<-done

	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		assert.Fail(t, "Канал должен быть закрыт")
	}
}
