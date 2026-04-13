package chain

import "fmt"

const Ethereum ID = "ethereum"

type Registry struct {
	clients map[ID]Client
}

func NewRegistry(clients ...Client) *Registry {
	r := &Registry{clients: make(map[ID]Client)}
	for _, c := range clients {
		r.clients[c.ChainID()] = c
	}
	return r
}

func (r *Registry) Client(id ID) (Client, error) {
	c, ok := r.clients[id]
	if !ok {
		return nil, fmt.Errorf("chain client not registered: %s", id)
	}
	return c, nil
}
