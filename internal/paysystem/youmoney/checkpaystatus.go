package youmoney

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"vpngigabot/internal/models"
)

type OperationHistory struct {
	// Define the fields based on the expected response from the History API
}

type OperationDetails struct {
	// Define the fields based on the expected response from the OperationDetails API
}

// OperationHistoryRequest содержит параметры запроса истории операций
type OperationHistoryRequest struct {
	Type        string    `json:"type,omitempty"`         // Тип операций
	Label       string    `json:"label,omitempty"`        // Метка для отбора платежей
	From        time.Time `json:"from,omitempty"`         // От момента времени
	Till        time.Time `json:"till,omitempty"`         // До момента времени
	StartRecord string    `json:"start_record,omitempty"` // Номер начальной записи
	Records     int       `json:"records,omitempty"`      // Количество записей
	Details     bool      `json:"details,omitempty"`      // Показывать подробности операций
}

// OperationHistoryResponse представляет ответ сервера
type OperationHistoryResponse struct {
	Operations []Operation `json:"operations"`
}

// Operation представляет отдельную операцию
type Operation struct {
	OperationID string    `json:"operation_id"`         // Идентификатор операции
	Status      string    `json:"status"`               // Статус платежа
	Datetime    time.Time `json:"datetime"`             // Дата и время совершения операции
	Title       string    `json:"title"`                // Краткое описание операции
	PatternID   string    `json:"pattern_id,omitempty"` // Идентификатор шаблона (только для платежей)
	Direction   string    `json:"direction"`            // Направление движения средств
	Amount      float64   `json:"amount"`               // Сумма операции
	Label       string    `json:"label,omitempty"`      // Метка платежа (если указана)
	Type        string    `json:"type"`                 // Тип операции
}

func (ymc *YouMoneyClient) CheckPayStatus(pay *models.Pay) (string, error) {

	// Установите параметры запроса
	// params := url.Values{}
	// params.Add("records", "3")

	reqBody := encodeURLParams(OperationHistoryRequest{Label: pay.PayId, Records: 3})

	// Создайте запрос
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://yoomoney.ru/api/operation-history", bytes.NewBufferString(reqBody))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Authorization", "Bearer "+ymc.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполните запрос
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// Проверьте статус ответа
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("Error: %s\nBody: %s", resp.Status, bodyString)
		return "", errors.New(fmt.Sprintf("Error: %s\nBody: %s", resp.Status, bodyString))
	}

	// Прочитайте и обработайте ответ
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var historyResponse OperationHistoryResponse
	err = json.Unmarshal(bodyBytes, &historyResponse)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Вывод информации о полученных операциях
	for _, operation := range historyResponse.Operations {
		if operation.Status != "" {
			return operation.Status, nil
		}

		log.Printf("Operation ID: %s\nStatus: %s\nDatetime: %s\nTitle: %s\nAmount: %s\n\n",
			operation.OperationID, operation.Status, operation.Datetime, operation.Title, operation.Amount)
	}

	return "", errors.New("reccords not found " + pay.PayId)
}

func encodeURLParams(req OperationHistoryRequest) string {
	params := url.Values{}

	if req.Type != "" {
		params.Add("type", req.Type)
	}
	if req.Label != "" {
		params.Add("label", req.Label)
	}
	if !req.From.IsZero() {
		params.Add("from", req.From.Format(time.RFC3339))
	}
	if !req.Till.IsZero() {
		params.Add("till", req.Till.Format(time.RFC3339))
	}
	if req.StartRecord != "" {
		params.Add("start_record", req.StartRecord)
	}
	if req.Records > 0 {
		params.Add("records", strconv.Itoa(req.Records))
	}
	if req.Details {
		params.Add("details", strconv.FormatBool(req.Details))
	}

	return params.Encode()
}
