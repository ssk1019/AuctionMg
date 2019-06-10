package Competitor

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"../DbMySql"
	"../MainApp"
	"../WebUtility"
)

// 取出該店家所有商品 itemId
func GetItemIdList(mainApp *MainApp.MainApp, shopId string) (map[int]string, bool, error) {

	// 取出商品列表：(可取出商品 ID)
	// https://shopee.tw/api/v2/search_items/?by=pop&limit=100&match_id=62140966&newest=0&order=desc&page_type=shop

	var strSQL string
	// var result sql.Result
	var sqlRows *sql.Rows
	var errSql error

	listItemId := make(map[int]string)

	nowTime := time.Now()
	condTimeStr := nowTime.Add(-600 * time.Minute) // 10小時之內的資料都有效

	// 先找資料庫的是否在有效期限內
	strSQL = fmt.Sprintf("SELECT ItemIdList, ItemIdCnt FROM ShopItemList WHERE ShopId='%s' AND UpdateTime>='%s'", shopId, condTimeStr)
	fmt.Println(strSQL)
	sqlRows, errSql = mainApp.DbMySql.Query(strSQL)
	defer sqlRows.Close()
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		// fmt.Printf("Run SQL result=%q", result)
		var itemIdList string
		var itemIdCnt int
		for sqlRows.Next() {
			if err := sqlRows.Scan(&itemIdList, &itemIdCnt); err != nil {
				fmt.Printf("dbMySql.Err=%s", err)
				return nil, false, err
			}

			// fmt.Println(itemIdList)
			err := json.Unmarshal([]byte(itemIdList), &listItemId)
			if err != nil {
				// panic(err)
				return nil, false, errors.New("Can't parse json! " + err.Error())
			}
		}
		if len(listItemId) > 0 {
			return listItemId, true, nil
		}
	}

	// 取出所有該店家的所有 itemID
	cntPerPage := 100 // 每一頁取得項目數

	idx := 0
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("https://shopee.tw/api/v2/search_items/?by=pop&limit=%v&match_id=%v&newest=%v&order=desc&page_type=shop", cntPerPage, shopId, idx)
		fmt.Println(url)
		htmlData, err := WebUtility.ReadWebPage(url)
		if err != nil {
			return nil, false, errors.New("Can't read web page from " + url + "  " + err.Error())
		}
		// fmt.Println(htmlData)
		bData := []byte(htmlData)

		var dat map[string]interface{}
		err = json.Unmarshal(bData, &dat)
		if err != nil {
			// panic(err)
			return nil, false, errors.New("Can't parse json! " + err.Error())
		}
		items := dat["items"].([]interface{})
		// fmt.Println(items)
		for _, value := range items {
			vvv := value.(map[string]interface{})
			itemId := strconv.FormatFloat(vvv["itemid"].(float64), 'f', 0, 64)
			// fmt.Println("idx:", idx, "Key:", key, "itemId", itemId, "Name:", vvv["name"])
			listItemId[idx] = itemId
			idx++
		}
		if len(items) < cntPerPage {
			break
		}
		time.Sleep(1)
	}

	// 新增一份清單到 DB
	nowTimeStr := nowTime.Format("2006-01-02 15:04:05")

	jsonString, err := json.Marshal(listItemId)
	// fmt.Println(len(jsonString))
	if err != nil {

	}
	strSQL = fmt.Sprintf("INSERT INTO ShopItemList(ShopId,UpdateTime,ItemIdList,ItemIdCnt) VALUES('%s','%s','%s', %d) ON DUPLICATE KEY UPDATE ShopId='%s',UpdateTime='%s',ItemIdList='%s',ItemIdCnt=%d", shopId, nowTimeStr, jsonString, len(listItemId), shopId, nowTimeStr, jsonString, len(listItemId))
	_, errSql = mainApp.DbMySql.Exec(strSQL)
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		// fmt.Printf("Run SQL result=%q", result)
	}

	return listItemId, false, nil
}

