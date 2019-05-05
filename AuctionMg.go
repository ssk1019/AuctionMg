package main

import (
	"fmt"

	"./Competitor"
	"./DbMySql"
	"./MainApp"
	"./Models"
)

// 計算出購買清單
func CaleBuyList(mainApp *MainApp.MainApp) {
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

		strSQL = fmt.Sprintf("SELECT TTT.ItemId, TTT.ItemName, TTT.PlatformItemId, TTT.ItemModel, TTT.Qty FROM ( "+
			"SELECT OrderInfoBuyDetail.ItemId AS ItemId, ProductInfo.NameCn AS ItemName, ProductInfo.PlatformItemId AS PlatformItemId, OrderInfoBuyDetail.ItemModel AS ItemModel, SUM(OrderInfoBuyDetail.ItemQty) AS Qty FROM OrderInfo, OrderInfobuyDetail, ProductInfo WHERE OrderInfo.OrderId = OrderInfoBuyDetail.OrderId AND OrderInfoBuyDetail.ItemId = ProductInfo.ItemId AND ProductInfo.PlatformItemId='%s' AND OrderInfo.OrderTime>='%s' AND OrderInfo.OrderTime<='%s' GROUP BY  orderinfobuydetail.ItemId, orderinfobuydetail.ItemModel "+
			") AS TTT ORDER BY TTT.Qty DESC", value, "2019-03-01", "2019-03-31")
		sqlRows, errSql := mainApp.DbMySql.Query(strSQL)

		defer sqlRows.Close()
		if errSql != nil {
			fmt.Printf("CaleBuyList: dbMySql.Err=%s", errSql)
			return
		}

		// fmt.Println(strSQL)

		var itemId, itemName, platformItemId, itemModel string
		var sellQty, totalNeed int
		totalNeed = 0
		isFirst := true
		for sqlRows.Next() {
			if err := sqlRows.Scan(&itemId, &itemName, &platformItemId, &itemModel, &sellQty); err != nil {
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
				if itemModel == name {
					diff := sellQty - stock
					need := diff
					if need < 0 {
						need = 0
					}
					totalNeed += need
					fmt.Printf("%-25v, %5d, %5d, %5d, %5d\n", itemModel, sellQty, stock, need, diff)
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
		salesImport.CsvImportFromShopee("./Data/Order.completed.20190401_20190420.csv")
	}

	if false {
		salesStatistics := new(Models.SalesStatistics)
		salesStatistics.Init(HanMainApp)
		salesStatistics.MonthlyStatistics("2018/01/01 00:00:00", "2018/12/31 23:59:59")
	}

	// Competitor.UpdateMyShopItemInfo(HanMainApp, "62140966") // 更新我的商品列表 ( ProductInfo )

	CaleBuyList(HanMainApp)

	// Competitor.CaleMonthlyIncome(HanMainApp, "62140966") // My
	// Competitor.CaleStockMoney(HanMainApp, "62140966") // My
	// Competitor.CaleMonthlyIncome(HanMainApp, "28876327") // 親子媽

}
