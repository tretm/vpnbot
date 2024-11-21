package outline

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"vpngigabot/internal/config"
	"vpngigabot/internal/models"
)

type Outline struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Outline {
	return &Outline{cfg: cfg}
}

type CreateKeyResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Method    string `json:"method"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	AccessUrl string `json:"accessUrl"`
}
type CreateKeyRequest struct {
	Name     string    `json:"name"`
	Method   string    `json:"method"`
	Password string    `json:"password"`
	Port     int       `json:"port"`
	Limit    dataLimit `json:"limit"`
}
type dataLimit struct {
	Bytes int `json:"bytes"`
}

func (o *Outline) UpdateKey(key *models.Link) error {
	return nil
}
func (o *Outline) KeyBuilder(key *models.Link) string {
	return ""
}
func (o *Outline) CreateKey(usr *models.User) (*models.Link, error) {

	keyReq := CreateKeyRequest{Name: strconv.FormatInt(usr.UserId, 10),
		Method:   "chacha20-ietf-poly1305",
		Password: strconv.FormatInt(usr.UserId, 10)}
	keyReqByte, err := json.Marshal(keyReq)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Пропускаем проверку сертификата
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", o.cfg.UrlOutline+"access-keys", bytes.NewBuffer(keyReqByte))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	keyResp := CreateKeyResponse{}
	err = json.Unmarshal(responseBody, &keyResp)
	if err != nil {
		return nil, err
	}
	pathword, err := getPaht(keyResp.AccessUrl)
	if err != nil {
		return nil, err
	}
	decodedString, err := base64.StdEncoding.DecodeString(pathword)
	if err != nil {
		return nil, err
	}
	pwd := ""
	extractedPathword := strings.Split(string(decodedString), ":")
	if len(extractedPathword) == 2 {
		pwd = extractedPathword[1]
	}
	link := &models.Link{UserId: usr.UserId,
		VpnLink:     keyResp.AccessUrl,
		VpnLinkId:   keyResp.Id,
		VpnPassword: pwd,
	}

	return link, nil
}

func (o *Outline) SetTraficLimit(key string, limitGb int) error {
	return nil
}

type DeleteKeyResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (o *Outline) DeleteKey(id string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Пропускаем проверку сертификата
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", o.cfg.UrlOutline+"access-keys/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(responseBody) == 0 {
		return nil
	}
	keyResp := DeleteKeyResponse{}
	err = json.Unmarshal(responseBody, &keyResp)
	if err != nil {
		return err
	}
	return errors.New(keyResp.Message)
}

func getPaht(accesKey string) (string, error) {
	re := regexp.MustCompile(`ss://(.*?)[@/]`)

	// Находим подстроку, соответствующую паттерну
	match := re.FindStringSubmatch(accesKey)
	if len(match) > 1 {
		// Извлекаем подстроку
		extractedString := match[1]
		return extractedString, nil
	} else {
		return "", errors.New("No match found")

	}

}
