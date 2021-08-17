package datamodels

type Product struct {
	ID           int64   `json:"id" sql:"ID" sales:"ID"`
	ProductName  string  `json:"ProductName" sql:"productName" sales:"ProductName"`
	ProductPrice float64 `json:"ProductPrice" sql:"productPrice" sales:"ProductPrice"`
	ProductNum   int64   `json:"ProductNum" sql:"productNum" sales:"ProductNum"`
	ProductImage string  `json:"ProductImage" sql:"productImage" sales:"ProductImage"`
	ProductUrl   string  `json:"ProductUrl" sql:"productUrl" sales:"ProductUrl"`
}
