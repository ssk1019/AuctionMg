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
	var totalMoney, totalProfit float64

	totalMoney = 0
	totalProfit = 0

	// 先取得 ProductInfo 列表
	mapProductInfo := make(map[string]map[string]ProductInfo)
	// SELECT EffectiveDate, ItemId, ItemModel, ItemModelName, ItemCost, ItemCostCourency, FreightCost, FreightCostCourency, FreightCostPerWeight FROM ProductCost WHERE ItemId="CC-0001-KLX" ORDER BY EffectiveDate DESC
	strSQL = fmt.Sprintf("SELECT ItemId, NameCN, ItemModel, ItemModelName, Weight, CostRMB FROM ProductInfo")
	sqlRows, errSql := mainApp.DbMySql.Query(strSQL)

	defer sqlRows.Close()
	if errSql != nil {
		fmt.Printf("CaleMyProfit: dbMySql.Err=%s", errSql)
		return
	}
	for sqlRows.Next() {
		var itemId, itemName, itemModel, itemModelName string
		var weight, costRMB float64
		if err := sqlRows.Scan(&itemId, &itemName, &itemModel, &itemModelName, &weight, &costRMB); err != nil {
			fmt.Printf("CaleMyProfit: dbMySql.Err=%s", err)
			break
		}
		if _, ok := mapProductInfo[itemId]; !ok {
			mapProductInfo[itemId] = map[string]ProductInfo{}
		}

		if itemModelName == "" {
			mapProductInfo[itemId]["default"] = ProductInfo{itemId: itemId, itemName: itemName, itemModel: itemModel, itemModelName: itemModelName, weight: weight, costRMB: costRMB}
		} else {
			mapProductInfo[itemId][itemModelName] = ProductInfo{itemId: itemId, itemName: itemName, itemModel: itemModel, itemModelName: itemModelName, weight: weight, costRMB: costRMB}
			// fmt.Printf("%v %v  %v\n", itemId, itemModel, weight)
		}
	}

	// 取得訂單列表
	strSQL = fmt.Sprintf("SELECT orderId, orderAmount, buyerFreight, payMethod FROM OrderInfo WHERE PayTime>='%s' AND PayTime<='%s'", dateStart, dateEnd)
	sqlRows1, errSql1 := mainApp.DbMySql.Query(strSQL)
	defer sqlRows1.Close()
	if errSql1 != nil {
		fmt.Printf("CaleMyProfit: dbMySql.Err=%s", errSql1)
		return
	}

	var orderId, payMethod string
	var orderAmount, buyerFreight int
	for sqlRows1.Next() {
		if err := sqlRows1.Scan(&orderId, &orderAmount, &buyerFreight, &payMethod); err != nil {
			fmt.Printf("CaleMyProfit: dbMySql.Err=%s", err)
			break
		}
		fmt.Printf("orderId %v %v %v %v =====================\n", orderId, orderAmount, buyerFreight, payMethod)

		totalMoney += float64(orderAmount)

		// 取得訂單明細
		var itemId, itemModel, itemModelName string
		var itemQty, itemPrice int
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
			if _, val1 := mapProductInfo[itemId]; !val1 {
				fmt.Errorf("ProductInfo 缺少 %s", itemId)
				return
			}

			productInfo := mapProductInfo[itemId][itemModelName]
			if _, val2 := mapProductInfo[itemId][itemModelName]; !val2 {
				productInfo = mapProductInfo[itemId]["default"]

			}
			fmt.Printf("    %v %v %v %v %v   %v %v\n", itemId, itemModel, itemModelName, itemQty, itemPrice, productInfo.weight, productInfo.itemModelName)

			orderProfit += (float64(itemPrice) - productInfo.costRMB*RATE_RMB_TWD) * float64(itemQty)
		}
		if buyerFreight == 0 {
			orderProfit -= 60
		}
		sqlRows2.Close()

		totalProfit += orderProfit
	}
	sqlRows1.Close()

	fmt.Printf("總收入:%v\n", totalMoney)
	fmt.Printf("總淨利:%v\n", totalProfit)

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
		salesImport.CsvImportFromShopee("./Data/Order.completed.20190501_20190531.csv")
	}

	if false {
		salesStatistics := new(Models.SalesStatistics)
		salesStatistics.Init(HanMainApp)
		salesStatistics.MonthlyStatistics("2018/01/01 00:00:00", "2018/12/31 23:59:59")
	}

	// Competitor.UpdateMyShopItemInfo(HanMainApp, "62140966") // 更新我的商品列表 ( ProductInfo )

	// CaleBuyList(HanMainApp, "2019-04-23", "2019-05-22") // 列出補貨清單
	CaleMyProfit(HanMainApp, "2019-05-01", "2019-05-30") // 計算區間(每月)淨利

	// Competitor.CaleMonthlyIncome(HanMainApp, "62140966") // My
	// Competitor.CaleStockMoney(HanMainApp, "62140966") // My 計算庫存商品總金額
	// Competitor.CaleMonthlyIncome(HanMainApp, "28876327") // 親子媽

}
