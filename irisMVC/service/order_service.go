package service

import (
	"go-flash-sale/irisMVC/datamodels"
	"go-flash-sale/irisMVC/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(order *datamodels.Order) error
	GetAllOrder()([]*datamodels.Order, error)
	InsertOrder(order *datamodels.Order) (int64, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error)
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: repository}
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOrderByID (orderId int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectByKey(orderId)
}


func (o *OrderService) DeleteOrderByID (orderId int64) (isOk bool)  {
	isOk = o.OrderRepository.Delete(orderId)
	return
}

func (o *OrderService) UpdateOrder (order *datamodels.Order) error{
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder (order *datamodels.Order) (orderID int64, err error)  {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order,error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.OrderRepository.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(
	message *datamodels.Message) (orderId int64, err error) {
	order := &datamodels.Order{
		UserId: message.UserID,
		ProductId: message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}
