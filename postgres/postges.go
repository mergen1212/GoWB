package postgres

import (
	"context"
	"log"
	"prisma/dto"
	"prisma/prisma-shema/db"
)

func CreateOrderOne(client *db.PrismaClient, ctx context.Context, Order dto.Order) error {
	_, err := client.Order.CreateOne(
		db.Order.OrderUId.Set(Order.OrderUid),
		db.Order.TrackNumber.Set(Order.TrackNumber),
		db.Order.Entry.Set(Order.Entry),
		db.Order.DeliveryID.Set(Order.DeliveryId),
		db.Order.Locale.Set(Order.Locale),
		db.Order.InternalSignature.Set(Order.InternalSignature),
		db.Order.CustomerID.Set(Order.CustomerId),
		db.Order.DeliveryService.Set(Order.DeliveryService),
		db.Order.Shardkey.Set(Order.Shardkey),
		db.Order.SmID.Set(Order.SmId),
		db.Order.OofShard.Set(Order.OofShard),
	).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CreateDeliveryOne(client *db.PrismaClient, ctx context.Context, Order dto.Order) error {
	_, err := client.Delivery.CreateOne(
		db.Delivery.Name.Set(Order.Delivery.Name),
		db.Delivery.Phone.Set(Order.Delivery.Phone),
		db.Delivery.Zip.Set(Order.Delivery.Zip),
		db.Delivery.City.Set(Order.Delivery.City),
		db.Delivery.Address.Set(Order.Delivery.Address),
		db.Delivery.Region.Set(Order.Delivery.Region),
		db.Delivery.Email.Set(Order.Delivery.Email),
		db.Delivery.Order.Link(
			db.Order.DeliveryID.Equals(Order.DeliveryId),
		),
	).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CreatePlaymentOne(client *db.PrismaClient, ctx context.Context, Order dto.Order) error {
	_, err := client.Payment.CreateOne(
		db.Payment.RequestID.Set(Order.Payment.RequestId),
		db.Payment.Provider.Set(Order.Payment.Provider),
		db.Payment.Amount.Set(Order.Payment.Amount),
		db.Payment.PaymentDt.Set(Order.Payment.PaymentDt),
		db.Payment.Bank.Set(Order.Payment.Bank),
		db.Payment.DeliveryCost.Set(Order.Payment.DeliveryCost),
		db.Payment.GoodsTotal.Set(Order.Payment.GoodsTotal),
		db.Payment.CustomFee.Set(Order.Payment.CustomFee),
		db.Payment.Order.Link(
			db.Order.OrderUId.Equals(Order.OrderUid),
		),
	).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CreateOneItem(client *db.PrismaClient, ctx context.Context, Order dto.Order) error {
	_, err := client.Items.CreateOne(
		db.Items.ChrtID.Set(Order.Items[0].ChrtId),
		db.Items.Price.Set(Order.Items[0].Price),
		db.Items.Rid.Set(Order.Items[0].Rid),
		db.Items.Name.Set(Order.Items[0].Name),
		db.Items.Sale.Set(Order.Items[0].Sale),
		db.Items.Size.Set(Order.Items[0].Size),
		db.Items.TotalPrice.Set(Order.Items[0].TotalPrice),
		db.Items.NmID.Set(Order.Items[0].NmId),
		db.Items.Brand.Set(Order.Items[0].Brand),
		db.Items.Status.Set(Order.Items[0].Status),
		db.Items.Post.Link(
			db.Order.TrackNumber.Equals(Order.Items[0].TrackNumber),
		),
	).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetItems(client *db.PrismaClient, ctx context.Context, TrackNumber string) (error, []db.ItemsModel) {
	Itemss, err := client.Items.FindMany(
		db.Items.Post.Where(
			db.Order.TrackNumber.Equals(TrackNumber),
		),
	).Exec(ctx)
	if err != nil {
		return err, nil
	}
	return nil, Itemss
}

func GetDelivery(client *db.PrismaClient, ctx context.Context, DeliveryId string) (error, *db.DeliveryModel) {
	Delivery, err := client.Delivery.FindFirst(
		db.Delivery.Order.Where(
			db.Order.DeliveryID.Equals(DeliveryId),
		),
	).Exec(ctx)
	if err != nil {
		return err, nil
	}
	return nil, Delivery
}

func GetPayment(client *db.PrismaClient, ctx context.Context, transaction string) (error, *db.PaymentModel) {
	Payment, err := client.Payment.FindFirst(
		db.Payment.Order.Where(
			db.Order.OrderUId.Equals(transaction),
		),
	).Exec(ctx)
	if err != nil {
		return err, nil
	}
	return nil, Payment
}

func GetOrderUid(client *db.PrismaClient, ctx context.Context, OrderUId string) (error, *db.OrderModel) {
	Order, err := client.Order.FindFirst(
		db.Order.OrderUId.Equals(OrderUId),
	).With(
		db.Order.Delivery.Fetch(),
		db.Order.Payment.Fetch(),
		db.Order.Items.Fetch(),
	).Exec(ctx)
	if err != nil {
		return err, nil
	}
	Items := Order.Items()
	for _, item := range Items {
		log.Printf("items: %v", item)
	}
	log.Printf("GetOrderUid: %v", Order)
	return nil, Order
}
