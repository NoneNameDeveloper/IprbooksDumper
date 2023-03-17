package engine

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type DumpData struct {
	Name      string
	BookBytes []byte
}

// DumpBookData получает декодированный ряд байтов книги и ее название
func DumpBookData(bookIdList []int) (resArray []DumpData) {

	for _, bookId := range bookIdList {
		resValue, err := dumpData(bookId)

		if err != nil {
			log.Fatal(err)
			continue
		}

		resArray = append(resArray, resValue)
	}

	return resArray
}

func dumpData(bookId int) (DumpData, error) {
	// создаем авторизованного клиента
	client := Auth()

	// ссылка на зашифрованный контент книги
	link := "https://www.iprbookshop.ru/pdfstream.php?publicationId=" + strconv.Itoa(bookId) + "&part=null"

	requestModel, err := http.NewRequest("GET", link, nil)

	// сайт упал или какие то другие неполадки
	if err != nil {
		return DumpData{}, errors.New("Site is down!")
	}

	// делаем запрос на сайт
	response, err := client.Do(requestModel)

	if err != nil {
		return DumpData{}, errors.New("Site is down!")
	}

	// закрываем запрос во избежание потерь ресурсов
	defer response.Body.Close()

	bodyText, err := io.ReadAll(response.Body)

	if err != nil {
		return DumpData{}, errors.New("Site is down!")
	}

	if len(bodyText) == 25462 {
		return DumpData{}, errors.New("Book doesn`t exists!")
	}

	return DumpData{Name: GetBookName(bookId), BookBytes: DecodeBytes(bodyText)}, nil
}

// Min поиск минимального значения в массиве
func Min(arr []int) int {
	min := arr[0]

	for _, val := range arr {
		if val < min {
			min = val
		}
	}
	return min
}

// Auth возвращает авторизованный клиент
func Auth() http.Client {
	authUrl := "https://www.iprbookshop.ru/95835"

	// данные для авторизации ;)
	data := url.Values{}
	data.Set("action", "login")
	data.Set("username", "mtuci")
	data.Set("password", "2xNTqGZL")
	data.Set("rememberme", "1")

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, authUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// запрос на авторизацию
	authReq, _ := client.Do(r)

	// создаем дальнейший контейнер для куки авторизации
	cookieJar, _ := cookiejar.New(nil)

	url, _ := url.Parse(authUrl)

	// устанавливаем куки авторизации
	cookieJar.SetCookies(url, authReq.Cookies())

	// помещаем их в клиент
	client = &http.Client{Jar: cookieJar}

	return *client
}

// DecodeBytes декодирует набор байтов
func DecodeBytes(b []byte) []byte {

	for i := 0; i < len(b); i += 2048 {
		for j := i; j < Min([]int{i + 100, len(b) - 1}); j += 2 {
			b[j], b[j+1] = b[j+1], b[j]
		}
	}
	return b
}

// Name контейнер для имени книги
type Name struct {
	name string
}

// Получает название книги с сайта
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
