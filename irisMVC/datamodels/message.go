package datamodels

type Message struct {
	ProductID int64
	UserID    int64
}

//create struct
func NewMessage(userId int64, productId int64) *Message {
	return &Message{}
}
