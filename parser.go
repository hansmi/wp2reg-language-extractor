package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/unicode"
)

const (
	langFileMagic uint32 = 0x89455A12
)

func parseLangFile(r io.Reader) ([]string, error) {
	var header struct {
		Magic   uint32
		_       uint32
		_       uint32
		Entries uint32
	}

	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, fmt.Errorf("reading %d byte header failed: %w", binary.Size(header), err)
	}

	if header.Magic != langFileMagic {
		return nil, fmt.Errorf("magic mismatch; got %#08x, want %#08x", header.Magic, langFileMagic)
	}

	offsets := make([]uint32, header.Entries)

	if err := binary.Read(r, binary.LittleEndian, &offsets); err != nil {
		return nil, fmt.Errorf("reading offset list failed: %w", err)
	}

	textData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading text data failed: %w", err)
	}

	var decodedReader *bufio.Reader
	var buf strings.Builder

	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	entries := make([]string, 0, header.Entries)
	rawReader := bytes.NewReader(textData)

	for _, offset := range offsets {
		if _, err := rawReader.Seek(int64(offset*2), os.SEEK_SET); err != nil {
			return nil, fmt.Errorf("seeking to entry failed: %w", err)
		}

		if dr := dec.Reader(rawReader); decodedReader == nil {
			decodedReader = bufio.NewReaderSize(dr, 16)
		} else {
			decodedReader.Reset(dr)
		}

		buf.Reset()

		for {
			if r, _, err := decodedReader.ReadRune(); err != nil {
				return nil, fmt.Errorf("reading decoded rune failed (buffer=%q): %w", buf.String(), err)
			} else if r != 0 {
				buf.WriteRune(r)
			} else {
				// End of string found
				break
			}
		}

		entries = append(entries, buf.String())
	}

	if uint32(len(entries)) != header.Entries {
		return nil, fmt.Errorf("header announces %d entries, found %d", header.Entries, len(entries))
	}

	return entries, nil
}
