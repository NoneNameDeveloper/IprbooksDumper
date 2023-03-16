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

	fmt.Print("Введите ID вашей книги -> ")

	text, _ := reader.ReadString('\n')

	num := strings.Replace(text, "\n", "", -1)

	res, err := strconv.Atoi(num)

	if err != nil {
		log.Fatal("Не целое число..")
	}

	name, data := engine.DumpBookData(res)

	fmt.Println("Название: " + name)

	engine.SaveToFile(name, data)

	fmt.Println("Файл записан.")
}
