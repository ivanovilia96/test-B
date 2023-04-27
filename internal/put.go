package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) putMessageToQueue(w http.ResponseWriter, r *http.Request) {
	matches := queue.FindStringSubmatch(r.URL.Path) // от сюда пытаемся получить имя очереди
	if len(matches) < 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("dont find queueName %s", r.URL.Path)))
		return
	}

	queueName := matches[1]

	h.Queue.RLock()
	// проверяем есть ли такая очередь и получаем её
	queue, has := h.Queue.Store[queueName]
	h.Queue.RUnlock()
	if !has && len(h.Queue.Store) < h.Queue.MaxQueuesCount { // если нету еще такой очереди и их кол-во меньше чем MaxQueuesCount, то создаем
		messages := make([]Message, 0)
		h.Queue.Lock()
		h.Queue.Store[queueName] = &messages
		queue = h.Queue.Store[queueName]
		h.Queue.Unlock()
	} else if !has && len(h.Queue.Store) >= h.Queue.MaxQueuesCount { // значит мы не можем больше создать очередей и пора возвращать ошибку
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("max queues count is exideed, max queues count: %d", h.Queue.MaxQueuesCount)))
		return
	}

	// блок проверки на кол-во сообщений в очереди относительно аргумента командной строки
	if len(*queue) >= h.Queue.MaxMessagesInOneQueueCount { // получается не можем положить сообщение
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("max messages count is exideed, max messages count: %d", h.Queue.MaxMessagesInOneQueueCount)))
		return
	}

	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error in req body, decode error  %s", err.Error())))
		return
	}

	h.Queue.Lock()
	*queue = append(*queue, message)
	h.Queue.Unlock()

	w.WriteHeader(http.StatusOK)
}
