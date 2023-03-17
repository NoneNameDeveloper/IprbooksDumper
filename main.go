package main

import (
	"IprbooksDumper/engine"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите ID вашей книги, если книг несколько, введите их ID через запятую -> ")

	text, _ := reader.ReadString('\n')

	num := strings.Replace(text, "\n", "", -1)
	idList := strings.Split(num, ",")

	var idListRes []int

	// цикл и последующя валидация введенных ID
	for _, val := range idList {
		convertId, err := strconv.Atoi(strings.TrimSpace(val))

		if err != nil {
			log.Fatal("%s is not valid id.", val)
		}

		idListRes = append(idListRes, convertId)
	}

	resInfoList := engine.DumpBookData(idListRes)

	if len(resInfoList) == 0 {
		panic("All ids is wrong.")
	}

	// процессинг айдишников
	for _, dumpBook := range resInfoList {
		fmt.Println("Название: " + dumpBook.Name)
		engine.SaveToFile(dumpBook.Name, dumpBook.BookBytes)
		fmt.Println("Файл записан.\n")
	}

}
