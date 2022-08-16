package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	OrderId   int64  `gorm:"Not Null" json:"orderid"`
	TotalPaid int64  `json:"totalpaid"`
	ReceiptId string `json:"receiptid"`
	CashierId int64  `json:"cashierid"`
	PaymentId int64  `json:"paymentid"`
}

type OrderProducts struct {
	gorm.Model
	ProductId int64 `gorm:"Not Null" json:"productid"`
	Qty       int64 `json:"qty"`
	Orderid   int64 `json:"orderid"`
}

type ReqCreateOrder struct {
	PaymentId int64 `json:"paymentid"`
	TotalPaid int64 `json:"totalpaid"`
	Products  []struct {
		ProductId int64 `json:"productid"`
		Qty       int64 `json:"qty"`
	} `json:"products"`
}

type CreateOrderStruct struct {
	CashiersId     int64  `json:"cashiersId"`
	PaymentTypesId int64  `json:"paymentTypesId"`
	TotalPrice     int64  `json:"totalPrice"`
	TotalPaid      int64  `json:"totalPaid"`
	TotalReturn    int64  `json:"totalReturn"`
	ReceiptId      string `json:"receiptId"`
}

type CreateProductStruct struct {
	ProductId        int64  `json:"productId"`
	Name             string `json:"name"`
	Price            int64  `json:"price"`
	Qty              int64  `json: "qty"`
	TotalNormalPrice int64  `json:"totalNormalPrice"`
	TotalFinalPrice  int64  `json:"TotalFinalPrice"`
}

type ResCreateOrder struct {
	CreatedOrder    CreateOrderStruct     `json:"order"`
	CreatedProducts []CreateProductStruct `json:"products"`
}

type OrderList struct {
	OrderId       int64  `json:"orderId"`
	CashiersId    int64  `json:"cashiersId"`
	PaymentTypeId int64  `json:"paymentTypeId"`
	TotalPrice    int64  `json:"totalPrice"`
	TotalPaid     int64  `json:"totalPaid"`
	TotalReturn   int64  `json:"totalReturn"`
	ReceiptId     string `json:"receiptId"`
	Cashier       struct {
		CashierId int64  `json:"cashierId"`
		Name      string `json:"name"`
	}
	PaymentType struct {
		PaymentTypeId int64  `json:"paymentTypeId"`
		Name          string `json:"name"`
		Logo          string `json:"logo"`
		Type          string `json:"type"`
	}
}

type ReqDetailSubtotalOrder struct {
	ProductId int64 `json:"productid"`
	Qty       int64 `json:"qty"`
}

type ResDetailProduct struct {
	ProductId        int64  `json:"productId"`
	Name             string `json:"name"`
	Price            int64  `json:"price"`
	Qty              int64  `json: "qty"`
	TotalNormalPrice int64  `json:"totalNormalPrice"`
	TotalFinalPrice  int64  `json:"TotalFinalPrice"`
}

type ResOrderDetail struct {
	DetailOrder struct {
		OrderId       int64  `json:"orderId"`
		CashiersId    int64  `json:"cashiersId"`
		PaymentTypeId int64  `json:"paymentTypeId"`
		TotalPrice    int64  `json:"totalPrice"`
		TotalPaid     int64  `json:"totalPaid"`
		TotalReturn   int64  `json:"totalReturn"`
		ReceiptId     string `json:"receiptId"`
		Cashier       struct {
			CashierId int64  `json:"cashierId"`
			Name      string `json:"name"`
		}
		PaymentType struct {
			PaymentTypeId int64  `json:"paymentTypeId"`
			Name          string `json:"name"`
			Logo          string `json:"logo"`
			Type          string `json:"type"`
		}
	} `json:"order"`
	DetailProducts []ResDetailProduct `json:"products"`
}

type ResSubTotalProduct struct {
	ProductId        int64   `json:"productid"`
	Name             string  `json:"name"`
	Stock            int64   `json:"stock"`
	Price            float64 `json:"price"`
	Image            string  `json:"image"`
	Sku              string  `json:"SKU"`
	Qty              int64   `json:qty`
	CategoryId       int64   `json:"categoryId"`
	DiscountId       int64   `json:"discount"`
	TotalNormalPrice int64   `json:"totalNormalPrice"`
	TotalFinalPrice  int64   `json:"totalFinalPrice"`
}
type ResSubTotal struct {
	SubTotal    int64                `json:"subtotal"`
	SubProducts []ResSubTotalProduct `json:"products"`
}

