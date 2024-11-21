package vtoray

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"vpngigabot/internal/config"
	"vpngigabot/internal/models"

	"github.com/google/uuid"
)

type Vray struct {
	cfg     *config.Config
	session string
}

type user struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Obj     struct {
		ID         int    `json:"id"`
		InboundID  int    `json:"inboundId"`
		Enable     bool   `json:"enable"`
		Email      string `json:"email"`
		Up         int    `json:"up"`
		Down       int    `json:"down"`
		ExpiryTime int    `json:"expiryTime"`
		Total      int    `json:"total"`
		Reset      int    `json:"reset"`
	} `json:"obj"`
}

type clientsJSON []struct {
	ID         string `json:"id"`
	Flow       string `json:"flow"`
	Email      string `json:"email"`
	LimitIP    int    `json:"limitIp"`
	TotalGB    int    `json:"totalGB"`
	ExpiryTime int64  `json:"expiryTime"`
	Enable     bool   `json:"enable"`
	TgID       string `json:"tgId"`
	SubID      string `json:"subId"`
	Reset      int    `json:"reset"`
}

type settingsJSON struct {
	Clients clientsJSON `json:"clients"`
}

func New(cfg *config.Config) *Vray {
	return login(cfg)
}

func login(cfg *config.Config) *Vray {
	const op = "vray.login"

	apiUrl := cfg.UrlVray + "login"
	// Параметры, которые нужно URL-кодировать
	data := url.Values{}
	data.Set("username", cfg.VrayLogin)
	data.Set("password", cfg.VrayPassword)

	// Создание HTTP-клиента
	client := &http.Client{}

	// Создание POST-запроса с URL-кодированными параметрами
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(op, err)
		panic(0)
	}

	// Установка заголовков
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(data.Encode())))

	// Отправка запроса
	resp, err := client.Do(req)
	if err != nil {
		log.Println(op, err)
		panic(0)
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == "session" {
			log.Println(op, "Vpn auth let`s ok")
			return &Vray{session: c.Value, cfg: cfg}

		}
	}
	panic(0)
}
func (v *Vray) CreateKey(usr *models.User) (*models.Link, error) {
	const op = "vtoray.CreateKey"

	urlW := v.cfg.UrlVray + "panel/api/inbounds/addClient?id=1&settings="

	var link models.Link
	link.UserId = usr.UserId
	link.VpnLinkId = uuid.New().String()
	link.Link = strconv.FormatInt(link.UserId, 10) + randString(v.cfg.BotConfig.LinkLen)

	settings := &settingsJSON{
		Clients: clientsJSON{
			{
				ID:      link.VpnLinkId,
				Flow:    "xtls-rprx-vision",
				Email:   link.Link,
				LimitIP: v.cfg.Vray.LimitConnection,
				TotalGB: 0,
				// ExpiryTime: time,
				Enable: true,
				TgID:   "",
				SubID:  "",
				Reset:  0,
			},
		},
	}
	err := v.requestDo(urlW, settings)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return &link, nil
}

func (v *Vray) KeyBuilder(link *models.Link) string {
	return "vless://" + link.VpnLinkId + "@serv1..:?type=tcp&security=reality&pbk=&fp=firefox&sni=google.com&sid=14ea9969&spx=%2F&flow=xtls-rprx-vision#vpn-vless-" + link.Link

}

func (v *Vray) UpdateKey(link *models.Link) error {
	const op = "vray.UpdateKey"

	urlW := v.cfg.UrlVray + "panel/api/inbounds/updateClient/" + link.VpnLinkId + "?id=1&settings="

	isEnable := false

	if link.State == models.StateAllowed {
		isEnable = true
	}
	settings := &settingsJSON{
		Clients: clientsJSON{
			{
				ID:      link.VpnLinkId,
				Flow:    "xtls-rprx-vision",
				Email:   link.Link,
				LimitIP: v.cfg.Vray.LimitConnection,
				TotalGB: 0,
				Enable:  isEnable,
				TgID:    "",
				SubID:   "",
				Reset:   0,
			},
		},
	}
	err := v.requestDo(urlW, settings)
	if err != nil {
		log.Println(op, err)
		return err
	}

	return nil
}

func (v *Vray) SetTraficLimit(key string, limitGb int) error {
	return nil
}
func (v *Vray) DeleteKey(string) error {
	return nil
}

func (v *Vray) requestDo(urlW string, settings *settingsJSON) error {
	const op = "vtoray.requestDo"
	myBytes, err := json.Marshal(settings)
	if err != nil {
		log.Println(op, err)
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", urlW+url.QueryEscape(string(myBytes)), nil)
	if err != nil {
		log.Println(op, err)
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.AddCookie(&http.Cookie{Name: "session", Value: v.session})

	res, err := client.Do(req)
	if err != nil {
		log.Println(op, err)
		return err
	}
	defer res.Body.Close()
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
