package main

import (
	"fmt"

	"./Competitor"
	"./DbMySql"
	"./MainApp"
	"./Models"
)

// 計算出購買清單
func CaleBuyList(mainApp *MainApp.MainApp, dateStart string, dateEnd string) {
	shopId := "62140966"

	var strSQL string
	// var result sql.Result
	// var errSql error

	mapItemId, _, err := Competitor.GetItemIdList(mainApp, shopId)
	if err != nil {
		fmt.Println("Error!" + err.Error())
		return
	}

	for _, value := range mapItemId {
		itemDetail, isCache, err := Competitor.GetItemDetail(mainApp, shopId, value)
		if err != nil {
			fmt.Println("Error!" + err.Error())
			return
		}
		// fmt.Println(itemDetail, isCache)
		tmpA := itemDetail["item"].(map[string]interface{})
		tmpB := tmpA["models"].([]interface{})
		// for _, modelDetail := range tmpB {
		// 	// fmt.Printf("aaaa %v  %v", idx, modelDetail)
		// 	tmpC := modelDetail.(map[string]interface{})
		// 	// price := tmpC["price"].(float64) / 100000
		// 	stock := tmpC["stock"].(float64)
		// 	name := tmpC["name"].(string)
		// 	// fmt.Printf("name=%v, stock=%v, price=%v\n", name, stock, price)
		// }

		// strSQL = fmt.Sprintf("SELECT TTT.ItemId, TTT.ItemName, TTT.PlatformItemId, TTT.ItemModel, TTT.ItemModelName, TTT.Qty FROM ( "+
		// 	"SELECT OrderInfoBuyDetail.ItemId AS ItemId, ProductInfo.NameCn AS ItemName, ProductInfo.PlatformItemId AS PlatformItemId, OrderInfoBuyDetail.ItemModel AS ItemModel, OrderInfoBuyDetail.ItemModelName AS ItemModelName, SUM(OrderInfoBuyDetail.ItemQty) AS Qty FROM OrderInfo, OrderInfobuyDetail, ProductInfo WHERE OrderInfo.OrderId = OrderInfoBuyDetail.OrderId AND ProductInfo.ItemModelName=OrderInfoBuyDetail.ItemModelName AND OrderInfoBuyDetail.ItemId = ProductInfo.ItemId AND ProductInfo.PlatformItemId='%s' AND OrderInfo.OrderTime>='%s' AND OrderInfo.OrderTime<='%s' GROUP BY  orderinfobuydetail.ItemId, orderinfobuydetail.ItemModel "+
		// 	") AS TTT ORDER BY TTT.Qty DESC", value, dateStart, dateEnd)

		strSQL = fmt.Sprintf("SELECT TTT.ItemId, TTT.ItemName, TTT.ItemModel, TTT.ItemModelName, TTT.Qty FROM ( "+
			"SELECT OrderInfoBuyDetail.ItemId AS ItemId, ProductInfo.NameCn AS ItemName, OrderInfoBuyDetail.ItemModel AS ItemModel, OrderInfoBuyDetail.ItemModelName AS ItemModelName, SUM(OrderInfoBuyDetail.ItemQty) AS Qty FROM OrderInfo, OrderInfobuyDetail,"+
			"(SELECT NameCn, ItemId FROM ProductInfo WHERE ProductInfo.PlatformItemId='%s' GROUP BY PlatformItemId ) AS ProductInfo"+
			" WHERE ProductInfo.ItemId=OrderInfoBuyDetail.ItemId AND OrderInfo.OrderId = OrderInfoBuyDetail.OrderId AND OrderInfo.OrderTime>='%s' AND OrderInfo.OrderTime<='%s' GROUP BY  orderinfobuydetail.ItemId, orderinfobuydetail.ItemModel"+
			" ) AS TTT ORDER BY TTT.Qty DESC", value, dateStart, dateEnd)
		sqlRows, errSql := mainApp.DbMySql.Query(strSQL)

		defer sqlRows.Close()
		if errSql != nil {
			fmt.Printf("CaleBuyList: dbMySql.Err=%s", errSql)
			return
		}

		// fmt.Println(strSQL)

		var itemId, itemName, itemModel, itemModelName string
		var sellQty, totalNeed int
		totalNeed = 0
		isFirst := true
		for sqlRows.Next() {
			if err := sqlRows.Scan(&itemId, &itemName, &itemModel, &itemModelName, &sellQty); err != nil {
				fmt.Printf("CaleBuyList: dbMySql.Err=%s", err)
				break
			}
			if isFirst {
				fmt.Printf("\nItemID: %v ------- %s\n", itemId, itemName)
				fmt.Printf("          ItemModel            Sell  stock  need  need(-)\n")
				isFirst = false
			}
			for _, modelDetail := range tmpB {
				tmpC := modelDetail.(map[string]interface{})
				stock := int(tmpC["stock"].(float64))
				name := tmpC["name"].(string)
				if itemModelName == name {
					diff := sellQty - stock
					need := diff
					if need < 0 {
						need = 0
					}
					totalNeed += need
					fmt.Printf("%-25v, %5d, %5d, %5d, %5d\n", itemModelName, sellQty, stock, need, diff)
				}
			}

			// fmt.Printf("%v  %v  %d", itemId, itemModel, sellQty)
		}
		fmt.Printf("      Total Need = %d\n", totalNeed)
		fmt.Printf("\n")

		_ = itemDetail
		_ = isCache
	}
}

