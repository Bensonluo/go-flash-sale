package datamodels

type Product struct {
	ID int64 `json:"id" sql:"ID" bensonl:"id"`
	ProductName string `json:"ProductName" sql:"productName" bensonl:"ProductName"`
	ProductNum	int64 `json:"ProductNum" sql:"productNum" bensonl:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" bensonl:"ProductImage"`
	ProductUrl string `json:"ProductUrl" sql:"productUrl" bensonl:"ProductUrl"`
}
