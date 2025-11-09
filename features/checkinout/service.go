package checkinout

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	LOGIN_URL      = "https://blueprint.cyberlogitec.com.vn/sso/login"
	CHECKINOUT_URL = "https://blueprint.cyberlogitec.com.vn/api/checkInOut/insert"
)

func setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
}

func doRequest(client *http.Client, method, url string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	setBrowserHeaders(req)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return client.Do(req)
}

func Run(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username, password are required")
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	resp, err := doRequest(client, "GET", LOGIN_URL, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		return fmt.Errorf("expected 302, got %d", resp.StatusCode)
	}
	location := resp.Header.Get("Location")
	resp2, err := doRequest(client, "GET", location, nil, "")
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp2.Body)
	if err != nil {
		return err
	}
	formAction, exists := doc.Find("form#login-form").Attr("action")
	if !exists {
		return errors.New("login form action not found")
	}

	loginData := fmt.Sprintf("username=%s&password=%s", username, password)
	resp3, err := doRequest(client, "POST", formAction, strings.NewReader(loginData), "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != 302 {
		return fmt.Errorf("expected 302 after login, got %d", resp3.StatusCode)
	}
	location = resp3.Header.Get("Location")
	resp4, err := doRequest(client, "GET", location, nil, "")
	if err != nil {
		return err
	}
	defer resp4.Body.Close()
	if resp4.StatusCode != 302 {
		return fmt.Errorf("expected 302 after redirect 1, got %d", resp4.StatusCode)
	}
	location = resp4.Header.Get("Location")
	resp5, err := doRequest(client, "GET", location, nil, "")
	if err != nil {
		return err
	}
	defer resp5.Body.Close()

	resp6, err := doRequest(client, "POST", CHECKINOUT_URL, nil, "")
	if err != nil {
		return err
	}
	defer resp6.Body.Close()
	body, _ := io.ReadAll(resp6.Body)
	fmt.Printf("[SUCCESS] Check-in/out API response: %s\n", string(body))

	return nil
}
