package main

type Database interface {
	FindClient(id int) (Client, error)
	UpdateClient(c Client) error
	UpdateTransactions(id int, t Transaction) error
}
