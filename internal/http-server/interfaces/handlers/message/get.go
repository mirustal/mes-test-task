package message

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"kafka-app/internal/domain"
)



// GetResponse представляет ответ на запрос получения сообщений
// swagger:response GetResponse
type GetResponse struct {
	// Код ошибки
    // in: body
	Err      int              `json:"error"`
	// Список сообщений
    // in: body
	Messages []domain.Message `json:"messages"`
}

type MessageGetter interface {
	GetUser(ctx context.Context, name string, limit int) ([]domain.Message, error)
}



// NewGetter создает новый обработчик для получения сообщений
// swagger:route GET /get getMessages
//
// Получить сообщения по ID и лимиту
//
// Список параметров запроса:
//   - name: id
//     in: query
//     description: ID сообщения
//     required: true
//     schema:
//       type: string
//   - name: limit
//     in: query
//     description: Лимит сообщений для возвращения
//     required: true
//     schema:
//       type: integer
//       format: int32
//
//     responses:
//       200:
//         description: Успешный ответ
//         content:
//           application/json:
//             schema:
//               $ref: '#/components/schemas/GetResponse'
//       400:
//         description: Неверные параметры запроса
//       404:
//         description: Не удалось получить данные
func NewGetter(messageGetter MessageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		limitStr := r.URL.Query().Get("limit")
		
		if id == "" || limitStr == "" {
			http.Error(w, "Missing id or limit parameter", http.StatusBadRequest)
			return
		}


		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}


		messages, err := messageGetter.GetUser(context.Background(), id, limit)
		if err != nil {
			log.Printf("Failed to get data: %v", err)
			http.Error(w, "Failed to get data", http.StatusNotFound)
			return
		}


		res := GetResponse{
			Err:      http.StatusOK,
			Messages: messages,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}