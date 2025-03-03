package main

import (
	"api"
	"db"
	"fmt"
	"net/http"
	"os"
	"static"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DATABASE_URL")

	if port == "" {
		port = "8080"
	}

	var err error

	err = db.Connect(databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %q\n", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/arma1.json", api.Arma1JsonHandler)
	r.HandleFunc("/arma1.xml", api.Arma1XmlHandler)
	r.HandleFunc("/arma2.json", api.Arma2JsonHandler)
	r.HandleFunc("/arma2.xml", api.Arma2XmlHandler)
	r.HandleFunc("/arma3.json", api.Arma3JsonHandler)
	r.HandleFunc("/arma3.xml", api.Arma3XmlHandler)
	r.HandleFunc("/ofp.json", api.OfpJsonHandler)
	r.HandleFunc("/ofp.xml", api.OfpXmlHandler)
	r.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
		AssetInfo: static.AssetInfo,
		Prefix:    "static",
	}))

	var handler http.Handler
	handler = handlers.CORS()(r)
	handler = handlers.CompressHandler(handler)

	fmt.Fprintf(os.Stdout, "Server launching on port %s\n", port)

	http.ListenAndServe(":"+port, handler)
}
