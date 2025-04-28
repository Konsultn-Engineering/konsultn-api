package main

import "konsultn-api/internal/db"

func main() {
	initDB, err := db.InitDB()

	if err != nil {
		return
	}

	initDB.DB()
}
