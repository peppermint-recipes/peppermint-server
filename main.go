package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/controller"
	"github.com/peppermint-recipes/peppermint-server/models"
)

var dbConnection *gorm.DB

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	config := config.GetConfig()

	err := models.Init(config.DB)
	if err != nil {
		fmt.Println(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/recipe", controller.CreateRecipe).Methods("POST")
	router.HandleFunc("/recipe", controller.GetRecipes).Methods("GET")
	router.HandleFunc("/recipe/{id}", controller.GetOneRecipe).Methods("GET")
	router.HandleFunc("/recipe/{id}", controller.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/recipe/{id}", controller.DeleteRecipe).Methods("DELETE")

	log.Fatal(http.ListenAndServe(config.Web.AddressAndPort, router))
}
