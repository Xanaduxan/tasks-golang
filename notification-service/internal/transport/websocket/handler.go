package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	gws "github.com/gorilla/websocket"
)

type Handler struct {
	manager  *Manager
	upgrader gws.Upgrader
}

func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
		upgrader: gws.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error:", err)
		return
	}

	h.manager.Register(userID, conn)
	defer h.manager.Unregister(userID, conn)

	log.Printf("websocket client connected: user_id=%s", userID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("websocket read error:", err)
			break
		}

		log.Printf("received from user %s: %s", userID, string(msg))
	}
}
