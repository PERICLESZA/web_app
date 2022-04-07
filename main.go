package main

import (
	"net/http"

	"./utils"

	"./models"
	"./routes"
)

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
