package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/toaster/pdf"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("cannot open file:", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatal("cannot stat file:", err)
	}

	pdfReader, err := pdf.NewReader(f, fi.Size())
	if err != nil {
		log.Fatal("failed to read PDF:", err)
	}

	totalPage, err := pdfReader.NumPage()
	if err != nil {
		log.Fatal("failed to count PDF pages:", err)
	}
	var output bytes.Buffer

	fmt.Println("I see", totalPage, "pages.")
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		page, err := pdfReader.Page(pageIndex)
		if err != nil {
			log.Fatal("failed to read PDF page:", err)
		}
		if page.V.IsNull() {
			continue
		}
		content, err := page.Content()
		if err != nil {
			log.Fatal("failed to read PDF page content:", err)
		}
		var lastLine, lastWordEnd float64
		var separator string
		var previousEndsWithHyphen bool
		for _, text := range content.Text {
			isNewLine := lastLine != 0 && lastLine != text.Y
			// empirical knowledge: +/- 1.1 is roughly close enough to be the same word
			// However, it could produce false results for tables with nearly no space between the contents.
			isNewWord := !isNewLine && lastWordEnd != 0 && (text.X-1.1 > lastWordEnd || lastWordEnd > text.X+1.1)
			if separator != "" {
				if isNewLine {
					isNewLine = false
				} else {
					output.WriteString(separator)
				}
				separator = ""
			}
			if text.S == "-" {
				separator = text.S
			} else {
				if isNewWord || (isNewLine && !previousEndsWithHyphen) {
					output.WriteString(" ")
				}
				lastWordEnd = text.X + text.W
				lastLine = text.Y
				output.WriteString(text.S)
				previousEndsWithHyphen = strings.HasSuffix(text.S, "-")
			}
		}
		if separator == "" && !previousEndsWithHyphen {
			output.WriteString(" ")
		}
		fmt.Println("after page", pageIndex, "the text is:", output.String())
	}
	cleaner := regexp.MustCompile("\\s+")
	log.Println("text:", cleaner.ReplaceAllString(output.String(), " "))
}