// 取出該店家該商品詳細資訊
func GetItemDetail(mainApp *MainApp.MainApp, shopId string, itemId string) (map[string]interface{}, bool, error) {

	// 取出某商品資訊: (銷售...等)
	// https://shopee.tw/api/v2/item/get?itemid=1149763457&shopid=62140966

	var strSQL string
	// var result sql.Result
	var sqlRows *sql.Rows
	var errSql error

	var itemDetail map[string]interface{}

	nowTime := time.Now()
	condTimeStr := nowTime.Add(-600 * time.Minute) // 10小時之內的資料都有效

	// 先找資料庫的是否在有效期限內
	strSQL = fmt.Sprintf("SELECT Data FROM ShopItemDetail WHERE ShopId='%s' AND ItemId='%s' AND UpdateTime>='%s'", shopId, itemId, condTimeStr)
	// fmt.Println(strSQL)
	sqlRows, errSql = mainApp.DbMySql.Query(strSQL)
	defer sqlRows.Close()
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		// fmt.Printf("Run SQL result=%q", result)
		var strItemDetail string
		if sqlRows.Next() {
			// fmt.Println("get from Db...")
			if err := sqlRows.Scan(&strItemDetail); err != nil {
				fmt.Printf("dbMySql.Err=%s", err)
				return nil, false, err
			}

			// fmt.Println(strItemDetail)
			// fmt.Println(itemIdList)
			err := json.Unmarshal([]byte(strItemDetail), &itemDetail)
			if err != nil {
				// panic(err)
				return nil, false, errors.New("Can't parse json! " + err.Error())
			}
			return itemDetail, true, nil
		}
	}
	// fmt.Println("get from web...")

	for i := 0; i < 1; i++ {
		url := fmt.Sprintf("https://shopee.tw/api/v2/item/get?itemid=%v&shopid=%v", itemId, shopId)
		// fmt.Println(url)
		htmlData, err := WebUtility.ReadWebPage(url)
		if err != nil {
			return nil, false, errors.New("Can't read web page from " + url + "  " + err.Error())
		}
		// fmt.Println(htmlData)
		bData := []byte(htmlData)

		err = json.Unmarshal(bData, &itemDetail)
		if err != nil {
			// panic(err)
			return nil, false, errors.New("Can't parse json! " + err.Error())
		}
	}

	// 新增一份清單到 DB
	nowTimeStr := nowTime.Format("2006-01-02 15:04:05")

	jsonString, err := json.Marshal(itemDetail)
	// fmt.Println(len(jsonString))
	if err != nil {

	}
	strJsonString := string(jsonString[:len(jsonString)])
	strJsonString = DbMySql.MysqlRealEscapeString(strJsonString)
	strSQL = fmt.Sprintf("INSERT INTO ShopItemDetail(ShopId,ItemId,UpdateTime,Data) VALUES('%s','%s','%s','%s') ON DUPLICATE KEY UPDATE ShopId='%s',ItemId='%s',UpdateTime='%s',Data='%s'", shopId, itemId, nowTimeStr, strJsonString, shopId, itemId, nowTimeStr, strJsonString)
	_, errSql = mainApp.DbMySql.Exec(strSQL)
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		// fmt.Printf("Run SQL result=%q", result)
	}
	return itemDetail, false, nil
}

// 計算月營收
func CaleMonthlyIncome(mainApp *MainApp.MainApp, shopId string) error {

	rand.Seed(time.Now().UTC().UnixNano())

	totalIncome := 0.0

	mapItemId, _, err := GetItemIdList(mainApp, shopId)
	if err != nil {
		fmt.Println("Error!" + err.Error())
		return err
	}

	// fmt.Println(mapItemId)
	// itemDetail, err := caleMonthlyIncome_getItemDetail(mainApp, shopId, "1116148327")
	// fmt.Println(err, itemDetail)

	cnt := 1
	totalCnt := len(mapItemId)
	for _, value := range mapItemId {
		itemDetail, isCache, err := GetItemDetail(mainApp, shopId, value)
		if err != nil {
			fmt.Println("Error!" + err.Error())
			return err
		}
		// fmt.Println(itemDetail)

		tmpA := itemDetail["item"].(map[string]interface{})
		price := tmpA["price"].(float64) / 100000
		sold := tmpA["sold"].(float64)
		price_min := tmpA["price_min"].(float64) / 100000
		price_max := tmpA["price_max"].(float64) / 100000
		name := tmpA["name"].(string)
		// stock := tmpA["stock"].(float64)
		totalIncome += (price_min + price_max) / 2 * sold
		fmt.Printf("Get ShopId=%s ItemId=%s [%d/%d] isCache=%v Name=%s\n", shopId, value, cnt, totalCnt, isCache, name)

		fmt.Printf("price=%f  sold=%f, totalIncome=%f\n", price, sold, totalIncome)
		cnt++
		if isCache == false {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}
	}

	return nil
}

