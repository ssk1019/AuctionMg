package main

import (
	"fmt"

	"./DbMySql"
	"./MainApp"
	"./Models"
)

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
	salesImport := new(Models.SalesImport)
	salesImport.Init(HanMainApp)
	salesImport.CsvImportFromShopee("./Data/fafafa1019.shopee-order.20181101-20181130.csv")
	// Competitor.CaleMonthlyIncome(HanMainApp, "62140966")	// My
	// Competitor.CaleMonthlyIncome(HanMainApp, "28876327") // 親子媽

}
