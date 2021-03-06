package Models

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"

	"../MainApp"
)

type SalesImport struct {
	mainApp *MainApp.MainApp
}

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
	shipmentMethod    string  // 出貨方式
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

type OrderInfoBuyDetail struct {
	ItemId        string // itemId
	ItemModel     string // itemModel
	ItemModelName string
	ItemQty       int
	ItemPrice     float64
}

func (v *SalesImport) Init(mainApp *MainApp.MainApp) {
	v.mainApp = mainApp
}

// 以2關鍵字範圍取得文字內容
func (v *SalesImport) getStringKeywordRange(strSrc *string, keyword1 string, keyword2 string) (string, error) {
	idxStart := strings.Index(*strSrc, keyword1)
	if idxStart < 0 {
		return "", errors.New("Keyword1 not found! ")
	}
	idxStart += len(keyword1)

	idxEnd := strings.Index((*strSrc)[idxStart:], keyword2)
	if idxEnd < 0 {
		return "", errors.New("Keyword2 not found! ")
	}

	idxEnd += idxStart

	// fmt.Println(idxStart, idxEnd, (*strSrc)[idxStart:idxEnd])

	return (*strSrc)[idxStart:idxEnd], nil
}

//convert UTF-8 to BIG5
func (v *SalesImport) findStringInArray(strList []string, str string, cmpAll bool) int {
	for i, n := range strList {
		if cmpAll {
			if str == n {
				return i
			}
		} else {
			if strings.Contains(n, str) {
				return i
			}
		}
	}
	return len(strList)
}

func (v *SalesImport) CsvImportFromShopee(csvFilePath string) {
	v.CsvImportFromShopee_20190613(csvFilePath)
}

