package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	db := NewInMemoryDatabase()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /clientes/{id}/transacoes", handleTransaction(db))
}

func handleTransaction(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid ID"))
			return
		}

		c, err := db.FindClient(id)
		if err == ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Client not found"))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Oops, something went wrong!"))
			return
		}

		t := Transaction{}
		err = json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid transaction"))
			return
		}

		if t.Tipo == "c" {
			c.Total += t.Valor
		} else if hasEnoughMoney(c.Limite, c.Total, t.Valor) {
			c.Total -= t.Valor
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Insufficient funds"))
			return
		}

		err = db.UpdateClient(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = db.UpdateTransactions(id, t)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		response := fmt.Sprintf(`{"limite": %d, "saldo": %d}`, c.Limite, c.Total)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

func hasEnoughMoney(limite, total, valor int) bool {
	return limite+(total-valor) >= 0
}