func FindOrderProductsByOrderId(id int64) []OrderProducts {
	var orderProducts []OrderProducts

	db := GetDB()
	err := GetDB().Where("orderid = ?", id).First(&orderProducts).Error
	if err != nil {
		return orderProducts
	}

	db.Find(&orderProducts)

	return orderProducts
}

func FindSubTotal(reqSubtotalOrder []ReqDetailSubtotalOrder) ResSubTotal {
	var resSubTotal ResSubTotal
	resSubTotal.SubProducts = make([]ResSubTotalProduct, len(reqSubtotalOrder))
	var subtotal int64 = 0
	for i := 0; i < len(reqSubtotalOrder); i++ {
		products, _ := FindProduct(int(reqSubtotalOrder[i].ProductId))
		resSubTotal.SubProducts[i].CategoryId = products.CategoryId
		resSubTotal.SubProducts[i].DiscountId = products.DiscountId
		resSubTotal.SubProducts[i].Image = products.Image
		resSubTotal.SubProducts[i].Name = products.Name
		resSubTotal.SubProducts[i].Price = products.Price
		resSubTotal.SubProducts[i].ProductId = products.ProductId
		resSubTotal.SubProducts[i].Sku = products.Sku
		resSubTotal.SubProducts[i].Stock = products.Stock
		resSubTotal.SubProducts[i].Qty = reqSubtotalOrder[i].Qty
		resSubTotal.SubProducts[i].TotalNormalPrice = int64(products.Price) * reqSubtotalOrder[i].Qty
		resSubTotal.SubProducts[i].TotalFinalPrice = int64(products.Price) * reqSubtotalOrder[i].Qty
		subtotal += resSubTotal.SubProducts[i].TotalNormalPrice
	}
	resSubTotal.SubTotal = subtotal
	return resSubTotal
}

func FindAllOrders(c *fiber.Ctx) []OrderList {
	var orders []Order

	db := GetDB()
	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}

	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}

	db.Find(&orders)
	var orderList = make([]OrderList, len(orders))

	for i := 0; i < len(orders); i++ {
		orderProducts := FindOrderProductsByOrderId(orders[i].OrderId)
		totalPrice := 0
		for j := 0; j < len(orderProducts); j++ {
			products, getProductErr := FindProduct(int(orderProducts[i].ProductId))
			if getProductErr == nil {
				totalPrice += int(products.Price) * int(orderProducts[j].Qty)
			}
		}
		cashier, _ := FindCashier(int(orders[i].CashierId))
		payment, _ := FindPayment(int(orders[i].PaymentId))
		orderList[i].OrderId = orders[i].OrderId
		orderList[i].CashiersId = orders[i].CashierId
		orderList[i].PaymentTypeId = orders[i].PaymentId
		orderList[i].TotalPaid = orders[i].TotalPaid
		orderList[i].TotalPrice = int64(totalPrice)
		orderList[i].TotalReturn = orders[i].TotalPaid - int64(totalPrice)
		orderList[i].ReceiptId = orders[i].ReceiptId
		orderList[i].Cashier.CashierId = cashier.CashierId
		orderList[i].Cashier.Name = cashier.Name
		orderList[i].PaymentType.PaymentTypeId = payment.PaymentId
		orderList[i].PaymentType.Name = payment.Name
		orderList[i].PaymentType.Logo = payment.Logo
		orderList[i].PaymentType.Type = payment.Type
	}

	return orderList
}

func GetOrderCount(c *fiber.Ctx) int64 {
	var orders []Order

	db := GetDB()

	count := db.Find(&orders).RowsAffected

	return count
}

