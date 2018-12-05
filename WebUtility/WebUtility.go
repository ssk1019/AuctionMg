package WebUtility

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
)

func ReadWebPage(url string) (string, error) {
	// 參考 https://dlintw.github.io/gobyexample/public/http-client.html

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "", err
	}

	// var aaa []byte
	// body2, _ := utf8.DecodeRune(body)
	// str1, _, _ := transform.String(traditionalchinese.Big5.NewEncoder(), string(body))
	// str1, err := iconv.NewReader(body, "big5", "utf-8")
	// nR, nW, err := iconv.Convert(body, aaa[:], "big5", "utf-8")
	return string(body), nil
}

func CutString(strSource string, tagHead string, tagTail string, findCnt int, revTag bool) (string, int) {
	for i := 1; i <= findCnt; i++ {
		idxHead := strings.Index(strSource, tagHead)
		if idxHead < 0 {
			// fmt.Println("CutString: Error! not found tagHead")
			return "", -1
		}
		if i < findCnt {
			strSource = strSource[(idxHead + 1):]
			continue
		}

		idxTail := strings.Index(strSource[idxHead:], tagTail)
		if idxTail < 0 {
			// fmt.Println("CutString: Error! not found tagTail")
			return "", -2
		}
		idxTail += idxHead

		if revTag {
			return strSource[idxHead:(idxTail + len(tagTail))], (idxTail + len(tagTail))
		} else {
			return strSource[(idxHead + len(tagHead)):idxTail], idxTail
		}

	}
	return "", -3
}
