package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const baseURL = "https://yoomoney.ru/api/"

type Client struct {
	Token   string
	BaseURL string
}

func NewClient(token string) *Client {
	return &Client{
		Token:   token,
		BaseURL: baseURL,
	}
}

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

type OperationHistory struct {
	// Define the fields based on the expected response from the History API
}

type OperationDetails struct {
	// Define the fields based on the expected response from the OperationDetails API
}

func (c *Client) AccountInfo() (*AccountInfo, error) {
	method := "account-info"
	url := c.BaseURL + method

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accountInfo AccountInfo
	if err := json.NewDecoder(resp.Body).Decode(&accountInfo); err != nil {
		return nil, err
	}

	return &accountInfo, nil
}

func (c *Client) OperationHistory(params map[string]string) (*OperationHistory, error) {
	method := "operation-history"
	url := c.BaseURL + method

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var operationHistory OperationHistory
	if err := json.NewDecoder(resp.Body).Decode(&operationHistory); err != nil {
		return nil, err
	}

	return &operationHistory, nil
}

func (c *Client) OperationDetails(operationID string) (*OperationDetails, error) {
	method := "operation-details"
	url := c.BaseURL + method

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)

	q := req.URL.Query()
	q.Add("operation_id", operationID)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var operationDetails OperationDetails
	if err := json.NewDecoder(resp.Body).Decode(&operationDetails); err != nil {
		return nil, err
	}

	return &operationDetails, nil
}

// type Quickpay struct {
// 	Receiver      string
// 	QuickpayForm  string
// 	Targets       string
// 	PaymentType   string
// 	Sum           float64
// 	Formcomment   string
// 	ShortDest     string
// 	Label         string
// 	Comment       string
// 	SuccessURL    string
// 	NeedFio       bool
// 	NeedEmail     bool
// 	NeedPhone     bool
// 	NeedAddress   bool
// 	BaseURL       string
// 	RedirectedURL string
// 	Response      *http.Response
// }

type Quickpay1 struct {
	Receiver     string
	QuickpayForm string
	// Targets      string
	PaymentType string
	Sum         float64
	// 	Formcomment   string
	// 	ShortDest     string
	// 	Label         string
	// 	Comment       string
	// 	SuccessURL    string
	// 	NeedFio       bool
	// 	NeedEmail     bool
	// 	NeedPhone     bool
	// 	NeedAddress   bool

	// 	RedirectedURL string
	// 	Response      *http.Response
}

func NewQuickpay1(receiver, quickpayForm, paymentType string, sum float64) string {
	q := &Quickpay{
		Receiver:     receiver,
		QuickpayForm: quickpayForm,
		// Targets:      targets,
		PaymentType: paymentType,
		Sum:         sum,
	}

	return q.request()
}

func (q *Quickpay) request() string {
	params := url.Values{}
	params.Add("receiver", q.Receiver)
	params.Add("quickpay_form", q.QuickpayForm)
	// params.Add("targets", q.Targets)
	params.Add("paymentType", q.PaymentType)
	params.Add("sum", strconv.FormatFloat(q.Sum, 'f', 2, 64))

	baseURL := "https://yoomoney.ru/quickpay/confirm.xml?"
	baseURL += params.Encode()

	response, err := http.PostForm(baseURL, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return response.Request.URL.String()
}

type Quickpay struct {
	Receiver      string
	QuickpayForm  string
	Targets       string
	PaymentType   string
	Sum           float64
	Formcomment   *string
	ShortDest     *string
	Label         *string
	Comment       *string
	SuccessURL    *string
	NeedFio       *bool
	NeedEmail     *bool
	NeedPhone     *bool
	NeedAddress   *bool
	Response      *http.Response
	RedirectedURL string
}

func NewQuickpay(
	receiver string,
	quickpayForm string,
	targets string,
	paymentType string,
	sum float64,
	formcomment *string,
	shortDest *string,
	label *string,
	comment *string,
	successURL *string,
	needFio *bool,
	needEmail *bool,
	needPhone *bool,
	needAddress *bool,
) *Quickpay {
	return &Quickpay{
		Receiver:     receiver,
		QuickpayForm: quickpayForm,
		Targets:      targets,
		PaymentType:  paymentType,
		Sum:          sum,
		Formcomment:  formcomment,
		ShortDest:    shortDest,
		Label:        label,
		Comment:      comment,
		SuccessURL:   successURL,
		NeedFio:      needFio,
		NeedEmail:    needEmail,
		NeedPhone:    needPhone,
		NeedAddress:  needAddress,
	}
}

func (q *Quickpay) Request() error {
	baseURL, err := url.Parse("https://yoomoney.ru/quickpay/confirm.xml")
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("receiver", q.Receiver)
	params.Set("quickpay-form", q.QuickpayForm)
	params.Set("targets", q.Targets)
	params.Set("paymentType", q.PaymentType)
	params.Set("sum", strconv.FormatFloat(q.Sum, 'f', -1, 64))

	baseURL.RawQuery = params.Encode()

	resp, err := http.Post(baseURL.String(), "", nil)
	if err != nil {
		return err
	}

	q.Response = resp
	q.RedirectedURL = resp.Request.URL.String()
	fmt.Println(q.RedirectedURL)
	return nil
}
func main() {
	q := NewQuickpay(
		"4100118696367561",
		"shop",
		"Sponsor this project",
		"SB",
		100.50,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	q.Request()
	fmt.Println("Response URL:", q)
}
