package tripay

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type InstructionStep struct {
	Title string   `json:"title"`
	Steps []string `json:"steps"`
}

func GetInstruction(code string) ([]InstructionStep, error) {
	apiKey := os.Getenv("TRIPAY_API_KEY")
	url := "https://tripay.co.id/api-sandbox/payment/instruction?code=" + code

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Data    []InstructionStep `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, errors.New(result.Message)
	}

	return result.Data, nil
}
