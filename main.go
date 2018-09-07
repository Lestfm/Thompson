package main

import (
	"flag"
	"./tompson"
	"./db"
	"fmt"
)

func main() {
	dbPath := flag.String("db", "/tompson", "Путь до файла базы данных")
	port := flag.Int("port", 8081, "Порт который слушаем")
	db, err := db.NewDb(*dbPath)
	if err != nil {
		panic("Невозможно создать/открыть файл БД")
	}
	storage := tompson.NewStorage(db)
	router := tompson.NewRouter(storage)
	router.ListenAndServe(fmt.Sprintf(":%d", *port))

}

