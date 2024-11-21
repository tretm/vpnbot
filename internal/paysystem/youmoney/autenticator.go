package youmoney

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type authorize struct {
	ClientID    string
	RedirectURI string
	Scope       []string
}

func (a *authorize) authorize() (string, error) {
	scope := strings.Join(a.Scope, "%20")
	authURL := fmt.Sprintf("https://yoomoney.ru/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
		url.QueryEscape(a.ClientID), url.QueryEscape(a.RedirectURI), scope)

	fmt.Println("Visit this website and confirm the application authorization request:")
	fmt.Println(authURL)

	fmt.Print("Enter redirected url (https://yourredirect_uri?code=XXXXXXXXXXXXX) or just code: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	code := scanner.Text()
	// code := "code=EAA925AEAB426EC63F024686BC058A9568C3726D71C240113203C6679A337742565ECBA2457D7FDCAC91090B68FFD6FC3CF1A97D79B406B9D12560837EF2DEDE0455A558D0579701A32DE93D8F66FA8FB994782D0590837B4EDAECCDABC678A3F8D353C0EDEA64F0747EEA0566DC6E560A89EA6B06B44695869292FDFB2A14F6"
	if strings.Contains(code, "code=") {
		code = strings.Split(code, "code=")[1]
	}

	tokenURL := fmt.Sprintf("https://yoomoney.ru/oauth/token?code=%s&client_id=%s&grant_type=authorization_code&redirect_uri=%s",
		url.QueryEscape(code), url.QueryEscape(a.ClientID), url.QueryEscape(a.RedirectURI))

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to obtain token")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if err, exists := result["error"]; exists {
		switch err {
		case "invalid_request":
			return "", errors.New("invalid request")
		case "unauthorized_client":
			return "", errors.New("unauthorized client")
		case "invalid_grant":
			return "", errors.New("invalid grant")
		default:
			return "", fmt.Errorf("unknown error: %v", err)
		}
	}

	tt := ""
	if token, exists := result["access_token"]; exists {
		fmt.Println("Your access token:")
		fmt.Println(token)
		tt = token.(string)
	} else {
		return "", errors.New("empty token")
	}

	return tt, nil
}
