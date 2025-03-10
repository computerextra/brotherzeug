package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

func main() {
	files, err := ioutil.ReadDir("files")
	if err != nil {
		log.Fatal(err)
	}

	// db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	csvFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"Artikel", "EAN", "UVP", "Preis"})

	eanRegex := regexp.MustCompile(`\b\d{13}\b`)
	priceRegex := regexp.MustCompile(`\b\d{1,3}(?:,\d{3})*,\d{2} €\b`)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pdf") {
			fmt.Println("Verarbeite Datei:", file.Name())
			content, err := extractTextFromPDF("files/" + file.Name())
			if err != nil {
				log.Println("Fehler beim Extrahieren des Textes:", err)
				continue
			}

			fmt.Println("Extrahierter Text:", content)

			eans := eanRegex.FindAllString(content, -1)
			prices := priceRegex.FindAllString(content, -1)

			fmt.Println("Gefundene EANs:", eans)
			fmt.Println("Gefundene Preise:", prices)

			for i, ean := range eans {
				if i < len(prices) {
					uvp := prices[i]
					uvp = strings.Replace(uvp, " €", "", 1)
					uvp = strings.Replace(uvp, ".", "", -1)
					uvp = strings.Replace(uvp, ",", ".", 1)

					// var artikel string
					// var preis float64
					// err := db.QueryRow("SELECT artikelnummer, preis FROM artikel WHERE ean = ?", ean).Scan(&artikel, &preis)
					// if err != nil {
					//     log.Println(err)
					//     continue
					// }

					// fmt.Println("Gefundener Artikel:", artikel, "Preis:", preis)
					writer.Write([]string{"artikel", ean, uvp, ""})
				}
			}
		}
	}
}

func extractTextFromPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var content strings.Builder
	r.SetPlainText(&content)

	return content.String(), nil
}
