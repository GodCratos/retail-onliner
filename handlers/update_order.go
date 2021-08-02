package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type OnlinerStatus struct {
	Status string `json:"status"`
}

func HandlerUpdateOrder(w http.ResponseWriter, r *http.Request) {
	var bodyReq OnlinerStatus
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	mapRetail, _ := ParserJSONWithoutStruct(body)

	if mapRetail["status"].(string) == "client-confirmed" {
		bodyReq.Status = "processing"
	} else if mapRetail["status"].(string) == "assembling-complete" {
		bodyReq.Status = "confirmed"
	} else if mapRetail["status"].(string) == "send-delivery-courier" {
		bodyReq.Status = "shipping"
	} else if mapRetail["status"].(string) == "complete" {
		bodyReq.Status = "delivered"
	}
	jsonBody, _ := json.Marshal(bodyReq)
	fmt.Println(string(jsonBody))
	fmt.Println("----------------------------------------------------------------------------")
	client := &http.Client{}
	url := fmt.Sprintf("https://cart.api.onliner.by/orders/%s", mapRetail["key"].(string))
	log.Println("[SERVICES:ONLINER] Start sending request : ", url)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("[SERVICES:ONLINER] Error while creating new request. Error description : ", err)
	}
	req.Header.Add("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJ1c2VyX3R5cGUiOiJiMmJfc2hvcCIsInVzZXJfZGF0YSI6eyJpZCI6MTgyMjEsIm5hbWUiOiJcdTA0MThcdTA0MjBcdTA0MWRcdTA0MTBcdTA0MjJcdTA0MTAifSwiYXBpX3R5cGUiOiJjYXJ0IiwiYjJiX29wZXJhdG9yIjoiMTgyMjFAYXBpLWNsaWVudC5jYXJ0IiwiaWF0IjoxNjI3OTM1NTA0LCJzY29wZXMiOltdLCJleHAiOjE2Mjc5MzkxMDR9.Nv0EZIC6uCzQg89Ivt8e9s2d6dgLqvyHTRfH7kpT5x_NBI7c2jsAIzpFeGZ2Menta2C1UJcN7LqHe2pRqssSCA")
	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[SERVICES:ONLINER] Error while sending new request in onliner. Error description : ", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	log.Println("[SERVICES:ONLINER] Response from Retail : ", string(body))
	if err != nil {
		fmt.Println("[SERVICES:ONLINER] Error while parsing body request. Error description : ", err)
	}
}

//ParserJSONWithoutStruct is function for update Form to struct
func ParserJSONWithoutStruct(value []byte) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(value, &jsonMap)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jsonMap, nil
}
