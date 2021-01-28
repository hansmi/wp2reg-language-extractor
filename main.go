package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func withFile(path string, fn func(*os.File) error) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}

	defer fh.Close()

	return fn(fh)
}

type fileData struct {
	name    string
	entries []string
}

func writeTable(data []fileData, writeRow func([]string) error) error {
	row := make([]string, 1, len(data)+1)
	row[0] = "Index"

	longest := 0

	for _, fd := range data {
		row = append(row, fd.name)

		if len(fd.entries) > longest {
			longest = len(fd.entries)
		}
	}

	if err := writeRow(row); err != nil {
		return err
	}

	for i := 0; i < longest; i++ {
		row = make([]string, 1, len(data)+1)
		row[0] = strconv.Itoa(i)

		for _, fd := range data {
			if i < len(fd.entries) {
				row = append(row, fd.entries[i])
			} else {
				row = append(row, "")
			}
		}

		if err := writeRow(row); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	data := []fileData{}

	for _, path := range os.Args[1:] {
		fd := fileData{
			name: filepath.Base(path),
		}

		if err := withFile(path, func(fh *os.File) error {
			var err error
			fd.entries, err = parseLangFile(fh)
			return err
		}); err != nil {
			log.Fatal(err)
		}

		data = append(data, fd)
	}

	w := csv.NewWriter(os.Stdout)

	if err := writeTable(data, w.Write); err != nil {
		log.Fatal(err)
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
