package OpenAIConnection

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func SendToOpenAI(prompt string) (OpenAIResponse, error) {

	var Messages []Message
	openAiMessage := Message{
		Role:    "user",
		Content: prompt,
	}
	Messages = append(Messages, openAiMessage)

	openAIRequest := OpenAIRequest{
		Model:    os.Getenv("OPEN_AI_MODEL_NAME"),
		Messages: Messages,
	}

	openAIResponse, err := sendOpenAIRequest(openAIRequest)
	if err != nil {
		return OpenAIResponse{}, err
	}
	return openAIResponse, err
}

func formChatCompletionRequest(jsonRequest []byte) (*http.Request, error) {

	req, err := http.NewRequest("POST", os.Getenv("OPEN_AI_HOST")+"/v1/chat/completions", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("OpenAIConnection -> formChatCompletionRequest -> Error constructing http request", err)
		return &http.Request{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPEN_AI_KEY"))

	return req, nil
}

func sendOpenAIRequest(openAIRequest OpenAIRequest) (OpenAIResponse, error) {

	client := &http.Client{}

	jsonData, err := json.Marshal(openAIRequest)
	if err != nil {
		return OpenAIResponse{}, err
	}

	req, err := formChatCompletionRequest(jsonData)
	if err != nil {
		return OpenAIResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("OpenAIConnection -> sendOpenAIRequest -> Error communicating to OpenAI", err)
		return OpenAIResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OpenAIResponse{}, err
	}

	var response OpenAIResponse

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return OpenAIResponse{}, err
	}

	return response, nil

}
