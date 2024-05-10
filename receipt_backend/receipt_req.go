// Copyright 2024, Yair Zadok, All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/IntuitDeveloper/OAuth2-Go/cache"
	"github.com/IntuitDeveloper/OAuth2-Go/config"
)

type TaxCodeRef struct {
    Value string `json:"value"`
}

type ItemRef struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

type SalesItemLineDetail struct {
    TaxCodeRef TaxCodeRef `json:"TaxCodeRef"`
    Qty        int        `json:"Qty"`
    UnitPrice  float64    `json:"UnitPrice"`
    ItemRef    ItemRef    `json:"ItemRef"`
}

type Line struct {
    Description         string              `json:"Description"`
    DetailType          string              `json:"DetailType"`
    SalesItemLineDetail SalesItemLineDetail `json:"SalesItemLineDetail"`
    LineNum             int                 `json:"LineNum"`
    Amount              float64             `json:"Amount"`
    Id                  string              `json:"Id"`
}

type SalesReceipt struct {
    Line []Line `json:"Line"`
}

// Adds a Sales Receipt to QuickBooks Servers
func post_receipt(receipt_data SalesReceipt) error {
	client := &http.Client{}

	realmId := cache.GetFromCache("realmId")
	if realmId == "" {
		log.Println("No realm ID.  QBO calls only work if the accounting scope was passed!")
        return errors.New("Realm ID empty")
	}

    json_data, err := json.Marshal(receipt_data)
    if err != nil {
        fmt.Printf("Error marshaling struct: %v\n", err)
        return err
    }

	request, err := http.NewRequest("POST", config.OAuthConfig.IntuitAccountingAPIHost+"/v3/company/"+realmId+"/salesreceipt?minorversion=70", bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("accept", "application/json")
	accessToken := cache.GetFromCache("access_token")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(request)
	if err != nil {
        fmt.Printf("Error: %v", err)
        return err
    }
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
        return err
	}

	responseString := string(body)
	log.Println(responseString)

	return nil
}
