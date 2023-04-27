package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) getMessageFromQueue(w http.ResponseWriter, r *http.Request) {
	matches := queue.FindStringSubmatch(r.URL.Path) // от сюда пытаемся получить имя очереди
	if len(matches) < 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("dont find queueName %s", r.URL.Path)))
		return
	}

	queueName := matches[1]

	h.Queue.RLock()
	messages, has := h.Queue.Store[queueName] // проверяем есть ли такая очередь
	h.Queue.RUnlock()
	if !has {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("queue %s is`nt exists", queueName)))
		return
	}

	// проверяем еть ли сообщения в очереди
	if len(*messages) == 0 {
		timeouts := r.URL.Query()["timeout"]
		timeout := 0
		if len(timeouts) != 0 { // значит хоть 1 таймаут не дефолтный задан
			t, err := strconv.Atoi(timeouts[0]) // даже если их много передано, то берем 1 на данный момент
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("did not find parameter timeout"))
				return
			}

			timeout = t
		} else {
			timeout = h.Queue.DefaultTimeout
		}

		// ждем сообщения в очередь timeout секунд
		result := messageWaiter(timeout, messages, &h.Queue.RWMutex)
		if result == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("queue is empty, timeout was " + strconv.Itoa(timeout)))
			return
		}

		m, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("unmarshal error " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(m)
		return
	}

	h.Queue.Lock()
	message := (*messages)[0]
	*messages = (*messages)[1:] // удаляем первое сообщение из очереди
	h.Queue.Unlock()

	m, err := json.Marshal(message) // в принципе тут ошибка не возможна, её можно игнорировать
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("marshal error" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(m)
	return
}
