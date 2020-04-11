package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/theErikss0n/peppermint-server/services/recipe-service/config"
	"github.com/theErikss0n/peppermint-server/services/recipe-service/controllers"
	"github.com/theErikss0n/peppermint-server/services/recipe-service/models"
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
	router.HandleFunc("/recipe", recipeController.CreateRecipe).Methods("POST")
	router.HandleFunc("/recipe", recipeController.GetRecipes).Methods("GET")
	router.HandleFunc("/recipe/{id}", recipeController.GetOneRecipe).Methods("GET")
	router.HandleFunc("/recipe/{id}", recipeController.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/recipe/{id}", recipeController.DeleteRecipe).Methods("DELETE")

	log.Fatal(http.ListenAndServe(config.Web.AddressAndPort, router))
}
