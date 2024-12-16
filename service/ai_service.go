package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// struct APIResponse untuk menyimpan format API dari huggingface
type APIResponse struct {
	Answer      string   `json:"answer"`
	Coordinates [][2]int `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	//validasi input
	if len(table) == 0 {
		return "", errors.New("tabel data kosong")
	}

	// Prepare request untuk Tapas API
	tapasURL := "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq"

	requestBody := map[string]interface{}{
		"inputs": map[string]interface{}{
			"table": table,
			"query": query,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("gagal memformat request: %v", err)
	}

	// Buat HTTP request
	req, err := http.NewRequest("POST", tapasURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Kirim request
	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Periksa status response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response error dengan status: %d", resp.StatusCode)
	}

	// Decode response
	var result APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("gagal decode response: %v", err)
	}

	// Extract answer dari response
	if len(result.Cells) != 0 {
		return result.Cells[0], nil
	}

	return "", errors.New("tidak ada jawaban dari model AI")
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	// Prepare request untuk Phi-3.5-mini-instruct
	phiURL := "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct"

	//Modifikasi prompt dengan instruction khusus
	instruction := "Respond with only one answer, explain shortly less than 50 words"

	// Gabungkan context dan query dengan format yang terstruktur
	prompt := fmt.Sprintf("Instruction: %s\nContext: %s\nQuestion: \"%s\"\nAnswer:", instruction, context, query)

	//setup parameter request
	requestBody := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"max_length":  2000, //membatasi panjang output
			"temperature": 0.7,
			"top_p":       0.9,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("gagal memformat request: %v", err)
	}

	// Buat HTTP request
	req, err := http.NewRequest("POST", phiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("gagal membuat request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Kirim request
	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Periksa status response
	if resp.StatusCode != http.StatusOK {
		return model.ChatResponse{}, fmt.Errorf("response error dengan status: %d", resp.StatusCode)
	}

	// Decode response
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.ChatResponse{}, fmt.Errorf("gagal decode response: %v", err)
	}

	// respon dari API berformat
	// [prompt]\n\nResponse:
	prefix := "Response:"
	text := result[0]["generated_text"].(string)[len(prompt):]
	index := strings.Index(text, prefix)
	var sliced string

	if index != -1 {
		sliced = text[index+len(prefix):]
	} else { // jika format API tidak sesuai
		sliced = text
	}

	// A.I terkadang mengulang prompt di akhir, jadi harus dipotong
	// cut off extra "Response:"
	index = strings.Index(sliced, prefix)
	if index != -1 {
		sliced = sliced[:index]
	}

	// cut off extra "Instruction:"
	prefix = "Instruction:"
	index = strings.Index(sliced, prefix)
	if index != -1 {
		sliced = sliced[:index]
	}

	// Extract generated text
	if len(result) > 0 && result[0]["generated_text"] != nil {
		return model.ChatResponse{
			GeneratedText: sliced,
		}, nil
	}

	return model.ChatResponse{}, errors.New("tidak ada jawaban dari model AI")
}
