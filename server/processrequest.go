package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"gopkg.in/mgo.v2/bson"
)

func ProcessRequestIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.ProcessRequest
	c := Database.C("processrequests")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var processrequestEntryList []models.ProcessRequestBundleEntry
	for _, processrequest := range result {
		var entry models.ProcessRequestBundleEntry
		entry.Title = "ProcessRequest " + processrequest.Id
		entry.Id = processrequest.Id
		entry.Content = processrequest
		processrequestEntryList = append(processrequestEntryList, entry)
	}

	var bundle models.ProcessRequestBundle
	bundle.Type = "Bundle"
	bundle.Title = "ProcessRequest Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = processrequestEntryList

	log.Println("Setting processrequest search context")
	context.Set(r, "ProcessRequest", result)
	context.Set(r, "Resource", "ProcessRequest")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadProcessRequest(r *http.Request) (*models.ProcessRequest, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("processrequests")
	result := models.ProcessRequest{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting processrequest read context")
	context.Set(r, "ProcessRequest", result)
	context.Set(r, "Resource", "ProcessRequest")
	return &result, nil
}

func ProcessRequestShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadProcessRequest(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "ProcessRequest"))
}

func ProcessRequestCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	processrequest := &models.ProcessRequest{}
	err := decoder.Decode(processrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("processrequests")
	i := bson.NewObjectId()
	processrequest.Id = i.Hex()
	err = c.Insert(processrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting processrequest create context")
	context.Set(r, "ProcessRequest", processrequest)
	context.Set(r, "Resource", "ProcessRequest")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/ProcessRequest/"+i.Hex())
}

func ProcessRequestUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	processrequest := &models.ProcessRequest{}
	err := decoder.Decode(processrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("processrequests")
	processrequest.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, processrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting processrequest update context")
	context.Set(r, "ProcessRequest", processrequest)
	context.Set(r, "Resource", "ProcessRequest")
	context.Set(r, "Action", "update")
}

func ProcessRequestDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("processrequests")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting processrequest delete context")
	context.Set(r, "ProcessRequest", id.Hex())
	context.Set(r, "Resource", "ProcessRequest")
	context.Set(r, "Action", "delete")
}