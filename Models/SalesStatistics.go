package Models

import (
	"database/sql"
	"fmt"

	"../MainApp"
)

type SalesStatistics struct {
	mainApp         *MainApp.MainApp
	exchangeRMB2TWD float64
	handlingFee     float64 // 平台手續費
}

type MonthlyStatistics struct {
	totalOrder  int     // 總訂單數
	totalProfit float64 // 總毛利
	totalIncome float64 // 總營收
	// itemSale:{},         // 商品銷售統計
	totalFreeTransCnt int     // 免運訂單數量()
	totalItemSaleCnt  int     // 商品銷售總數量
	totalHandlingFee  float64 // 總手續費
	// buyerList:{},           // 紀錄買家清單, 統計不重複買家數
	// dateList:{},            // 日期清單
}

type ItemCostInfo struct {
	ItemCostCourency     string
	FreightCostCourency  string
	Weight               float64
	ItemCost             float64
	FreightCost          float64
	FreightCostPerWeight float64
}

func (v *SalesStatistics) Init(mainApp *MainApp.MainApp) {
	v.mainApp = mainApp
	v.exchangeRMB2TWD = 4.6
	v.handlingFee = 0.0149
}

func (v *SalesStatistics) MonthlyStatistics_FindCost(itemId string, itemModelName string, orderTime string) (*ItemCostInfo, bool) {
	var InfoItemCost = new(ItemCostInfo)

	// InfoItemCost.ItemCostCourency = "RMB"
	// InfoItemCost.FreightCostCourency = "RMB"
	// InfoItemCost.Weight = 0
	// InfoItemCost.ItemCost = 10
	// InfoItemCost.FreightCost = 10
	// InfoItemCost.FreightCostPerWeight = 0
	// return InfoItemCost, true

	strSQL := fmt.Sprintf("SELECT Weight, ItemCost, ItemCostCourency, FreightCost, FreightCostCourency, FreightCostPerWeight FROM ProductCost WHERE ItemId='%s' AND ItemModelName='%s' AND EffectiveDate<='%s' ORDER BY EffectiveDate DESC LIMIT 1", itemId, itemModelName, orderTime)
	sqlRows, errSql := v.mainApp.DbMySql.Query(strSQL)

	defer sqlRows.Close()
	if errSql != nil {
		fmt.Printf("SalesStatistics.MonthlyStatistics_FindCost: dbMySql.Err=%s", errSql)
		return InfoItemCost, false
	}

	if sqlRows.Next() { // 有找到相同的 ModelName
		if err := sqlRows.Scan(&InfoItemCost.Weight, &InfoItemCost.ItemCost, &InfoItemCost.ItemCostCourency, &InfoItemCost.FreightCost, &InfoItemCost.FreightCostCourency, &InfoItemCost.FreightCostPerWeight); err != nil {
			fmt.Printf("SalesStatistics.MonthlyStatistics_FindCost: dbMySql.Err=%s", err)
			return InfoItemCost, false
		}
	} else { // 沒有相同的 ModelName 則找預設 ""

		// InfoItemCost.ItemCostCourency = "RMB"
		// InfoItemCost.FreightCostCourency = "RMB"
		// InfoItemCost.Weight = 0
		// InfoItemCost.ItemCost = 10
		// InfoItemCost.FreightCost = 10
		// InfoItemCost.FreightCostPerWeight = 0
		// return InfoItemCost, true

		sqlRows.Close()

		strSQL2 := fmt.Sprintf("SELECT Weight, ItemCost, ItemCostCourency, FreightCost, FreightCostCourency, FreightCostPerWeight FROM ProductCost WHERE ItemId='%s' AND ItemModelName='' AND EffectiveDate<='%s' ORDER BY EffectiveDate DESC LIMIT 1", itemId, orderTime)
		sqlRows2, errSql2 := v.mainApp.DbMySql.Query(strSQL2)
		defer sqlRows2.Close()
		if errSql2 != nil {
			fmt.Printf("SalesStatistics.MonthlyStatistics_FindCost: dbMySql.Err=%s", errSql)
			return InfoItemCost, false
		}
		if sqlRows2.Next() { // 有找到預設
			if err := sqlRows2.Scan(&InfoItemCost.Weight, &InfoItemCost.ItemCost, &InfoItemCost.ItemCostCourency, &InfoItemCost.FreightCost, &InfoItemCost.FreightCostCourency, &InfoItemCost.FreightCostPerWeight); err != nil {
				fmt.Printf("SalesStatistics.MonthlyStatistics_FindCost: dbMySql.Err=%s", err)
				return InfoItemCost, false
			}
		}
		sqlRows2.Close()
	}

	// 如果單位是 RMB, 轉換為 TWD
	if InfoItemCost.ItemCostCourency == "RMB" {
		InfoItemCost.ItemCost = InfoItemCost.ItemCost * v.exchangeRMB2TWD
	} else if InfoItemCost.ItemCostCourency == "TWD" {
		// InfoItemCost.ItemCost
	} else {
		fmt.Printf("SalesStatistics.MonthlyStatistics_FindCost: Error Courency. ItemCostCourency=%s", InfoItemCost.ItemCostCourency)
		return InfoItemCost, false
	}

	// 如果單位是 RMB, 轉換為 TWD
	if InfoItemCost.FreightCostCourency == "RMB" {
		InfoItemCost.FreightCost = InfoItemCost.FreightCost * v.exchangeRMB2TWD
		InfoItemCost.FreightCostPerWeight = InfoItemCost.FreightCostPerWeight * v.exchangeRMB2TWD
	} else if InfoItemCost.FreightCostCourency == "TWD" {
		// InfoItemCost.FreightCost
	} else {
		fmt.Printf("SalesStatistics.MonthlyStatistics_FindCost: Error Courency. FreightCostCourency=%s", InfoItemCost.FreightCostCourency)
		return InfoItemCost, false
	}

	return InfoItemCost, true
}

