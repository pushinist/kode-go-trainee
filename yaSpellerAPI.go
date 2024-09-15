package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type SpellError struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Row         int      `json:"row"`
	Col         int      `json:"col"`
	Len         int      `json:"len"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

func CheckSpelling(text string) ([]SpellError, error) {
	baseURL := "https://speller.yandex.net/services/spellservice.json/checkText"
	data := url.Values{}
	data.Set("text", text)

	// Отправляем GET-запрос к API
	resp, err := http.Get(baseURL + "?" + data.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Декодируем ответ в структуру
	var spellErrors []SpellError
	if err := json.NewDecoder(resp.Body).Decode(&spellErrors); err != nil {
		return nil, err
	}

	return spellErrors, nil
}

func spellHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "No text provided", http.StatusBadRequest)
		return
	}

	// Проверяем орфографию
	spellErrors, err := CheckSpelling(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем ошибки в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spellErrors)
}
