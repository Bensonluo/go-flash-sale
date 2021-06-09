package repositories

import (
	"database/sql"
	"go-flash-sale/irisMVC/common"
	"go-flash-sale/irisMVC/datamodels"
	"strconv"
)


type IOrderRepository interface{
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey (int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}


func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
		return &OrderManagerRepository{table: table, mysqlConn: sql}
}


type OrderManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productId int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}

	sqlS := "INSERT " + o.table + " set userId=?,productId=?,orderStatus="
	stmt, errStmt := o.mysqlConn.Prepare(sqlS)
	if errStmt != nil {
		return productId, err
	}
	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)

	if errResult != nil {
		return productId, errResult
	}

	return result.LastInsertId()
}

func (o *OrderManagerRepository) Delete(orderId int64) (isOk bool) {
	if err := o.Conn(); err != nil {
		return
	}
	sqlS := "delete from " + o.table + " where ID=?"
	stmt, errStmt := o.mysqlConn.Prepare(sqlS)
	if errStmt != nil {
		return
	}
	_, err := stmt.Exec(orderId)
	if err != nil {
		return
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if errConn := o.Conn(); errConn != nil {
		return errConn
	}

	sqlS := "Update " + o.table + " set userTd=?,productId=?,orderStatus=? Where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(sqlS)
	if errStmt != nil {
		return errStmt
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)

	return
}

func (o *OrderManagerRepository) SelectByKey (orderId int64) (order *datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}
	sqlS := "Select * From " + o.table + " where ID=" + strconv.FormatInt(orderId, 10)

	row, errRow := o.mysqlConn.Query(sqlS)
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManagerRepository)SelectAll ()(orderArray []*datamodels.Order,err error)  {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sqlS := "Select * from " + o.table
	rows ,errRows := o.mysqlConn.Query(sqlS)
	defer rows.Close()
	if errRows != nil {
		return nil, errRows
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil ,err
	}

	for _, v := range result{
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sqlS := "Select o.ID,p.productName,o.orderStatus From sales.order as o left join product as p on o.productId=p.ID"
	rows, errRows := o.mysqlConn.Query(sqlS)
	defer rows.Close()
	if errRows != nil {
		return nil, errRows
	}
	return common.GetResultRows(rows), err
}
