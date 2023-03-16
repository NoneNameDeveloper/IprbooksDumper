package engine

import (
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"strconv"
)

func DumpBookData(bookId int) (string, []byte) {
	client := &http.Client{}

	link := "https://www.iprbookshop.ru/pdfstream.php?publicationId=" + strconv.Itoa(bookId) + "&part=null"

	req, err := http.NewRequest("GET", link, nil)

	if err != nil {
		log.Fatal(err)
	}

	//заголовки и куки для авторизации :TODO: autoreload auth
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, identity")
	req.Header.Set("Range", "bytes=0-2047")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "_ym_uid=167697153823982682; _ym_d=1676971538; privacy-policy=1; IPRSMARTLogin=423238ac7f8404b0ef12a8c9eaf19222%7Ce527c51102f6df1855f4f6b3be343327; _ym_isad=2; SN4f61b1c8b1bd0=g69ntv4i5js3qpj27kmcp7rqu3")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	// книга недоступна (куки устарели)
	// if len(bodyText) == 25462 {
	// 	log.Fatal("Wrong auth credentionals.")
	// }

	return (GetBookName(bookId)), (DecodeBytes(bodyText))
}

func Min(arr []int) int {
	min := arr[0]

	for _, val := range arr {
		if val < min {
			min = val
		}
	}
	return min
}

func DecodeBytes(b []byte) []byte {

	for i := 0; i < len(b); i += 2048 {
		for j := i; j < Min([]int{i + 100, len(b) - 1}); j += 2 {
			b[j], b[j+1] = b[j+1], b[j]
		}
	}
	return b
}

type Name struct {
	name string
}

func GetBookName(bookId int) string {
	link := "https://www.iprbookshop.ru/" + strconv.Itoa(bookId) + ".html"

	Name := Name{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.iprbookshop.ru", "iprbookshop.ru"),
	)

	c.OnHTML("h4.header-orange", func(e *colly.HTMLElement) {
		Name.name = e.Text
	})

	c.Visit(link)

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("OnError", err)
	})

	return Name.name
}
