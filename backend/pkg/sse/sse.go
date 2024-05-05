package sse

import (
	"sync"
)

type SSEConn struct {
	mu      sync.Mutex
	clients map[string][]chan string
}

func NewSSEConn() *SSEConn {
	return &SSEConn{clients: make(map[string][]chan string)}
}

func (p *SSEConn) AddClient(id string) *chan string {
	p.mu.Lock()

	defer func() {
		p.mu.Unlock()
	}()

	c, ok := p.clients[id]
	if !ok {
		client := []chan string{make(chan string)}
		p.clients[id] = client
		return &client[0]
	}

	newCh := make(chan string)
	p.clients[id] = append(c, newCh)
	return &newCh
}

func (p *SSEConn) RemoveClient(id string, conn chan string) {
	p.mu.Lock()

	defer func() {
		p.mu.Unlock()
	}()

	c, ok := p.clients[id]
	if !ok {
		return
	}

	pos := -1
	for i, ch := range c {
		if ch == conn {
			pos = i
		}
	}
	if pos == -1 {
		return
	}
	close(c[pos])
	// c = append(c[:pos], c[pos+1:]...)
	if pos == 0 {
		delete(p.clients, id)
	}
}

func (p *SSEConn) Broadcast(id string, data string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	c, ok := p.clients[id]
	if !ok {
		return
	}

	for _, ch := range c {
		ch <- data
	}
}
