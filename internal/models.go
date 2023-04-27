package internal

import (
	"sync"
)

type Message struct {
	Message string `json:"message"`
}

type Queue struct {
	MaxQueuesCount             int
	MaxMessagesInOneQueueCount int
	DefaultTimeout             int
	Store                      map[string]*[]Message // key - queue name, value - FIFO queue  // указатель на массив сделан что бы мы могли при таймауте видеть изменения переменных в реальном времени без труда
	sync.RWMutex
}