func (v *SalesImport) CsvImportFromShopee_20190613(csvFilePath string) {

	var strSQL string
	var errSql error

	csvFile, _ := os.Open(csvFilePath)
	csvFileBig5 := transform.NewReader(csvFile, traditionalchinese.Big5.NewDecoder()) //使用 big5 讀檔案
	reader := csv.NewReader(bufio.NewReader(csvFileBig5))
	defer csvFile.Close()

	lineFirst, errorFirst := reader.Read()
	if errorFirst == io.EOF {
		return
	} else if errorFirst != nil {
		log.Fatal(errorFirst)
	}

	idxOrderId := v.findStringInArray(lineFirst, "訂單編號", false)
	idxOrderStatus := v.findStringInArray(lineFirst, "訂單狀態", false)
	idxOrderReturn := v.findStringInArray(lineFirst, "退貨 / 退款狀態", false)
	idxBuyerAccount := v.findStringInArray(lineFirst, "買家帳號", false)
	idxOrderTime := v.findStringInArray(lineFirst, "訂單成立時間", false)
	idxPayTime := v.findStringInArray(lineFirst, "付款時間", false)
	// idxBuyDetail := v.findStringInArray(lineFirst, "商品資訊", false)
	idxRecvAddr := v.findStringInArray(lineFirst, "收件地址", false)
	idxCity := v.findStringInArray(lineFirst, "城市", false)
	idxDistrict := v.findStringInArray(lineFirst, "行政區", false)
	idxPostalCode := v.findStringInArray(lineFirst, "郵遞區號", false)
	idxBuyerName := v.findStringInArray(lineFirst, "收件者姓名", false)
	idxPhone := v.findStringInArray(lineFirst, "電話", false)
	idxShippingMethod := v.findStringInArray(lineFirst, "寄送方式", false)
	idxShipmentMethod := v.findStringInArray(lineFirst, "出貨方式", false)
	// idxOrderType := v.findStringInArray(lineFirst, "訂單類型", false)
	idxPayMethod := v.findStringInArray(lineFirst, "付款方式", false)
	idxCcLast4 := v.findStringInArray(lineFirst, "信用卡後四碼", false)
	idxLastShippingTime := v.findStringInArray(lineFirst, "最晚出貨日期", false)
	idxTrackingNum := v.findStringInArray(lineFirst, "包裹查詢號碼", false)
	idxRealShippingTime := v.findStringInArray(lineFirst, "實際出貨時間", false)
	idxOrderCompleteTime := v.findStringInArray(lineFirst, "訂單完成時間", false)
	idxBuyerComment := v.findStringInArray(lineFirst, "買家備註", false)
	// idxComment := v.findStringInArray(lineFirst, "備註", true)

	idxOrderAmount := v.findStringInArray(lineFirst, "訂單小計", false)
	idxBuyerFreight := v.findStringInArray(lineFirst, "買家支付的運費", false)
	idxTotalPay := v.findStringInArray(lineFirst, "訂單總金額", false)
	idxShopeeCoin := v.findStringInArray(lineFirst, "蝦幣折抵", false)
	// idxShopeeDiscount := v.findStringInArray(lineFirst, "蝦皮發放折扣券", false)
	// idxSellerDiscount := v.findStringInArray(lineFirst, "賣家自設折扣券", false)

	totalIncome := 0.0
	var orderInfo []OrderInfo
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		for i := 0; i < len(line); i++ {
			// Encodebig5(line[i], &line[i])
		}

		var oneOrder OrderInfo
		oneOrder.orderId = line[idxOrderId]
		oneOrder.orderStatus = line[idxOrderStatus]
		oneOrder.orderReturn = line[idxOrderReturn]
		oneOrder.buyerAccount = line[idxBuyerAccount]
		oneOrder.orderTime = line[idxOrderTime]
		oneOrder.payTime = line[idxPayTime]
		oneOrder.orderAmount = 0
		oneOrder.buyerFreight = 0
		oneOrder.totalPay = 0
		oneOrder.shopeeCoin = 0
		oneOrder.shopeeDiscount = 0
		oneOrder.sellerDiscount = 0
		oneOrder.buyDetail = "" //line[idxBuyDetail]
		oneOrder.tmpA = line[13]
		oneOrder.tmpB = line[14]
		oneOrder.recvAddr = line[idxRecvAddr]
		oneOrder.country = line[16]
		oneOrder.city = line[idxCity]
		oneOrder.district = line[idxDistrict]
		oneOrder.postalCode = line[idxPostalCode]
		oneOrder.buyerName = line[idxBuyerName]
		oneOrder.phone = line[idxPhone]
		oneOrder.shippingMethod = line[idxShippingMethod]
		oneOrder.shipmentMethod = line[idxShipmentMethod]
		oneOrder.orderType = "" //line[idxOrderType]
		oneOrder.payMethod = line[idxPayMethod]
		oneOrder.ccLast4 = line[idxCcLast4]
		oneOrder.lastShippingTime = line[idxLastShippingTime]
		oneOrder.trackingNum = line[idxTrackingNum]
		oneOrder.realShippingTime = line[idxRealShippingTime]
		oneOrder.orderCompleteTime = line[idxOrderCompleteTime]
		oneOrder.buyerComment = line[idxBuyerComment]
		oneOrder.comment = "" //line[idxComment]
		oneOrder.orderAmount, _ = strconv.ParseFloat(line[idxOrderAmount], 32)
		oneOrder.buyerFreight, _ = strconv.ParseFloat(line[idxBuyerFreight], 32)
		oneOrder.totalPay, _ = strconv.ParseFloat(line[idxTotalPay], 32)
		oneOrder.shopeeCoin, _ = strconv.ParseFloat(line[idxShopeeCoin], 32)
		// oneOrder.shopeeDiscount, _ = strconv.ParseFloat(line[idxShopeeDiscount], 32)
		// oneOrder.sellerDiscount, _ = strconv.ParseFloat(line[idxSellerDiscount], 32)
		orderInfo = append(orderInfo, oneOrder)
		strSQL = fmt.Sprintf("REPLACE INTO OrderInfo(Platform,OrderId,OrderStatus,OrderReturn,BuyerAccount,OrderTime,PayTime,OrderAmount,BuyerFreight,TotalPay,ShopeeCoin,ShopeeDiscount,SellerDiscount,BuyDetail,RecvAddr,Country,City,District,PostalCode,BuyerName,Phone,ShippingMethod,ShipmentMethod,OrderType,PayMethod,CcLast4,LastShippingTime,TrackingNum,RealShippingTime,OrderCompleteTime,BuyerComment,Comment) VALUES('shopee','%s','%s','%s','%s','%s','%s','%f','%f','%f','%f','%f','%f','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", oneOrder.orderId, oneOrder.orderStatus, oneOrder.orderReturn, oneOrder.buyerAccount, oneOrder.orderTime, oneOrder.payTime, oneOrder.orderAmount, oneOrder.buyerFreight, oneOrder.totalPay, oneOrder.shopeeCoin, oneOrder.shopeeDiscount, oneOrder.sellerDiscount, "暫時不填", oneOrder.recvAddr, oneOrder.country, oneOrder.city, oneOrder.district, oneOrder.postalCode, oneOrder.buyerName, oneOrder.phone, oneOrder.shippingMethod, oneOrder.shipmentMethod, oneOrder.orderType, oneOrder.payMethod, oneOrder.ccLast4, oneOrder.lastShippingTime, oneOrder.trackingNum, oneOrder.realShippingTime, oneOrder.orderCompleteTime, oneOrder.buyerComment, oneOrder.comment)
		_, errSql = v.mainApp.DbMySql.Exec(strSQL)
		if errSql != nil {
			fmt.Printf("dbMySql.Err=%s\n", errSql)
		} else {
			// fmt.Printf("Run SQL result=%q", result)
		}

		idxItemModelName := v.findStringInArray(lineFirst, "商品選項名稱", false)
		idxItemModel := v.findStringInArray(lineFirst, "商品選項貨號", false)
		idxItemId := v.findStringInArray(lineFirst, "主商品貨號", false)
		idxItemPrice := v.findStringInArray(lineFirst, "商品活動價格", false)
		if line[idxItemPrice] == "" {
			idxItemPrice = v.findStringInArray(lineFirst, "商品原價", false)
		}
		idxItemQty := v.findStringInArray(lineFirst, "數量", false)

		itemModelName := line[idxItemModelName]
		itemModel := line[idxItemModel]
		itemId := line[idxItemId]
		itemPriceF, _ := strconv.ParseFloat(line[idxItemPrice], 32)
		itemPrice := int(itemPriceF)
		itemQtyF, _ := strconv.ParseFloat(line[idxItemQty], 32)
		itemQty := int(itemQtyF)

		fmt.Println("orderId:", oneOrder.orderId, "itemModelName:", itemModelName, "itemModel:", itemModel, "itemId:", itemId, "itemPrice:", itemPrice, "itemQty:", itemQty)
		strSQL = fmt.Sprintf("REPLACE INTO OrderInfoBuyDetail(Platform,OrderId,ItemId,ItemModel,ItemModelName,ItemQty,ItemPrice) VALUES('shopee','%s','%s','%s','%s',%d,%d)", oneOrder.orderId, itemId, itemModel, itemModelName, itemQty, itemPrice)
		_, errSql = v.mainApp.DbMySql.Exec(strSQL)
		if errSql != nil {
			fmt.Printf("dbMySql.Err=%s", errSql)
		} else {
			// fmt.Printf("Run SQL result=%q", result)
		}

		totalIncome += oneOrder.orderAmount
		fmt.Printf("orderAmount=%v  totalIncome=%v", oneOrder.orderAmount, totalIncome)
		fmt.Println("---------")
	}

}

