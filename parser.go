package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)
import "golang.org/x/net/html/charset"

type Valute struct {
	NumCode  string
	CharCode string
	Nominal  int
	Name     string
	Value    float32
}

type Result struct {
	ValCurs []Valute `xml:"Valute"`
}


func main() {
	cmdargcurrency := flag.String("currency", "USD", "currency")
	cmdargvalue := flag.Int("value", 500, "value")
	flag.Parse()

	ratesURL := "http://www.cbr.ru/scripts/XML_daily.asp"

	resp, err := http.Get(ratesURL)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	bodyString := strings.Replace(string(body), ",", ".", -1)

	v := Result{}
	d := xml.NewDecoder(strings.NewReader(bodyString))
	d.CharsetReader = charset.NewReaderLabel
	err = d.Decode(&v)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	for _, curr := range v.ValCurs {
		if curr.CharCode == *cmdargcurrency {
			fmt.Printf("%.2f RUB\n", curr.Value*float32(*cmdargvalue)/float32(curr.Nominal))
			os.Exit(0)
		}
	}

	fmt.Println("currency not found")
}