// a
func (v *SalesStatistics) MonthlyStatistics(dateStart string, dateEnd string) {
	var strSQL string
	// var result sql.Result
	var sqlRows *sql.Rows
	var errSql error

	monthlyStatistics := MonthlyStatistics{
		totalOrder:  0, // 總訂單數
		totalProfit: 0, // 總毛利
		totalIncome: 0, // 總營收
		// itemSale:{},         // 商品銷售統計
		totalFreeTransCnt: 0, // 免運訂單數量()
		totalItemSaleCnt:  0, // 商品銷售總數量
		totalHandlingFee:  0, // 總手續費
		// buyerList:{},           // 紀錄買家清單, 統計不重複買家數
		// dateList:{},            // 日期清單
	}

	var orderInfo []OrderInfo
	// 先撈出所有指定時間內的訂單
	strSQL = fmt.Sprintf("SELECT OrderId, OrderReturn, BuyerAccount, OrderTime, BuyerFreight, TotalPay, ShippingMethod, PayMethod FROM OrderInfo WHERE OrderTime>='%s' AND OrderTime<='%s'", dateStart, dateEnd)
	// fmt.Println(strSQL)
	sqlRows, errSql = v.mainApp.DbMySql.Query(strSQL)
	defer sqlRows.Close()
	if errSql != nil {
		fmt.Printf("SalesStatistics.MonthlyStatistics: dbMySql.Err=%s", errSql)
	} else {
		// fmt.Printf("Run SQL result=%q", result)
		var orderId, orderReturn, buyerAccount, orderTime, shippingMethod, payMethod string
		var buyerFreight, totalPay float64
		for sqlRows.Next() {
			if err := sqlRows.Scan(&orderId, &orderReturn, &buyerAccount, &orderTime, &buyerFreight, &totalPay, &shippingMethod, &payMethod); err != nil {
				fmt.Printf("SalesStatistics.MonthlyStatistics: dbMySql.Err=%s", err)
				break
			}
			// fmt.Println(orderId, orderReturn, buyerAccount, orderTime, buyerFreight, shippingMethod, payMethod)
			orderInfo = append(orderInfo, OrderInfo{
				orderId:        orderId,
				orderReturn:    orderReturn,
				buyerAccount:   buyerAccount,
				orderTime:      orderTime,
				buyerFreight:   buyerFreight,
				totalPay:       totalPay,
				shippingMethod: shippingMethod,
				payMethod:      payMethod})

		}
	}
	sqlRows.Close()

	// 依據每一筆訂單撈出明細併統計
	var buyDetail []OrderInfoBuyDetail
	fmt.Println("Total:", len(orderInfo))
	for i := 0; i < len(orderInfo); i++ {
		orderHandlingFee := 0        // 單筆訂單平台手續費
		orderHandlingFeeFloat := 0.0 // 單筆訂單平台手續費
		creditCardHandlingFee := 0   // 單筆訂單信用卡手續費
		orderProfit := 0.0           // 單筆訂單毛利
		fmt.Println(i, "", orderInfo[i].orderId)
		strSQL = fmt.Sprintf("SELECT ItemId, ItemModel, ItemModelName, ItemQty, ItemPrice FROM OrderInfoBuyDetail WHERE OrderId='%s'", orderInfo[i].orderId)
		// fmt.Println(strSQL)
		sqlRows, errSql = v.mainApp.DbMySql.Query(strSQL)
		var oneBuyDetail OrderInfoBuyDetail
		for sqlRows.Next() {
			if err := sqlRows.Scan(&oneBuyDetail.ItemId, &oneBuyDetail.ItemModel, &oneBuyDetail.ItemModelName, &oneBuyDetail.ItemQty, &oneBuyDetail.ItemPrice); err != nil {
				fmt.Printf("SalesStatistics.MonthlyStatistics: dbMySql.Err=%s", err)
				break
			}
			//fmt.Println(ItemId, ItemModel, ItemModelName, ItemPrice, ItemQty, ItemPrice)
			buyDetail = append(buyDetail, oneBuyDetail)
		}
		sqlRows.Close()

		for j := 0; j < len(buyDetail); j++ {
			// 找出成本資料
			itemCostInfo, bResult := v.MonthlyStatistics_FindCost(buyDetail[j].ItemId, buyDetail[j].ItemModelName, orderInfo[i].orderTime)

			if bResult == false {
				fmt.Printf("SalesStatistics.MonthlyStatistics: No Product cost info. ItemId=%s, ItemModelName=%s", buyDetail[j].ItemId, buyDetail[j].ItemModelName)
				return
			}
			// 商品成本
			itemCost := 0.0
			if itemCostInfo.Weight == 0 { // 如果重量為 0, 則不使用重量計算成本
				itemCost = itemCostInfo.ItemCost + itemCostInfo.FreightCost
			} else {
				itemCost = itemCostInfo.ItemCost + itemCostInfo.Weight
			}

			orderHandlingFeeFloat += (buyDetail[j].ItemPrice * v.handlingFee) * float64(buyDetail[j].ItemQty)
			orderProfit += (buyDetail[j].ItemPrice - itemCost)
		}
		orderHandlingFee = int(orderHandlingFeeFloat + .5)
		creditCardHandlingFee = int(orderInfo[i].totalPay + .5)
		monthlyStatistics.totalProfit += (orderProfit - float64(orderHandlingFee) - float64(creditCardHandlingFee))
	}

	fmt.Println("總淨利:", monthlyStatistics.totalProfit)

}