func (v *SalesImport) CsvImportFromShopee_old(csvFilePath string) {

	var strSQL string
	var errSql error

	csvFile, _ := os.Open(csvFilePath)
	csvFileBig5 := transform.NewReader(csvFile, traditionalchinese.Big5.NewDecoder()) //使用 big5 讀檔案
	reader := csv.NewReader(bufio.NewReader(csvFileBig5))
	defer csvFile.Close()

	lineFirst, errorFirst := reader.Read()
	if errorFirst == io.EOF {
		return
	} else if errorFirst != nil {
		log.Fatal(errorFirst)
	}

	idxOrderId := v.findStringInArray(lineFirst, "訂單編號", false)
	idxOrderStatus := v.findStringInArray(lineFirst, "訂單狀態", false)
	idxOrderReturn := v.findStringInArray(lineFirst, "退貨 / 退款狀態", false)
	idxBuyerAccount := v.findStringInArray(lineFirst, "買家帳號", false)
	idxOrderTime := v.findStringInArray(lineFirst, "訂單成立時間", false)
	idxPayTime := v.findStringInArray(lineFirst, "付款時間", false)
	idxBuyDetail := v.findStringInArray(lineFirst, "商品資訊", false)
	idxRecvAddr := v.findStringInArray(lineFirst, "收件地址", false)
	idxCity := v.findStringInArray(lineFirst, "城市", false)
	idxDistrict := v.findStringInArray(lineFirst, "行政區", false)
	idxPostalCode := v.findStringInArray(lineFirst, "郵遞區號", false)
	idxBuyerName := v.findStringInArray(lineFirst, "收件者姓名", false)
	idxPhone := v.findStringInArray(lineFirst, "電話", false)
	idxShippingMethod := v.findStringInArray(lineFirst, "寄送方式", false)
	idxShipmentMethod := v.findStringInArray(lineFirst, "出貨方式", false)
	idxOrderType := v.findStringInArray(lineFirst, "訂單類型", false)
	idxPayMethod := v.findStringInArray(lineFirst, "付款方式", false)
	idxCcLast4 := v.findStringInArray(lineFirst, "信用卡後四碼", false)
	idxLastShippingTime := v.findStringInArray(lineFirst, "最晚出貨日期", false)
	idxTrackingNum := v.findStringInArray(lineFirst, "包裹查詢號碼", false)
	idxRealShippingTime := v.findStringInArray(lineFirst, "實際出貨時間", false)
	idxOrderCompleteTime := v.findStringInArray(lineFirst, "訂單完成時間", false)
	idxBuyerComment := v.findStringInArray(lineFirst, "買家備註", false)
	idxComment := v.findStringInArray(lineFirst, "備註", true)

	idxOrderAmount := v.findStringInArray(lineFirst, "訂單小計", false)
	idxBuyerFreight := v.findStringInArray(lineFirst, "買家支付的運費", false)
	idxTotalPay := v.findStringInArray(lineFirst, "訂單總金額", false)
	idxShopeeCoin := v.findStringInArray(lineFirst, "蝦幣折抵", false)
	idxShopeeDiscount := v.findStringInArray(lineFirst, "蝦皮發放折扣券", false)
	idxSellerDiscount := v.findStringInArray(lineFirst, "賣家自設折扣券", false)

	totalIncome := 0.0
	var orderInfo []OrderInfo
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		for i := 0; i < len(line); i++ {
			// Encodebig5(line[i], &line[i])
		}

		// oneOrder := OrderInfo{
		// 	orderId:           line[0],
		// 	orderStatus:       line[1],
		// 	orderReturn:       line[2],
		// 	buyerAccount:      line[3],
		// 	orderTime:         line[4],
		// 	payTime:           line[5],
		// 	orderAmount:       0,
		// 	buyerFreight:      0,
		// 	totalPay:          0,
		// 	shopeeCoin:        0,
		// 	shopeeDiscount:    0,
		// 	sellerDiscount:    0,
		// 	buyDetail:         line[12],
		// 	tmpA:              line[13],
		// 	tmpB:              line[14],
		// 	recvAddr:          line[15],
		// 	country:           line[16],
		// 	city:              line[17],
		// 	district:          line[18],
		// 	postalCode:        line[19],
		// 	buyerName:         line[20],
		// 	phone:             line[21],
		// 	shippingMethod:    line[22],
		// 	shipmentMethod:    line[23],
		// 	orderType:         line[24],
		// 	payMethod:         line[25],
		// 	ccLast4:           line[26],
		// 	lastShippingTime:  line[27],
		// 	trackingNum:       line[28],
		// 	realShippingTime:  line[29],
		// 	orderCompleteTime: line[30],
		// 	buyerComment:      line[31],
		// 	comment:           line[32],
		// }
		// oneOrder.orderAmount, _ = strconv.ParseFloat(line[6], 32)
		// oneOrder.buyerFreight, _ = strconv.ParseFloat(line[7], 32)
		// oneOrder.totalPay, _ = strconv.ParseFloat(line[8], 32)
		// oneOrder.shopeeCoin, _ = strconv.ParseFloat(line[9], 32)
		// oneOrder.shopeeDiscount, _ = strconv.ParseFloat(line[10], 32)
		// oneOrder.sellerDiscount, _ = strconv.ParseFloat(line[11], 32)

		var oneOrder OrderInfo
		oneOrder.orderId = line[idxOrderId]
		oneOrder.orderStatus = line[idxOrderStatus]
		oneOrder.orderReturn = line[idxOrderReturn]
		oneOrder.buyerAccount = line[idxBuyerAccount]
		oneOrder.orderTime = line[idxOrderTime]
		oneOrder.payTime = line[idxPayTime]
		oneOrder.orderAmount = 0
		oneOrder.buyerFreight = 0
		oneOrder.totalPay = 0
		oneOrder.shopeeCoin = 0
		oneOrder.shopeeDiscount = 0
		oneOrder.sellerDiscount = 0
		oneOrder.buyDetail = "" //line[idxBuyDetail]
		oneOrder.tmpA = line[13]
		oneOrder.tmpB = line[14]
		oneOrder.recvAddr = line[idxRecvAddr]
		oneOrder.country = line[16]
		oneOrder.city = line[idxCity]
		oneOrder.district = line[idxDistrict]
		oneOrder.postalCode = line[idxPostalCode]
		oneOrder.buyerName = line[idxBuyerName]
		oneOrder.phone = line[idxPhone]
		oneOrder.shippingMethod = line[idxShippingMethod]
		oneOrder.shipmentMethod = line[idxShipmentMethod]
		oneOrder.orderType = "" //line[idxOrderType]
		oneOrder.payMethod = line[idxPayMethod]
		oneOrder.ccLast4 = line[idxCcLast4]
		oneOrder.lastShippingTime = line[idxLastShippingTime]
		oneOrder.trackingNum = line[idxTrackingNum]
		oneOrder.realShippingTime = line[idxRealShippingTime]
		oneOrder.orderCompleteTime = line[idxOrderCompleteTime]
		oneOrder.buyerComment = line[idxBuyerComment]
		oneOrder.comment = "" //line[idxComment]

		// oneOrder := OrderInfo{
		// 	orderId:           line[idxOrderId],
		// 	orderStatus:       line[idxOrderStatus],
		// 	orderReturn:       line[idxOrderReturn],
		// 	buyerAccount:      line[idxBuyerAccount],
		// 	orderTime:         line[idxOrderTime],
		// 	payTime:           line[idxPayTime],
		// 	orderAmount:       0,
		// 	buyerFreight:      0,
		// 	totalPay:          0,
		// 	shopeeCoin:        0,
		// 	shopeeDiscount:    0,
		// 	sellerDiscount:    0,
		// 	buyDetail:         line[idxBuyDetail],
		// 	tmpA:              line[13],
		// 	tmpB:              line[14],
		// 	recvAddr:          line[idxRecvAddr],
		// 	country:           line[16],
		// 	city:              line[idxCity],
		// 	district:          line[idxDistrict],
		// 	postalCode:        line[idxPostalCode],
		// 	buyerName:         line[idxBuyerName],
		// 	phone:             line[idxPhone],
		// 	shippingMethod:    line[idxShippingMethod],
		// 	shipmentMethod:    line[idxShipmentMethod],
		// 	orderType:         line[idxOrderType],
		// 	payMethod:         line[idxPayMethod],
		// 	ccLast4:           line[idxCcLast4],
		// 	lastShippingTime:  line[idxLastShippingTime],
		// 	trackingNum:       line[idxTrackingNum],
		// 	realShippingTime:  line[idxRealShippingTime],
		// 	orderCompleteTime: line[idxOrderCompleteTime],
		// 	buyerComment:      line[idxBuyerComment],
		// 	comment:           line[idxComment],
		// }
		oneOrder.orderAmount, _ = strconv.ParseFloat(line[idxOrderAmount], 32)
		oneOrder.buyerFreight, _ = strconv.ParseFloat(line[idxBuyerFreight], 32)
		oneOrder.totalPay, _ = strconv.ParseFloat(line[idxTotalPay], 32)
		oneOrder.shopeeCoin, _ = strconv.ParseFloat(line[idxShopeeCoin], 32)
		// oneOrder.shopeeDiscount, _ = strconv.ParseFloat(line[idxShopeeDiscount], 32)
		// oneOrder.sellerDiscount, _ = strconv.ParseFloat(line[idxSellerDiscount], 32)
		orderInfo = append(orderInfo, oneOrder)
		if idxBuyDetail < len(line) {
			oneOrder.buyDetail = line[idxBuyDetail]
			oneOrder.orderType = line[idxOrderType]
			oneOrder.comment = line[idxComment]
			oneOrder.shopeeDiscount, _ = strconv.ParseFloat(line[idxShopeeDiscount], 32)
			oneOrder.sellerDiscount, _ = strconv.ParseFloat(line[idxSellerDiscount], 32)
		}

		// strSQL = fmt.Sprintf("INSERT INTO ShopItemList(ShopId,UpdateTime,ItemIdList,ItemIdCnt) VALUES('%s','%s','%s', %d) ON DUPLICATE KEY UPDATE ShopId='%s',UpdateTime='%s',ItemIdList='%s',ItemIdCnt=%d", shopId, nowTimeStr, jsonString, len(listItemId), shopId, nowTimeStr, jsonString, len(listItemId))
		strSQL = fmt.Sprintf("REPLACE INTO OrderInfo(Platform,OrderId,OrderStatus,OrderReturn,BuyerAccount,OrderTime,PayTime,OrderAmount,BuyerFreight,TotalPay,ShopeeCoin,ShopeeDiscount,SellerDiscount,BuyDetail,RecvAddr,Country,City,District,PostalCode,BuyerName,Phone,ShippingMethod,ShipmentMethod,OrderType,PayMethod,CcLast4,LastShippingTime,TrackingNum,RealShippingTime,OrderCompleteTime,BuyerComment,Comment) VALUES('shopee','%s','%s','%s','%s','%s','%s','%f','%f','%f','%f','%f','%f','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", oneOrder.orderId, oneOrder.orderStatus, oneOrder.orderReturn, oneOrder.buyerAccount, oneOrder.orderTime, oneOrder.payTime, oneOrder.orderAmount, oneOrder.buyerFreight, oneOrder.totalPay, oneOrder.shopeeCoin, oneOrder.shopeeDiscount, oneOrder.sellerDiscount, "暫時不填", oneOrder.recvAddr, oneOrder.country, oneOrder.city, oneOrder.district, oneOrder.postalCode, oneOrder.buyerName, oneOrder.phone, oneOrder.shippingMethod, oneOrder.shipmentMethod, oneOrder.orderType, oneOrder.payMethod, oneOrder.ccLast4, oneOrder.lastShippingTime, oneOrder.trackingNum, oneOrder.realShippingTime, oneOrder.orderCompleteTime, oneOrder.buyerComment, oneOrder.comment)
		_, errSql = v.mainApp.DbMySql.Exec(strSQL)
		if errSql != nil {
			fmt.Printf("dbMySql.Err=%s\n", errSql)
		} else {
			// fmt.Printf("Run SQL result=%q", result)
		}

		if idxBuyDetail < len(line) {
			// fmt.Println( line[0], line[3], strings.Split(line[12], ";"))
			buyDetail := strings.Split(line[idxBuyDetail], "\n")
			for i := 0; i < len(buyDetail); i++ {
				// fmt.Println(buyDetail[i])
				// itemName, err := getStringKeywordRange(&buyDetail[i], "商品名稱:", ";")
				itemModelName, err := v.getStringKeywordRange(&buyDetail[i], "商品規格名稱:", ";")
				itemModel, err := v.getStringKeywordRange(&buyDetail[i], "商品選項貨號: ", ";")
				itemId, err := v.getStringKeywordRange(&buyDetail[i], "主商品貨號: ", ";")
				itemPriceS, err := v.getStringKeywordRange(&buyDetail[i], "價格: $ ", ";")
				itemQtyS, err := v.getStringKeywordRange(&buyDetail[i], "數量: ", ";")
				itemPrice, _ := strconv.Atoi(itemPriceS)
				itemQty, _ := strconv.Atoi(itemQtyS)

				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("orderId:", oneOrder.orderId, "itemModelName:", itemModelName, "itemModel:", itemModel, "itemId:", itemId, "itemPrice:", itemPrice, "itemQty:", itemQty)
				strSQL = fmt.Sprintf("REPLACE INTO OrderInfoBuyDetail(Platform,OrderId,ItemId,ItemModel,ItemModelName,ItemQty,ItemPrice) VALUES('shopee','%s','%s','%s','%s',%d,%d)", oneOrder.orderId, itemId, itemModel, itemModelName, itemQty, itemPrice)
				_, errSql = v.mainApp.DbMySql.Exec(strSQL)
				if errSql != nil {
					fmt.Printf("dbMySql.Err=%s", errSql)
				} else {
					// fmt.Printf("Run SQL result=%q", result)
				}
			}
		} else {
			idxItemModelName := v.findStringInArray(lineFirst, "商品規格名稱", false)
			idxItemModel := v.findStringInArray(lineFirst, "商品選項貨號", false)
			idxItemId := v.findStringInArray(lineFirst, "主商品貨號", false)
			idxItemPrice := v.findStringInArray(lineFirst, "商品活動價格", false)
			if line[idxItemPrice] == "" {
				idxItemPrice = v.findStringInArray(lineFirst, "商品原價", false)
			}
			idxItemQty := v.findStringInArray(lineFirst, "數量", false)

			itemModelName := line[idxItemModelName]
			itemModel := line[idxItemModel]
			itemId := line[idxItemId]
			itemPriceF, _ := strconv.ParseFloat(line[idxItemPrice], 32)
			itemPrice := int(itemPriceF)
			itemQtyF, _ := strconv.ParseFloat(line[idxItemQty], 32)
			itemQty := int(itemQtyF)

			fmt.Println("orderId:", oneOrder.orderId, "itemModelName:", itemModelName, "itemModel:", itemModel, "itemId:", itemId, "itemPrice:", itemPrice, "itemQty:", itemQty)
			strSQL = fmt.Sprintf("REPLACE INTO OrderInfoBuyDetail(Platform,OrderId,ItemId,ItemModel,ItemModelName,ItemQty,ItemPrice) VALUES('shopee','%s','%s','%s','%s',%d,%d)", oneOrder.orderId, itemId, itemModel, itemModelName, itemQty, itemPrice)
			_, errSql = v.mainApp.DbMySql.Exec(strSQL)
			if errSql != nil {
				fmt.Printf("dbMySql.Err=%s", errSql)
			} else {
				// fmt.Printf("Run SQL result=%q", result)
			}
		}

		totalIncome += oneOrder.orderAmount
		fmt.Printf("orderAmount=%v  totalIncome=%v", oneOrder.orderAmount, totalIncome)
		fmt.Println("---------")
	}

}
