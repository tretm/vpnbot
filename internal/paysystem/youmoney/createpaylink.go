package youmoney

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"vpngigabot/internal/models"
)

func (ymc *YouMoneyClient) CreatePayLink(summAmount float64) (*models.Pay, error) {
	const op = "youmoney.CreatePayLink"

	comment := randString(10)
	q := newQuickpay(
		ymc.account,
		"shop",
		"AC", // "SB",
		summAmount,
		comment,
		ymc.domenName+"ym/"+comment, //ymc.domenName, //
	)
	err := q.Request()
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	p := &models.Pay{PayLink: q.Response.Request.URL.String(), PayId: comment}
	// fmt.Println("Response URL:", q.Response.Request.URL)
	return p, nil

}

type Quickpay struct {
	Receiver     string `json:"receiver"`
	QuickpayForm string `json:"quickpay-form"`
	// Targets       string
	PaymentType string  `json:"paymentType"`
	Sum         float64 `json:"sum"`
	// Formcomment   string
	// ShortDest     string
	Label string `json:"lable"`
	// Comment       string
	SuccessURL string `json:"successURL"`
	// NeedFio       bool
	// NeedEmail     bool
	// NeedPhone     bool
	// NeedAddress   bool
	Response      *http.Response
	RedirectedURL string
}

func newQuickpay(
	receiver string,
	quickpayForm string,
	paymentType string,
	sum float64,
	label string,
	successURL string,
) *Quickpay {
	return &Quickpay{
		Receiver:     receiver,
		QuickpayForm: quickpayForm,
		PaymentType:  paymentType,
		Sum:          sum,
		Label:        label,
		SuccessURL:   successURL,
	}
}

func (q *Quickpay) Request() error {
	const op = "youmoney.Request"
	baseURL, err := url.Parse("https://yoomoney.ru/quickpay/confirm")
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("receiver", q.Receiver)
	params.Set("quickpay-form", q.QuickpayForm)
	// params.Set("targets", q.Targets)
	params.Set("paymentType", q.PaymentType)
	params.Set("sum", strconv.FormatFloat(q.Sum, 'f', -1, 64))
	params.Set("label", q.Label)
	params.Set("successURL", q.SuccessURL)

	baseURL.RawQuery = params.Encode()

	resp, err := http.Post(baseURL.String(), "", nil)
	if err != nil {
		return err
	}
	q.Response = resp
	q.RedirectedURL = resp.Request.URL.String()
	log.Println(op, "-------- lable:", q.Label, "PayLink:", q.RedirectedURL)
	return nil
}

func randString(length int) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
