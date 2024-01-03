### Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func asChan(vs ...int) <-chan int {
    c := make(chan int)

    go func() {
        for _, v := range vs {
            c <- v
            time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
        }

        close(c)
    }()
    return c
}

func merge(a, b <-chan int) <-chan int {
    c := make(chan int)
    go func() {
        for {
            select {
            case v := <-a:
                c <- v
            case v := <-b:
                c <- v
            }
        }
    }()
    return c
}

func main() {

    a := asChan(1, 3, 5, 7)
    b := asChan(2, 4 ,6, 8)
    c := merge(a, b)
    for v := range c {
        fmt.Println(v)
    }
}
```

### Ответ:
Программа выведет все значения из каналов a и b в случайном порядке (порядок значений переданных в канал сохранится), а дальше будет выводить 0. Случайный порядок обоснован функцией `select`, так как она при возможности чтения с нескольких канало выбирает `case` произвольного порядка. Нули будут выводится из-за того, что функция `merge` будет читать из закрытых каналов, в связи с отсутствием вызова функции закрытия канала, а, так как, тип каналов `int`, прочитанные значение будут равны базовому, т.е. нули.