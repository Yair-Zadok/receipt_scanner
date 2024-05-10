// Copyright 2024, Yair Zadok, All rights reserved.

package main 

import (
	"io/ioutil"
	"log"
    	"fmt"
   	"net/url"
	"net/http"
	"errors"
	"github.com/IntuitDeveloper/OAuth2-Go/cache"
	"github.com/IntuitDeveloper/OAuth2-Go/config"
)

// Gets all expense categories from QuickBooks servers
func get_categories() error {
	client := &http.Client{}

	realmId := cache.GetFromCache("realmId")
	if realmId == "" {
		log.Println("No realm ID.  QBO calls only work if the accounting scope was passed!")
		return errors.New("Realm ID empty")
	}
    sql_query := "SELECT Name FROM Account WHERE Classification='Expense'"
    sql_query = url.QueryEscape(sql_query)

    request, err := http.NewRequest("GET", config.OAuthConfig.IntuitAccountingAPIHost+"/v3/company/"+realmId+"/query?query="+sql_query, nil)
	if err != nil {fmt.Println(err)}
    
    request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("accept", "application/json")
	accessToken := cache.GetFromCache("access_token")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(request)

	if err != nil {fmt.Println(err)}
	
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {fmt.Println(err)}

    responseString := string(body)
	log.Println(responseString)

	return nil
}