// 計算區間營收淨利
type ProductInfo struct {
	itemId        string
	itemName      string
	itemModel     string
	itemModelName string
	weight        float64
	costRMB       float64
}

func CaleMyProfit(mainApp *MainApp.MainApp, dateStart string, dateEnd string) {
	RATE_RMB_TWD := float64(4.6)

	var strSQL string
	var orderCnt, totalItemSellCnt int
	var totalMoney, totalProfit, totalPlatformFee, totalCreditCardFee float64

	orderCnt = 0
	totalItemSellCnt = 0
	totalMoney = 0
	totalProfit = 0
	totalPlatformFee = 0
	totalCreditCardFee = 0

	// 取得訂單列表
	strSQL = fmt.Sprintf("SELECT orderId, orderAmount, totalPay, buyerFreight, payMethod FROM OrderInfo WHERE PayTime>='%s' AND PayTime<='%s'", dateStart, dateEnd)
	sqlRows1, errSql1 := mainApp.DbMySql.Query(strSQL)
	defer sqlRows1.Close()
	if errSql1 != nil {
		fmt.Printf("CaleMyProfit: dbMySql.Err=%s", errSql1)
		return
	}

	var orderId, payMethod string
	var orderAmount, totalPay, buyerFreight int
	var orderFreightRemit, orderPlatformFee, orderCreditCardFee float64
	for sqlRows1.Next() {
		if err := sqlRows1.Scan(&orderId, &orderAmount, &totalPay, &buyerFreight, &payMethod); err != nil {
			fmt.Printf("CaleMyProfit: dbMySql.Err=%s", err)
			break
		}
		fmt.Printf("orderId %v %v %v %v =====================\n", orderId, orderAmount, buyerFreight, payMethod)

		orderCnt += 1
		totalMoney += float64(orderAmount)

		orderPlatformFee = 0
		// 取得訂單明細
		var itemId, itemModel, itemModelName string
		var itemQty, itemPrice int
		var itemCost = float64(0)
		var itemFreightCost = float64(0)
		var orderProfit = float64(0)
		strSQL = fmt.Sprintf("SELECT itemId, itemModel, itemModelName, itemQty, itemPrice FROM OrderInfoBuyDetail WHERE OrderId='%s'", orderId)
		sqlRows2, errSql2 := mainApp.DbMySql.Query(strSQL)
		defer sqlRows2.Close()
		if errSql2 != nil {
			fmt.Printf("CaleMyProfit: dbMySql.Err=%s", errSql2)
			return
		}
		for sqlRows2.Next() {
			if err := sqlRows2.Scan(&itemId, &itemModel, &itemModelName, &itemQty, &itemPrice); err != nil {
				fmt.Printf("CaleMyProfit: dbMySql.Err=%s", err)
				break
			}

			if orderId == "19052923376HRSQ" {
				fmt.Println()
			}

			// 取得價格清單
			var findCost = false
			var findFreightCost = false
			var itemPlatformFee float64
			var tmpItemModel, tmpItemModelName, tmpItemCostCourency, tmpFreightCostCourency string
			var tmpItemCost, tmpItemWeight, tmpFreightCost, tmpFreightCostPerWeight float64
			strSQL = fmt.Sprintf("SELECT ItemModel, ItemModelName, ItemCost, ItemCostCourency, Weight, FreightCost, FreightCostCourency, FreightCostPerWeight FROM ProductCost WHERE ItemId='%s' ORDER BY EffectiveDate DESC", itemId)
			sqlRows3, errSql3 := mainApp.DbMySql.Query(strSQL)
			defer sqlRows3.Close()
			if errSql3 != nil {
				fmt.Printf("CaleMyProfit: dbMySql.Err=%s", errSql2)
				return
			}
			for sqlRows3.Next() {
				if err := sqlRows3.Scan(&tmpItemModel, &tmpItemModelName, &tmpItemCost, &tmpItemCostCourency, &tmpItemWeight, &tmpFreightCost, &tmpFreightCostCourency, &tmpFreightCostPerWeight); err != nil {
					fmt.Printf("CaleMyProfit: dbMySql.Err=%s", err)
					break
				}
				if tmpItemModel == "" || tmpItemModel == itemModel {
					// 計算商品成本
					if tmpItemCostCourency == "RMB" {
						itemCost = tmpItemCost * RATE_RMB_TWD
						findCost = true
					} else if tmpItemCostCourency == "TWD" {
						itemCost = tmpItemCost
						findCost = true
					}

					// 計算商品運費成本
					if tmpItemWeight == 0.0 { // 直接價格計算
						if tmpFreightCostCourency == "RMB" {
							itemFreightCost = tmpFreightCost * RATE_RMB_TWD
							findFreightCost = true
						} else if tmpFreightCostCourency == "TWD" {
							itemFreightCost = tmpFreightCost
							findFreightCost = true
						}
					} else { // 以重量計算
						if tmpFreightCostCourency == "RMB" {
							itemFreightCost = tmpFreightCostPerWeight * tmpItemWeight / 1000 * RATE_RMB_TWD
							findFreightCost = true
						} else if tmpFreightCostCourency == "TWD" {
							itemFreightCost = tmpFreightCostPerWeight * tmpItemWeight / 1000
							findFreightCost = true
						}
					}
					break
				}

			}
			sqlRows3.Close()
			if findCost == false || findFreightCost == false {
				fmt.Printf("CaleMyProfit: Couldn't find item cost! itemId=%s, itemCost=%f, itemFreightCost=%f", itemId, itemCost, itemFreightCost)
				return
			}

			totalItemSellCnt += itemQty
			itemPlatformFee = (float64(itemPrice) * 0.0149) * float64(itemQty)
			orderPlatformFee += itemPlatformFee

			allItemProfit := (float64(itemPrice) - itemCost - itemFreightCost) * float64(itemQty)
			orderProfit += allItemProfit

			fmt.Printf("    %v %v %v qty=%v price=%v cost=%v freightCost=%v 平台手續費=%v, 商品總淨利=%v\n", itemId, itemModel, itemModelName, itemQty, itemPrice, itemCost, itemFreightCost, itemPlatformFee, allItemProfit)
		}
		if buyerFreight == 0 {
			orderFreightRemit = 60
		} else {
			orderFreightRemit = 0
		}

		if payMethod == "信用卡" {
			orderCreditCardFee = float64(totalPay) * 0.015
		} else {
			orderCreditCardFee = 0
		}

		sqlRows2.Close()

		totalPlatformFee += orderPlatformFee
		totalCreditCardFee += orderCreditCardFee
		totalProfit += (orderProfit - orderFreightRemit - orderPlatformFee - orderCreditCardFee)

		fmt.Printf("    減免運費=%v, 平台交易費=%v, 信用卡費=%v, 訂單淨利=%v, 累積淨利=%v, 累積營收=%v\n\n", orderFreightRemit, orderPlatformFee, orderCreditCardFee, orderProfit, totalProfit, totalMoney)
	}
	sqlRows1.Close()

	fmt.Printf("訂單筆數:%v\n", orderCnt)
	fmt.Printf("銷售商品總數:%v  平均訂單購買商品數:%v\n", totalItemSellCnt, float32(totalItemSellCnt)/float32(orderCnt))
	fmt.Printf("總收入:%v\n", totalMoney)
	fmt.Printf("總信用卡手續費:%v\n", totalCreditCardFee)
	fmt.Printf("總平台手續費:%v\n", totalPlatformFee)
	fmt.Printf("總淨利:%v (%v)\n", totalProfit, (totalProfit * 100 / float64(totalMoney)))
}

