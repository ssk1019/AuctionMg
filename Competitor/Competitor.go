package Competitor

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"../MainApp"
	"../WebUtility"
)

// 取出該店家所有商品 itemId
func caleMonthlyIncome_getItemIdList(mainApp *MainApp.MainApp, shopId string) (map[int]string, error) {

	// 取出商品列表：(可取出商品 ID)
	// https://shopee.tw/api/v2/search_items/?by=pop&limit=100&match_id=62140966&newest=0&order=desc&page_type=shop

	var strSQL string
	var result sql.Result
	var sqlRows *sql.Rows
	var errSql error

	listItemId := make(map[int]string)

	nowTime := time.Now()
	condTimeStr := nowTime.Add(-1440 * time.Minute) // 一天之內的資料都有效

	// 先找資料庫的是否在有效期限內
	strSQL = fmt.Sprintf("SELECT ItemIdList, ItemIdCnt FROM ShopItemList WHERE ShopId='%s' AND UpdateTime>='%s'", shopId, condTimeStr)
	fmt.Println(strSQL)
	sqlRows, errSql = mainApp.DbMySql.Query(strSQL)
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		// fmt.Printf("Run SQL result=%q", result)
		var itemIdList string
		var itemIdCnt int
		for sqlRows.Next() {
			if err := sqlRows.Scan(&itemIdList, &itemIdCnt); err != nil {
				fmt.Printf("dbMySql.Err=%s", err)
				return nil, err
			}
			// fmt.Println(itemIdList)
			err := json.Unmarshal([]byte(itemIdList), &listItemId)
			if err != nil {
				// panic(err)
				return nil, errors.New("Can't parse json! " + err.Error())
			}
		}
		if len(listItemId) > 0 {
			return listItemId, nil
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
			return nil, errors.New("Can't read web page from " + url + "  " + err.Error())
		}
		// fmt.Println(htmlData)
		bData := []byte(htmlData)

		var dat map[string]interface{}
		err = json.Unmarshal(bData, &dat)
		if err != nil {
			// panic(err)
			return nil, errors.New("Can't parse json! " + err.Error())
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
	result, errSql = mainApp.DbMySql.Exec(strSQL)
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		fmt.Printf("Run SQL result=%q", result)
	}

	return listItemId, nil
}

// 取出該店家該商品詳細資訊
func caleMonthlyIncome_getItemDetail(mainApp *MainApp.MainApp, shopId string, itemId string) (map[string]interface{}, error) {

	// 取出某商品資訊: (銷售...等)
	// https://shopee.tw/api/v2/item/get?itemid=1149763457&shopid=62140966

	var strSQL string
	var result sql.Result
	var sqlRows *sql.Rows
	var errSql error

	var itemDetail map[string]interface{}

	nowTime := time.Now()
	condTimeStr := nowTime.Add(-1440 * time.Minute) // 一天之內的資料都有效

	// 先找資料庫的是否在有效期限內
	strSQL = fmt.Sprintf("SELECT Data FROM ShopItemDetail WHERE ShopId='%s' AND ItemId='%s' AND UpdateTime>='%s'", shopId, itemId, condTimeStr)
	fmt.Println(strSQL)
	sqlRows, errSql = mainApp.DbMySql.Query(strSQL)
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		// fmt.Printf("Run SQL result=%q", result)
		hasData := false
		var strItemDetail string
		fmt.Println("a0001")
		for sqlRows.Next() {
			fmt.Println("a0002")

			if err := sqlRows.Scan(&strItemDetail); err != nil {
				fmt.Printf("dbMySql.Err=%s", err)
				return nil, err
			}
			// fmt.Println(itemIdList)
			err := json.Unmarshal([]byte(strItemDetail), &itemDetail)
			if err != nil {
				// panic(err)
				return nil, errors.New("Can't parse json! " + err.Error())
			}
			hasData = true
		}

		if hasData {
			return itemDetail, nil
		}
	}
	fmt.Println("get from web...")
	for i := 0; i < 1; i++ {
		url := fmt.Sprintf("https://shopee.tw/api/v2/item/get?itemid=%v&shopid=%v", itemId, shopId)
		fmt.Println(url)
		htmlData, err := WebUtility.ReadWebPage(url)
		if err != nil {
			return nil, errors.New("Can't read web page from " + url + "  " + err.Error())
		}
		// fmt.Println(htmlData)
		bData := []byte(htmlData)

		err = json.Unmarshal(bData, &itemDetail)
		if err != nil {
			// panic(err)
			return nil, errors.New("Can't parse json! " + err.Error())
		}
	}

	// 新增一份清單到 DB
	nowTimeStr := nowTime.Format("2006-01-02 15:04:05")

	jsonString, err := json.Marshal(itemDetail)
	// fmt.Println(len(jsonString))
	if err != nil {

	}
	strSQL = fmt.Sprintf("INSERT INTO ShopItemDetail(ShopId,ItemId,UpdateTime,Data) VALUES('%s','%s','%s','%s') ON DUPLICATE KEY UPDATE ShopId='%s',ItemId='%s',UpdateTime='%s',Data='%s'", shopId, itemId, nowTimeStr, jsonString, shopId, itemId, nowTimeStr, jsonString)
	result, errSql = mainApp.DbMySql.Exec(strSQL)
	if errSql != nil {
		fmt.Printf("dbMySql.Err=%s", errSql)
	} else {
		fmt.Printf("Run SQL result=%q", result)
	}
	return itemDetail, nil
}

func CaleMonthlyIncome(mainApp *MainApp.MainApp, shopId string) error {

	mapItemId, err := caleMonthlyIncome_getItemIdList(mainApp, shopId)
	if err != nil {
		fmt.Println("Error!" + err.Error())
		return err
	}

	// fmt.Println(mapItemId)

	for _, value := range mapItemId {
		itemDetail, err := caleMonthlyIncome_getItemDetail(mainApp, shopId, value)
		if err != nil {
			fmt.Println("Error!" + err.Error())
			return err
		}
		fmt.Println(itemDetail)
		break
	}

	return nil
}
