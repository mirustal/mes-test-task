package message

import (
	"context"
	"net/http"
)

type MarkResponse struct {
	Err      int    `json:"error"`
	Messages string `json:"messages"`
}

type markDB interface {
	MarkAsRead(context.Context, string) error
}

type kafkaConsumer interface {
	Signal()
}

func NewMark(db markDB, ch kafkaConsumer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		// msg, err := ch.Si
		// if err != nil {
		// 	http.Error(w, fmt.Sprintf("Failed to get message: %v", err), http.StatusInternalServerError)
		// 	return
		// }

		// // Предполагаем, что ID сообщения находится в поле `ID` после декодирования сообщения
		// var messageID string
		// if err := json.Unmarshal(msg.Value, &messageID); err != nil {
		// 	http.Error(w, fmt.Sprintf("Failed to decode message: %v", err), http.StatusInternalServerError)
		// 	return
		// }

		// // Вызов метода базы данных для пометки сообщения как прочитанного
		// ctx := context.Background() 
		// if err := db.MarkAsRead(ctx, messageID); err != nil {
		// 	http.Error(w, fmt.Sprintf("Failed to mark as read: %v", err), http.StatusInternalServerError)
		// 	return
		// }

		// // Отправляем успешный ответ
		// response := MarkResponse{
		// 	Err:      0,
		// 	Messages: messageID,
		// }
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(response)
	}
}
