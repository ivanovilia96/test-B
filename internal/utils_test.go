package internal

import (
	"sync"
	"testing"
	"time"
)

// технически правильно написанный тест (+-), но его надо еще визуально в порядок привести
func TestChecker(t *testing.T) {
	t.Run("correct work, expected an valid result in sleepTime time", func(t *testing.T) {
		var messages *[]Message = &[]Message{}
		v1 := Message{
			Message: "1",
		}
		sleepTime := 1 // sec
		timeout := 2   //sec
		curTime := time.Now()
		go func() {
			time.Sleep(time.Duration(sleepTime) * time.Second)
			*messages = append(*messages, v1)
		}()

		message := messageWaiter(timeout, messages, new(sync.RWMutex))

		if message == nil { // sleepTime < timeout , значит мы ожидаем что вернется v1 переменная
			t.Error("un expected  error, message == nil but we expect correct result")
			return
		}

		if message.Message != v1.Message {
			t.Error("un expected  error, returned message != expected message", message, v1)
			return
		}

		if time.Since(curTime) > (time.Duration(timeout) * time.Second) {
			t.Error("un expected  error , sleepTime < timeout мы ожидаем что результат никак не может вернутся позже таймаута < timeout")
			return
		}
	})
}
