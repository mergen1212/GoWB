package main

import (
	"context"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"prisma/Handle"
	"prisma/prisma-shema/db"
	"time"
)

func main() {
	client := db.NewClient()
	ctx := context.Background()
	err := client.Prisma.Connect()
	if err != nil {
		log.Fatal("not connect db")
	}
	conn, err := stan.Connect("test-cluster", "test-cluster")
	if err != nil {
		log.Fatal("Not nats connected")
	}
	c := cache.New(5*time.Minute, 10*time.Minute)

	defer func() {
		client.Prisma.Disconnect()
	}()

	defer func() {
		conn.Close()
	}()
	http.HandleFunc("/NatsOrderPost", Handle.NatsOrderPost(conn, client, ctx, c))
	http.HandleFunc("/GetOrderUid/", Handle.NatsGetOrderUid(c, client, ctx, conn))
	http.HandleFunc("/Getitemstotrek/", Handle.GetItemsToTrek(conn, client, ctx, c))
	http.HandleFunc("/DeliveryId", Handle.GetIdToDelivery(client, ctx))
	http.HandleFunc("/GetTransactionToPayment", Handle.GetTransactionToPayment(client, ctx))
	http.ListenAndServe(":8080", nil)
}
