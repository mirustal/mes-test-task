package message

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)



// SetRequest представляет запрос для установки нового сообщения
// swagger:parameters SetMessage
type SetRequest struct {
	// Текст сообщения
    // in: body
    // required: true
    // example: "Hello, World!"
	Text string `json:"text"`
}

// SetResponse представляет ответ на запрос установки сообщения
// swagger:response SetResponse
type SetResponse struct {
	// ID нового сообщения
    // example: "111"
	ID      string `json:"id,omitempty"`
}

type messageSetter interface {
	Insert(ctx context.Context, text string) (string, error)
}


type Producer interface {
	ProduceMessage(id, text string) error
}

// NewSetter создает новый обработчик для установки сообщения
// swagger:route POST /set setMessage
//
// Установить новое сообщение
//
// Запрос должен содержать:
//   - name: body
//     in: body
//     description: Запрос для установки сообщения
//     required: true
//     schema:
//       $ref: '#/components/schemas/SetRequest'
//
//     responses:
//       200:
//         description: Сообщение успешно установлено
//         content:
//           application/json:
//             schema:
//               $ref: '#/components/schemas/SetResponse'
//       400:
//         description: Неверные данные запроса
//       500:
//         description: Ошибка сервера
func NewSetter(messageSetter messageSetter, produce Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		id, err := messageSetter.Insert(context.Background(), req.Text)
		if err != nil {
			log.Printf("Failed to insert data: %v", err)
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		err = produce.ProduceMessage(id, req.Text)
		if err != nil {
			log.Printf("Failed to send to Kafka: %v", err)
			http.Error(w, "Failed to send to Kafka", http.StatusInternalServerError)
			return
		}

		res := SetResponse{
			ID:      id,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