// 計算庫存金額
func CaleStockMoney(mainApp *MainApp.MainApp, shopId string) error {

	rand.Seed(time.Now().UTC().UnixNano())

	mapItemId, _, err := GetItemIdList(mainApp, shopId)
	if err != nil {
		fmt.Println("Error!" + err.Error())
		return err
	}

	// fmt.Println(mapItemId)
	// itemDetail, err := caleMonthlyIncome_getItemDetail(mainApp, shopId, "1116148327")
	// fmt.Println(err, itemDetail)

	cnt := 1
	var totalMoney float64 = 0
	// totalCnt := len(mapItemId)
	for _, value := range mapItemId {
		itemDetail, isCache, err := GetItemDetail(mainApp, shopId, value)
		if err != nil {
			fmt.Println("Error!" + err.Error())
			return err
		}
		// fmt.Println(itemDetail)

		// fmt.Println(itemDetail["name"])

		var itemTotalMoney float64 = 0
		tmpA := itemDetail["item"].(map[string]interface{})
		tmpB := tmpA["models"].([]interface{})
		for _, modelDetail := range tmpB {
			// fmt.Printf("aaaa %v  %v", idx, modelDetail)
			tmpC := modelDetail.(map[string]interface{})
			price := tmpC["price"].(float64) / 100000
			stock := tmpC["stock"].(float64)
			name := tmpC["name"].(string)
			itemTotalMoney = itemTotalMoney + (stock * price)
			fmt.Printf("name=%v, stock=%v, price=%v\n", name, stock, price)
		}
		totalMoney = totalMoney + itemTotalMoney

		fmt.Printf("Item %v %v %v\n", tmpA["name"], itemTotalMoney, totalMoney)
		cnt++
		if isCache == false {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}
	}

	return nil
}

// 取得並更新商品列表
func UpdateMyShopItemInfo(mainApp *MainApp.MainApp, shopId string) error {

	var strSQL string
	// var result sql.Result
	var errSql error

	rand.Seed(time.Now().UTC().UnixNano())

	mapItemId, _, err := GetItemIdList(mainApp, shopId)
	if err != nil {
		fmt.Println("Error!" + err.Error())
		return err
	}

	cnt := 1
	// totalCnt := len(mapItemId)
	for _, value := range mapItemId {
		itemDetail, isCache, err := GetItemDetail(mainApp, shopId, value)
		if err != nil {
			fmt.Println("Error!" + err.Error())
			return err
		}
		// fmt.Println(itemDetail)

		// fmt.Println(itemDetail["itemid"], itemDetail["name"])

		itemInfo := itemDetail["item"].(map[string]interface{})
		fmt.Printf("%d  %v  %s", cnt, int(itemInfo["itemid"].(float64)), itemInfo["name"])

		strSQL = fmt.Sprintf("UPDATE ProductInfo SET NameCN='%s' WHERE PlatformItemId='%d'", itemInfo["name"], int(itemInfo["itemid"].(float64)))
		_, errSql = mainApp.DbMySql.Exec(strSQL)
		if errSql != nil {
			fmt.Printf("dbMySql.Err=%s", errSql)
		} else {
			// fmt.Printf("Run SQL result=%q", result)
		}

		// tmpB := tmpA["models"].([]interface{})
		// for _, modelDetail := range tmpB {
		// 	// fmt.Printf("aaaa %v  %v", idx, modelDetail)
		// 	tmpC := modelDetail.(map[string]interface{})
		// 	price := tmpC["price"].(float64) / 100000
		// 	stock := tmpC["stock"].(float64)
		// 	name := tmpC["name"].(string)
		// 	// fmt.Printf("name=%v, stock=%v, price=%v\n", name, stock, price)
		// }

		cnt++
		if isCache == false {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}
	}

	return nil
}
