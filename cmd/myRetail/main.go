package main

import (
	"net/http"

	"aaronmeza.io/myRetail/pkg/config"

	"go.uber.org/zap"

	"aaronmeza.io/myRetail/pkg/retailDb"

	"aaronmeza.io/myRetail/pkg/api"
)

func main() {
	conf, err := config.Parse()
	zapLog, _ := zap.NewDevelopment()
	defer zapLog.Sync()
	logger := zapLog.Sugar()

	logger.Info("connecting to price database")
	db, err := retailDb.NewPriceDatabase()
	defer db.Close()
	if err != nil {
		logger.Fatalf("failed to connect to price database, #{err}")
	}
	logger.Info("successfully connected to database")

	myRetailAPI := api.NewAPI(conf, db, logger)

	productsHandler := myRetailAPI.Products()
	http.Handle("/products", productsHandler)
	http.Handle("/products/", productsHandler)

	logger.Fatal(http.ListenAndServe(":9000", nil))
}
