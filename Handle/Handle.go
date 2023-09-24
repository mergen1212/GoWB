package Handle

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"io"
	"log"
	"net/http"
	"prisma/dto"
	"prisma/postgres"
	"prisma/prisma-shema/db"
	"strings"
)

func GetItemsToTrek(conn stan.Conn, client *db.PrismaClient, ctx context.Context, c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Замените на адрес вашего клиентского приложения

		// Разрешаем определенные методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Разрешаем определенные заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Разрешаем отправку куки
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method != http.MethodGet {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(r.URL.Path, "/")
		userID := parts[2]

		TrackNumber := dto.TrackNumber{TrackNumber: userID}

		if Items, found := c.Get(TrackNumber.TrackNumber); found {
			// Преобразуем объект к нужному типу
			cacheItems, ok := Items.([]db.ItemsModel)
			if !ok {
				log.Printf("err Order mem type")
			}
			log.Printf("OK %v", cacheItems)
			Ord, _ := json.Marshal(&Items)
			w.Write(Ord)
		} else {
			callback := func(msg *stan.Msg) {
				err, data := postgres.GetItems(client, ctx, TrackNumber.TrackNumber)
				if err != nil {
					log.Printf("db not GetItems")
				}
				c.Set(TrackNumber.TrackNumber, data, cache.DefaultExpiration)
			}
			_, err := conn.Subscribe(TrackNumber.TrackNumber, callback)
			if err != nil {
				log.Fatalf("Ошибка подписки на тему: %s", err)
			}
			BP, err := json.Marshal(TrackNumber.TrackNumber)
			if err != nil {
				log.Fatalf("Ошибка отправки сообщения: %v", err)
			}
			err = conn.Publish(TrackNumber.TrackNumber, BP)
			err, data := postgres.GetItems(client, ctx, TrackNumber.TrackNumber)
			if err != nil {
				log.Printf("db not GetItems")
			}
			Ord, _ := json.Marshal(data)
			w.Write(Ord)
		}

	}
}

func GetIdToDelivery(client *db.PrismaClient, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело POST-запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		// Закрываем тело запроса после чтения
		defer r.Body.Close()
		// Создаем структуру для размаршализации JSON
		var DeliveryId dto.DeliveryUid
		// Размаршализуем JSON из тела запроса
		if err := json.Unmarshal(body, &DeliveryId); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		//log.Printf("DeliveryId: %v", DeliveryId)
		//err, data := postgres.GetDelivery(client, ctx, DeliveryId.DeliveryUid)
		//Ord, _ := json.Marshal(&data)
		//w.Write(Ord)

	}
}

func GetTransactionToPayment(client *db.PrismaClient, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело POST-запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		// Закрываем тело запроса после чтения
		defer r.Body.Close()
		// Создаем структуру для размаршализации JSON
		var Transaction dto.Transaction
		// Размаршализуем JSON из тела запроса
		if err := json.Unmarshal(body, &Transaction); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		log.Printf("DeliveryId: %v", Transaction)
		err, data := postgres.GetPayment(client, ctx, Transaction.Transaction)
		Ord, _ := json.Marshal(&data)
		w.Write(Ord)
	}
}

func NatsOrderPost(conn stan.Conn, client *db.PrismaClient, ctx context.Context, c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}
		// Читаем тело POST-запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		// Закрываем тело запроса после чтения
		defer r.Body.Close()
		// Создаем структуру для размаршализации JSON
		var Order dto.Order
		//var Order dto.Order
		// Размаршализуем JSON из тела запроса
		if err := json.Unmarshal(body, &Order); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		if Items, found := c.Get(Order.OrderUid); found {
			cacheItems, _ := Items.(dto.Order)
			log.Printf("OK %v", cacheItems)
			Ord, _ := json.Marshal(&Items)
			w.Write(Ord)
		} else {
			callback := func(msg *stan.Msg) {
				if err := json.Unmarshal(msg.Data, &Order); err != nil {
					http.Error(w, "Error decoding JSON", http.StatusBadRequest)
					return
				}
				postgres.CreateOrderOne(client, ctx, Order)
				postgres.CreatePlaymentOne(client, ctx, Order)
				postgres.CreateDeliveryOne(client, ctx, Order)
				postgres.CreateOneItem(client, ctx, Order)
				if err != nil {
					log.Printf("db not GetItdgfhdfhems")
				}
				err, data := postgres.GetOrderUid(client, ctx, Order.OrderUid)
				if err != nil {
					log.Printf("db not GetItems")
				}

				c.Set(Order.OrderUid, data, cache.DefaultExpiration)
			}
			_, err = conn.Subscribe(Order.OrderUid, callback)
			if err != nil {
				log.Fatalf("Ошибка подписки на тему: %s", err)
			}
			BP, err := json.Marshal(Order)
			if err != nil {
				log.Fatalf("Ошибка отправки сообщения: %v", err)
			}
			err = conn.Publish(Order.OrderUid, BP)
			w.Write(BP)
		}
	}
}

func NatsGetOrderUid(c *cache.Cache, client *db.PrismaClient, ctx context.Context, conn stan.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Замените на адрес вашего клиентского приложения

		// Разрешаем определенные методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Разрешаем определенные заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Разрешаем отправку куки
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method != http.MethodGet {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(r.URL.Path, "/")
		userID := parts[2]

		Order := dto.OrderUd{OrderUid: userID}
		//var Order dto.Order
		// Размаршализуем JSON из тела запроса

		log.Printf(Order.OrderUid)
		if Items, found := c.Get(Order.OrderUid); found {
			cacheItems, _ := Items.(db.OrderModel)
			log.Printf("OK %v", cacheItems)
			Ord, _ := json.Marshal(&Items)
			w.Write(Ord)
		} else {
			callback := func(msg *stan.Msg) {
				if err := json.Unmarshal(msg.Data, &Order); err != nil {
					http.Error(w, "Error decoding JSON", http.StatusBadRequest)
					return
				}
				err, data := postgres.GetOrderUid(client, ctx, Order.OrderUid)
				if err != nil {
					log.Printf("db not GetItems")
				}
				c.Set(Order.OrderUid, data, cache.DefaultExpiration)
			}
			_, err := conn.Subscribe(Order.OrderUid, callback)
			if err != nil {
				log.Fatalf("Ошибка подписки на тему: %s", err)
			}
			BP, err := json.Marshal(Order)
			if err != nil {
				log.Fatalf("Ошибка отправки сообщения: %v", err)
			}
			err = conn.Publish(Order.OrderUid, BP)
			if err != nil {
				log.Printf("fdgdfhfjkghfdfghjhgfdghj %v", err)
			}
			err, data := postgres.GetOrderUid(client, ctx, Order.OrderUid)
			if err != nil {
				log.Printf("db not GetItems")
			}
			Org, err := json.Marshal(data)
			w.Write(Org)
		}
	}
}
