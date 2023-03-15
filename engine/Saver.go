package engine

import (
	"fmt"
	"os"
)

func SaveToFile(name string, data []byte) {
	f, err := os.Create(name + ".pdf")

	if err != nil {
		fmt.Println(err)
	}

	f.Write(data)
}
