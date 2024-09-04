package Translator

import (
	"encoding/json"
	OpenAIConnection "ginni-ai-task/src/OpenAIConnection"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupMockServer(mockResponse string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))

	os.Setenv("OPEN_AI_HOST", server.URL)
	os.Setenv("OPEN_AI_MODEL_NAME", "gpt-test")

	return server
}

func TestTranslateOneMessage(t *testing.T) {

	mockResponse := `{
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "[{\"sentence\":\"Testing string\",\"index\":0}]",
					"refusal": null
				}
			}
		]
	}`

	var mockedTargetTranscription []TargetTranscription
	var openAiResponse OpenAIConnection.OpenAIResponse
	err := json.Unmarshal([]byte(mockResponse), &openAiResponse)
	if err != nil {
		t.Errorf("Error in unmarshling mocked response")
	}

	_ = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &mockedTargetTranscription)

	server := setupMockServer(mockResponse)
	defer server.Close()

	callTranscriptions := []CallTranscription{
		{
			Speaker:  "Seif",
			Time:     "20:00",
			Sentence: "تيست",
		},
	}

	Translate(&callTranscriptions)
	assert.Equal(t, callTranscriptions[0].Sentence, mockedTargetTranscription[0].Sentence)
}

func TestTranslateMultipleMessages(t *testing.T) {

	mockResponse := `{
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "[{\"sentence\":\"Testing string\",\"index\":0}]",
					"refusal": null
				}
			},
			{
				"index": 2,
				"message": {
					"role": "assistant",
					"content": "[{\"sentence\":\"Testing string\",\"index\":0}]",
					"refusal": null
				}
			}
		]
	}`

	var mockedTargetTranscription []TargetTranscription
	var openAiResponse OpenAIConnection.OpenAIResponse
	err := json.Unmarshal([]byte(mockResponse), &openAiResponse)
	if err != nil {
		t.Errorf("Error in unmarshling mocked response")
	}

	_ = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &mockedTargetTranscription)

	server := setupMockServer(mockResponse)
	defer server.Close()

	callTranscriptions := []CallTranscription{
		{
			Speaker:  "Seif",
			Time:     "20:00",
			Sentence: "تيست",
		},
		{
			Speaker:  "Ali",
			Time:     "20:05",
			Sentence: "Test string",
		},
		{
			Speaker:  "Omar",
			Time:     "20:10",
			Sentence: "كيف حالك",
		},
	}

	Translate(&callTranscriptions)
	for i, target := range mockedTargetTranscription {
		assert.Equal(t, target.Sentence, callTranscriptions[i].Sentence)
	}
}

func TestTranslateWithMixedLanguages(t *testing.T) {

	mockResponse := `{
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "[{\"sentence\":\"Testing string\",\"index\":0}]",
					"refusal": null
				}
			},
			{
				"index": 2,
				"message": {
					"role": "assistant",
					"content": "[{\"sentence\":\"Testing string\",\"index\":0}]",
					"refusal": null
				}
			}
		]
	}`

	var mockedTargetTranscription []TargetTranscription
	var openAiResponse OpenAIConnection.OpenAIResponse
	err := json.Unmarshal([]byte(mockResponse), &openAiResponse)
	if err != nil {
		t.Errorf("Error in unmarshling mocked response")
	}

	_ = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &mockedTargetTranscription)

	server := setupMockServer(mockResponse)
	defer server.Close()

	callTranscriptions := []CallTranscription{
		{
			Speaker:  "Seif",
			Time:     "20:00",
			Sentence: "Test تيست",
		},
		{
			Speaker:  "Ali",
			Time:     "20:05",
			Sentence: "Test string",
		},
		{
			Speaker:  "Omar",
			Time:     "20:10",
			Sentence: "كيف حالك",
		},
	}

	Translate(&callTranscriptions)
	for i, target := range mockedTargetTranscription {
		assert.Equal(t, target.Sentence, callTranscriptions[i].Sentence)
	}
}

func TestTranslateNoMessages(t *testing.T) {
	callTranscriptions := []CallTranscription{
		{
			Speaker:  "Seif",
			Time:     "20:00",
			Sentence: "Test string 1",
		},
		{
			Speaker:  "Ali",
			Time:     "20:05",
			Sentence: "Test string 2",
		},
		{
			Speaker:  "Omar",
			Time:     "20:10",
			Sentence: "Testing string 3",
		},
	}

	Translate(&callTranscriptions)
}

func TestMissingObjectInOpenAIResponse(t *testing.T) {

	mockResponse := `{
		"choices": [
		]
	}`

	var openAiResponse OpenAIConnection.OpenAIResponse
	err := json.Unmarshal([]byte(mockResponse), &openAiResponse)
	if err != nil {
		t.Errorf("Error in unmarshling mocked response")
	}

	server := setupMockServer(mockResponse)
	defer server.Close()

	callTranscriptions := []CallTranscription{
		{
			Speaker:  "Seif",
			Time:     "20:00",
			Sentence: "تيست",
		},
	}

	translateError := Translate(&callTranscriptions)
	assert.NotEqual(t, nil, translateError)
	assert.Equal(t, "OpenAI Response didn't contain choices", translateError.Error())
}