func main() {
	// testWespai()

	// getStr, idxFindStrEnd := WebUtility.CutString("aa<aa>bb<aa>ccddeeffgghhjiijj</aa>kk</aa>", "<aa>", "</aa>", true)
	// fmt.Println(getStr, idxFindStrEnd)
	HanMainApp := new(MainApp.MainApp)

	HanMainApp.DbMySql = new(DbMySql.DbMySql)
	err := HanMainApp.DbMySql.Create("127.0.0.1", "3306", "AuctionMgUser", "auctionMgUser123", "AuctionMg")
	if err != nil {
		fmt.Printf("Fail to create MySQL db...%s", err)
	}

	// strSQL := fmt.Sprintf("INSERT INTO ShopItemList(ShopId,UpdateTime,ItemIdList,ItemIdCnt) VALUES('%s','%s','%s', %d) ON DUPLICATE KEY UPDATE ShopId='%s',UpdateTime='%s',ItemIdList='%s',ItemIdCnt=%d", "aaa", "2018-01-01 20:30:00", "{}", 1, "aaa", "2018-01-01 20:30:00", "{}", 1)
	// result, errSql := HanMainApp.DbMySql.Exec(strSQL)
	// fmt.Println(result, errSql)
	if false {
		salesImport := new(Models.SalesImport)
		salesImport.Init(HanMainApp)
		// salesImport.CsvImportFromShopee("./Data/fafafa1019.shopee-order.20181101-20181130.csv")
		// salesImport.CsvImportFromShopee("./Data/Order.completed.20190201_20190228.csv")
		// salesImport.CsvImportFromShopee("./Data/Order.completed.20190301_20190331.csv")
		// salesImport.CsvImportFromShopee("./Data/Order.completed.20190101_20190131.csv")
		salesImport.CsvImportFromShopee("./Data/Order.completed.20190601_20190630.csv")
	}

	if false {
		salesStatistics := new(Models.SalesStatistics)
		salesStatistics.Init(HanMainApp)
		salesStatistics.MonthlyStatistics("2019/06/01 00:00:00", "2019/06/30 23:59:59")
	}

	// Competitor.UpdateMyShopItemInfo(HanMainApp, "62140966") // 更新我的商品列表 ( ProductInfo )

	// CaleBuyList(HanMainApp, "2019-04-23", "2019-05-22") // 列出補貨清單
	// CaleMyProfit(HanMainApp, "2019-06-01", "2019-06-30") // 計算區間(每月)淨利

	Competitor.CaleMonthlyIncome(HanMainApp, "62140966") // My
	// Competitor.CaleStockMoney(HanMainApp, "62140966") // My 計算庫存商品總金額
	// Competitor.CaleMonthlyIncome(HanMainApp, "28876327") // 親子媽

}
