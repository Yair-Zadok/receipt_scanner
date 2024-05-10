// Copyright 2024, Yair Zadok, All rights reserved.

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"io"
	"strings"
	"errors"
	"github.com/Yair-Zadok/godeeby"
)


// Request Structs
type ImageURL struct {
	URL string `json:"url"`
}

type Content struct {
	Type string    `json:"type"`
	Text string    `json:"text,omitempty"`
	Image ImageURL `json:"image_url,omitempty"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Payload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	MaxTokens int      `json:"max_tokens"`
}
/////////////////////////////////////////////////////////



// Response Structs
type responseMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type Choice struct {
    Index         int      `json:"index"`
    Message       responseMessage  `json:"message"`
    Logprobs      []string `json:"logprobs,omitempty"`
    FinishReason  string   `json:"finish_reason"`
}

type Usage struct {
    PromptTokens      int `json:"prompt_tokens"`
    CompletionTokens  int `json:"completion_tokens"`
    TotalTokens       int `json:"total_tokens"`
}

type ChatCompletion struct {
    ID                string    `json:"id"`
    Object            string    `json:"object"`
    Created           int `json:"created"`
    Model             string    `json:"model"`
    Choices           []Choice  `json:"choices"`
    Usage             Usage     `json:"usage"`
    SystemFingerprint string    `json:"system_fingerprint,omitempty"`
}
/////////////////////////////////////////////////////////


// Retrieves a string response from GPT-4 with vision
func getReceiptString(base64Image, textPrompt, apiKey string) (string, error) {

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", apiKey), }
	payload := Payload{
		Model: "gpt-4-vision-preview",
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: textPrompt,
					},
					{
						Type:  "image_url",
						Image: ImageURL{URL: fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)},
					},
				},
			},
		},
		MaxTokens: 100,
	}

	payloadBytes, err := json.Marshal(payload)
	    if err != nil { return "", err }

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
    	if err != nil { return "", err }

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
    	if err != nil { return "", err }
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
    	if err != nil { return "", err }

	var completion ChatCompletion 
	err = json.Unmarshal(responseData, &completion)
    	if err != nil { return "", err }
	
	fmt.Println(completion.Choices[0].Message.Content)
	return completion.Choices[0].Message.Content, nil
}

// Organizes openai response into a struct
func get_receipt_struct(base_encoding, api_key string) (godeeby.Receipt, error) {
    prompt := `IT IS EXTREMELY IMPORTANT THAT ALL INSTRUCTIONS BE FOLLOWED WITH PRECISION: Given the provided image, give the subtotal, total, tax, date, and tips if they are present (otherwise leave the fields as XXX) in EXACTLY the following format, IF THE TAX IS INLCUDED IN THE SUBTOTAL, SUBTRACT IT FROM THE SUBTOTAL, GIVE NO OTHER RESPONSE BUT THE FOLLOWING: subtotal: XXX, total: XXX, tax: XXX, tips: XXX, date: YEAR/MONTH/DAY`
    
    response_str, err := getReceiptString(base_encoding, prompt, api_key)  
    	if err != nil { return godeeby.Receipt{}, err }

    response_str = strings.Replace(response_str, "$", "", -1)
    response_str = strings.Replace(response_str, "xxx", "0", -1)
    response_parts := strings.Split(response_str, ", ")
    if len(response_parts) != 5 { return godeeby.Receipt{}, errors.New("incorrect number of responses in list") }
    
    for _, response_part := range response_parts {
        if len(strings.Split(response_part, ": ")) != 2 { return godeeby.Receipt{}, 
        errors.New("response part has wrong length") }
    }

    subtotal := strings.Split(response_parts[0], ": ")[1]
    total := strings.Split(response_parts[1], ": ")[1]
    tax := strings.Split(response_parts[2], ": ")[1]
    tips := strings.Split(response_parts[3], ": ")[1]
    date := strings.Split(response_parts[4], ": ")[1]

    return godeeby.Receipt{Subtotal: subtotal, Total: total, Tax: tax, 
    Tips: tips, Date: date, Account: "account_TEST", Supplier: "OTHERTEST", Encoded_image: base_encoding}, nil
}



