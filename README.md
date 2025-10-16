# or-channel

Пакет для объединения нескольких done-каналов в один.

## Установка

```bash
go get github.com/sunr3d/or-channel
```

## Использование

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/sunr3d/or-channel"
)

func main() {
    // Создаем каналы
    ch1 := make(chan interface{})
    ch2 := make(chan interface{})

    // Запускаем горутины
    go func() {
        time.Sleep(1 * time.Second)
        close(ch1)
    }()

    go func() {
        time.Sleep(2 * time.Second)
        close(ch2)
    }()

    // Ждем закрытия любого канала
    done := orchan.Or(ch1, ch2)
    <-done

    fmt.Println("Один из каналов закрыт!")
}
```

## API

```go
func Or(channels ...<-chan interface{}) <-chan interface{}
```

Принимает done-каналы и возвращает общий done-канал, который закрывается при закрытии любого из входных каналов.

## Когда использовать

- Динамическое количество каналов
- Нужно ждать завершения любого из нескольких процессов
- Упрощение сложных `select` с множественными `case`

## Тестирование

```bash
go test -v
```
