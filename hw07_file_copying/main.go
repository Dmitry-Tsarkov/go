package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" {
		slog.Error("Ошибка: не указан исходный файл")
		os.Exit(1)
	}

	if to == "" {
		slog.Error("Ошибка: не указан файл назначения")
		os.Exit(1)
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		slog.Error("Ошибка", "error", err)
		os.Exit(1)
	}

	fmt.Println("Copy succeeded")
}
