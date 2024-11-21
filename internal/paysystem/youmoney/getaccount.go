package youmoney

import (
	"encoding/json"
	"log"
	"net/http"
)

type AccountInfo struct {
	Account        string           `json:"account"`
	Balance        float64          `json:"balance"`
	Currency       string           `json:"currency"`
	AccountStatus  string           `json:"account_status"`
	AccountType    string           `json:"account_type"`
	BalanceDetails *BalanceDetails  `json:"balance_details,omitempty"`
	CardsLinked    []CardLinkedInfo `json:"cards_linked,omitempty"`
}

type BalanceDetails struct {
	Total             float64  `json:"total"`
	Available         float64  `json:"available"`
	DepositionPending *float64 `json:"deposition_pending,omitempty"`
	Blocked           *float64 `json:"blocked,omitempty"`
	Debt              *float64 `json:"debt,omitempty"`
	Hold              *float64 `json:"hold,omitempty"`
}

type CardLinkedInfo struct {
	PanFragment string  `json:"pan_fragment"`
	Type        *string `json:"type,omitempty"`
}

func (ymc *YouMoneyClient) getAccount() error {
	const op = "youmoney.getAccount"
	method := "account-info"
	url := baseURL + method
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(op, err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+ymc.token)
	// req.Header.Set("Upgrade", "h2c")
	// req.Header.Set("Connection", "Upgrade")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(op, err)
		return err
	}
	defer resp.Body.Close()

	var accountInfo AccountInfo
	if err := json.NewDecoder(resp.Body).Decode(&accountInfo); err != nil {
		log.Println(op, err)
		return err

	}
	ymc.account = accountInfo.Account
	return nil

}
