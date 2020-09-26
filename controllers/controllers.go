package controllers

import (
	"fmt"
	"golang-api/repository"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ReturnData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	s := strings.Split(key, ",")
	limit, _ := strconv.Atoi(vars["limit"])
	data := repository.ReadData(s, limit)
	for _, d := range data {
		fmt.Fprintf(w, d)
	}

}
func AddNewData(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	if repository.InsertData(reqBody) {
		fmt.Fprintf(w, "%+v", "Data added")
	} else {
		fmt.Fprintf(w, "%+v", "Error storing data")
	}
}
