package main

import (
	"sync"
	"time"
)

type InMemoryDatabase struct {
	clientMut, transactionMut *sync.RWMutex
	clients                   map[int]Client
	transactions              map[int][]Transaction
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		clientMut:      &sync.RWMutex{},
		transactionMut: &sync.RWMutex{},
		clients:        map[int]Client{},
		transactions:   map[int][]Transaction{},
	}
}

func (d *InMemoryDatabase) FindClient(id int) (Client, error) {
	d.clientMut.RLock()
	defer d.clientMut.RUnlock()
	client, ok := d.clients[id]
	if !ok {
		return Client{}, ErrNotFound
	}
	return client, nil
}

func (d *InMemoryDatabase) UpdateClient(c Client) error {
	d.clientMut.Lock()
	defer d.clientMut.Unlock()
	d.clients[c.ID] = c
	return nil
}

func (d *InMemoryDatabase) UpdateTransactions(id int, t Transaction) error {
	d.transactionMut.Lock()
	defer d.transactionMut.Unlock()
	now := time.Now()
	t.RealizadaEm = &now
	d.transactions[id] = append(d.transactions[id], t)
	return nil
}
