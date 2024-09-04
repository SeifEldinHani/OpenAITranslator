package Translator

import (
	"encoding/json"
	"errors"
	OpenAIConnection "ginni-ai-task/src/OpenAIConnection"
	"log"
)

type CallTranscription struct {
	Speaker  string `json:"speaker" binding:"required"`
	Time     string `json:"time"    binding:"required"`
	Sentence string `json:"sentence" binding:"required"`
}
type TargetTranscription struct {
	Sentence string `json:"sentence"`
	Index    int    `json:"index"`
}

func Translate(callTranscription *[]CallTranscription) error {
	var targetTranscriptions []TargetTranscription
	for i, transcriptionObj := range *callTranscription {

		if checkIfArabic(transcriptionObj.Sentence) {
			transcription := TargetTranscription{
				Sentence: transcriptionObj.Sentence,
				Index:    i,
			}
			targetTranscriptions = append(targetTranscriptions, transcription)
		}

	}
	if len(targetTranscriptions) == 0 { // All Transcriptions are in english
		return nil
	}

	prompt, err := getTranslationPrompt(targetTranscriptions)

	if err != nil {
		log.Fatal("Translate -> getTranslationPrompt -> Error consutrcting prompt", err)
	}

	openAIResponse, err := OpenAIConnection.SendToOpenAI(prompt)
	if err != nil {
		log.Fatal("Translate -> SendToOpenAI -> Error fetching request to OpenAI")
		return err
	}
	if openAIResponse.Choices == nil || len(openAIResponse.Choices) == 0 {
		return errors.New("OpenAI Response didn't contain choices")
	}

	content := openAIResponse.Choices[0].Message.Content
	translatedTranscriptions, err := extractContent(content)
	if err != nil {
		return err
	}
	for _, translatedResp := range translatedTranscriptions {
		(*callTranscription)[translatedResp.Index].Sentence = translatedResp.Sentence
	}

	return nil
}

func checkIfArabic(sentence string) bool {
	log.Println("checkIfArabic -> ", sentence)
	for i := 0; i < len(sentence); i++ {
		if int(sentence[i]) < 0 || int(sentence[i]) > 127 {
			return true
		}
	}
	return false
}

func getTranslationPrompt(targetTranscriptions []TargetTranscription) (string, error) {
	promptObj, err := json.Marshal(targetTranscriptions)
	if err != nil {
		return "", err
	}
	prompt := "Translate the sentences in these objects to English, return objects of sentence and index only" + string(promptObj)
	return prompt, nil
}

func extractContent(content string) ([]TargetTranscription, error) {
	var translatedTranscriptions []TargetTranscription
	err := json.Unmarshal([]byte(content), &translatedTranscriptions)
	if err != nil {
		return []TargetTranscription{}, err
	}
	return translatedTranscriptions, nil
}
