package Models

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type OrderInfo struct {
	orderId           string  // 訂單編號
	orderStatus       string  // 訂單狀態
	orderReturn       string  // 退貨 / 退款狀態
	buyerAccount      string  // 買家帳號
	orderTime         string  // 訂單成立時間
	payTime           string  // 買家完成付款時間
	orderAmount       float64 // 訂單小計 (TWD)
	buyerFreight      float64 // 買家支付的運費
	totalPay          float64 // 訂單總金額
	shopeeCoin        float64 // 蝦幣折抵
	shopeeDiscount    float64 // 蝦皮發放折扣券
	sellerDiscount    float64 // 賣家自設折扣券
	buyDetail         string  // 商品資訊
	tmpA              string  // 促銷組合指標
	tmpB              string  // 蝦皮促銷組合折扣
	recvAddr          string  // 收件地址
	country           string  // 國家
	city              string  // 城市
	district          string  // 行政區
	postalCode        string  // 郵遞區號
	buyerName         string  // 收件者姓名
	phone             string  // 電話號碼
	shippingMethod    string  // 寄送方式
	ShipmentMethod    string  // 出貨方式
	orderType         string  // 訂單類型
	payMethod         string  // 付款方式
	ccLast4           string  // 信用卡後四碼
	lastShippingTime  string  // 最晚出貨日期
	trackingNum       string  // 包裹查詢號碼
	realShippingTime  string  // 實際出貨時間
	orderCompleteTime string  // 訂單完成時間
	buyerComment      string  // 買家備註
	comment           string  // 備註
}

//convert BIG5 to UTF-8    https://gist.github.com/zhangbaohe/c691e1da5bbdc7f41ca5
func Decodebig5(strSource string, strOutput *string) error {
	s := []byte(strSource)
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return e
	}
	*strOutput = string(d)
	return nil
}

//convert UTF-8 to BIG5
func Encodebig5(strSource string, strOutput *string) error {
	s := []byte(strSource)
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return e
	}
	*strOutput = string(d)
	return nil
}

func SalesImport(csvFilePath string) {

	csvFile, _ := os.Open(csvFilePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var orderInfo []OrderInfo
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		for i := 0; i < len(line); i++ {
			Decodebig5(line[i], &line[i])
		}

		oneOrder := OrderInfo{
			orderId:           line[0],
			orderStatus:       line[1],
			orderReturn:       line[2],
			buyerAccount:      line[3],
			orderTime:         line[4],
			payTime:           line[5],
			orderAmount:       0,
			buyerFreight:      0,
			totalPay:          0,
			shopeeCoin:        0,
			shopeeDiscount:    0,
			sellerDiscount:    0,
			buyDetail:         line[12],
			tmpA:              line[13],
			tmpB:              line[14],
			recvAddr:          line[15],
			country:           line[16],
			city:              line[17],
			district:          line[18],
			postalCode:        line[19],
			buyerName:         line[20],
			phone:             line[21],
			shippingMethod:    line[22],
			ShipmentMethod:    line[23],
			orderType:         line[24],
			payMethod:         line[25],
			ccLast4:           line[26],
			lastShippingTime:  line[27],
			trackingNum:       line[28],
			realShippingTime:  line[29],
			orderCompleteTime: line[30],
			buyerComment:      line[31],
			comment:           line[32],
		}
		oneOrder.orderAmount, _ = strconv.ParseFloat(line[6], 32)
		oneOrder.buyerFreight, _ = strconv.ParseFloat(line[7], 32)
		oneOrder.totalPay, _ = strconv.ParseFloat(line[8], 32)
		oneOrder.shopeeCoin, _ = strconv.ParseFloat(line[9], 32)
		oneOrder.shopeeDiscount, _ = strconv.ParseFloat(line[10], 32)
		oneOrder.sellerDiscount, _ = strconv.ParseFloat(line[11], 32)
		orderInfo = append(orderInfo, oneOrder)

		fmt.Println(line[0], line[3], line[12])
	}

}
