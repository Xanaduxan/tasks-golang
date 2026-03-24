package websocket

import (
	"sync"

	"github.com/google/uuid"
	gws "github.com/gorilla/websocket"
)

type Manager struct {
	mu sync.RWMutex

	clients map[uuid.UUID]map[*gws.Conn]bool
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[uuid.UUID]map[*gws.Conn]bool),
	}
}

func (m *Manager) Register(userID uuid.UUID, conn *gws.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.clients[userID] == nil {
		m.clients[userID] = make(map[*gws.Conn]bool)
	}

	m.clients[userID][conn] = true
}

func (m *Manager) Unregister(userID uuid.UUID, conn *gws.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conns, ok := m.clients[userID]; ok {
		delete(conns, conn)

		if len(conns) == 0 {
			delete(m.clients, userID)
		}
	}

	_ = conn.Close()
}

func (m *Manager) SendToUser(userID uuid.UUID, message []byte) {
	m.mu.RLock()

	conns, ok := m.clients[userID]
	if !ok {
		m.mu.RUnlock()
		return
	}

	var toRemove []*gws.Conn

	for conn := range conns {
		if err := conn.WriteMessage(gws.TextMessage, message); err != nil {
			toRemove = append(toRemove, conn)
		}
	}

	m.mu.RUnlock()

	for _, conn := range toRemove {
		m.Unregister(userID, conn)
	}
}