func FindOrder(id int, reqDetailOrder []ReqDetailSubtotalOrder) (ResOrderDetail, error) {
	var order Order

	err := GetDB().Where("order_id = ?", id).First(&order).Error
	if err != nil {
		err = errors.New("order not found")
	}

	orderProducts := FindOrderProductsByOrderId(int64(id))
	var resOrderDetail ResOrderDetail
	var orderDetailProducts []OrderProducts
	for i := 0; i < len(reqDetailOrder); i++ {
		for j := 0; j < len(orderProducts); j++ {
			if reqDetailOrder[i].ProductId == orderProducts[j].ProductId &&
				reqDetailOrder[i].Qty == orderProducts[j].Qty {
				orderDetailProducts = append(orderDetailProducts, orderProducts[j])
				break
			}
		}
	}
	resOrderDetail.DetailOrder.OrderId = order.OrderId
	resOrderDetail.DetailOrder.CashiersId = order.CashierId
	resOrderDetail.DetailOrder.PaymentTypeId = order.PaymentId
	resOrderDetail.DetailOrder.TotalPaid = order.TotalPaid

	resOrderDetail.DetailProducts = make([]ResDetailProduct, len(orderDetailProducts))
	for i := 0; i < len(orderDetailProducts); i++ {
		products, getProductErr := FindProduct(int(orderDetailProducts[i].ProductId))
		if getProductErr == nil {
			resOrderDetail.DetailOrder.TotalPrice += int64(products.Price) * orderProducts[i].Qty
			resOrderDetail.DetailProducts[i].ProductId = products.ProductId
			resOrderDetail.DetailProducts[i].Name = products.Name
			resOrderDetail.DetailProducts[i].Price = int64(products.Price)
			resOrderDetail.DetailProducts[i].Qty = orderDetailProducts[i].Qty
			resOrderDetail.DetailProducts[i].TotalNormalPrice = int64(products.Price) * orderDetailProducts[i].Qty
			resOrderDetail.DetailProducts[i].TotalFinalPrice = int64(products.Price) * orderDetailProducts[i].Qty
		}
	}
	resOrderDetail.DetailOrder.TotalReturn = resOrderDetail.DetailOrder.TotalPaid - resOrderDetail.DetailOrder.TotalPrice
	resOrderDetail.DetailOrder.ReceiptId = order.ReceiptId
	cashier, _ := FindCashier(int(order.CashierId))
	payment, _ := FindPayment(int(order.PaymentId))
	resOrderDetail.DetailOrder.Cashier.CashierId = cashier.CashierId
	resOrderDetail.DetailOrder.Cashier.Name = cashier.Name
	resOrderDetail.DetailOrder.PaymentType.PaymentTypeId = payment.PaymentId
	resOrderDetail.DetailOrder.PaymentType.Type = payment.Type
	resOrderDetail.DetailOrder.PaymentType.Name = payment.Name
	resOrderDetail.DetailOrder.PaymentType.Logo = payment.Logo

	return resOrderDetail, err
}

func CreateOrder(reqOrder ReqCreateOrder) (ResCreateOrder, error) {
	var maxOrder Order
	var payment Payments
	order := new(Order)
	var orderProducts = make([]OrderProducts, len(reqOrder.Products))

	GetDB().Raw(`
		SELECT COALESCE(MAX(order_id) + 1, 1) as order_id
		FROM orders
		`).Scan(
		&maxOrder,
	)

	order.CashierId = 1
	order.PaymentId = reqOrder.PaymentId
	order.TotalPaid = reqOrder.TotalPaid
	order.OrderId = maxOrder.OrderId
	order.ReceiptId = "S1001"

	var resCreateOrder ResCreateOrder
	resCreateOrder.CreatedOrder.PaymentTypesId = order.PaymentId
	resCreateOrder.CreatedOrder.CashiersId = order.CashierId
	resCreateOrder.CreatedOrder.TotalPaid = order.TotalPaid
	count := GetDB().Table("payments").Where("payment_id = ?", reqOrder.PaymentId).Find(&payment).RowsAffected
	if count == 0 {
		return resCreateOrder, errors.New("PaymentType not found")
	}

	orderErr := GetDB().Create(order).Error
	if orderErr != nil {
		return resCreateOrder, orderErr
	}
	resCreateOrder.CreatedProducts = make([]CreateProductStruct, len(reqOrder.Products))
	for i := 0; i < len(reqOrder.Products); i++ {
		orderProducts[i].ProductId = reqOrder.Products[i].ProductId
		orderProducts[i].Qty = reqOrder.Products[i].Qty
		orderProducts[i].Orderid = order.OrderId
		orderProductErr := GetDB().Create(&orderProducts[i]).Error
		if orderProductErr != nil {
			return resCreateOrder, orderProductErr
		}

		products, getProductErr := FindProduct(int(orderProducts[i].ProductId))
		if getProductErr == nil {
			resCreateOrder.CreatedOrder.TotalPrice += int64(products.Price) * orderProducts[i].Qty
			resCreateOrder.CreatedProducts[i].ProductId = products.ProductId
			resCreateOrder.CreatedProducts[i].Name = products.Name
			resCreateOrder.CreatedProducts[i].Price = int64(products.Price)
			resCreateOrder.CreatedProducts[i].Qty = orderProducts[i].Qty
			resCreateOrder.CreatedProducts[i].TotalNormalPrice = int64(products.Price) * orderProducts[i].Qty
			resCreateOrder.CreatedProducts[i].TotalFinalPrice = int64(products.Price) * orderProducts[i].Qty
		}
	}
	resCreateOrder.CreatedOrder.TotalReturn = resCreateOrder.CreatedOrder.TotalPaid - resCreateOrder.CreatedOrder.TotalPrice
	return resCreateOrder, nil
}
