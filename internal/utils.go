package internal

import (
	"context"
	"sync"
	"time"
)

// messageWaiter - проверяет в течении timeout с интервалом checkerIntervalSec появились ли сообщения , если появились, то возвращает 1 из них ( при возвращении удаляет это сообщение из очереди ), иначе по timeout nil
func messageWaiter(timeout int, messages *[]Message, mutex *sync.RWMutex) *Message {
	ticker := time.NewTicker(checkerIntervalSec)
	defer ticker.Stop()

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	for {
		select {
		case <-ticker.C: // проверяем есть ли сообщения в очереди c определенным интервалом
			if len(*messages) != 0 {
				mutex.Lock()
				message := &(*messages)[0]
				*messages = (*messages)[1:] // удаляем первое сообщение из очереди
				mutex.Unlock()
				return message
			}
		case <-ctx.Done():
			return nil
		}
	}
}
