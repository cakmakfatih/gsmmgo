package internals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

var C Config

type Config struct {
	PbAdmin *hClient
	BaseURL string
}

type hClient struct {
	c       *http.Client
	baseURL string
}

func (h *hClient) Post(path string, body io.Reader) (resp *http.Response, err error) {
	return h.c.Post(fmt.Sprintf("%v%v", h.baseURL, path), "application/json", body)
}

func (h *hClient) Get(path string) (resp *http.Response, err error) {
	return h.c.Get(fmt.Sprintf("%v%v", h.baseURL, path))
}

func (h *hClient) NewRequest(method string, path string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", h.baseURL, path), body)
	if err != nil {
		return nil, err
	}

	return h.c.Do(req)
}

func PostWithToken(path string, body io.Reader, token string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%v%v", C.BaseURL, path), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	return client.Do(req)
}

func GetWithToken(path string, token string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v%v", C.BaseURL, path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	return client.Do(req)
}

func (h *hClient) NewRequestWithToken(method string, path string, body io.Reader, token string) (resp *http.Response,
	err error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", C.BaseURL, path), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	return client.Do(req)
}

func newHClient(token string, baseURL string) *hClient {
	return &hClient{
		c: &http.Client{
			Transport: &transport{
				UnderlyingTransport: http.DefaultTransport,
				Token:               token,
			},
		},
		baseURL: baseURL,
	}
}

type AdminAuthRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type AdminAuthResponse struct {
	Token string `json:"token"`
}

type transport struct {
	UnderlyingTransport http.RoundTripper
	Token               string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", t.Token))

	return t.UnderlyingTransport.RoundTrip(req)
}

func pbAdminAuth() (string, error) {
	url := fmt.Sprintf("%v/api/admins/auth-with-password", os.Getenv("PB_URL"))

	jsonData, err := json.Marshal(&AdminAuthRequest{
		Identity: os.Getenv("PB_ADMIN_IDENTITY"),
		Password: os.Getenv("PB_ADMIN_PASSWORD"),
	})

	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return "", err
	}

	var result AdminAuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return "", err
	}

	return result.Token, nil
}

func InitConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token, err := pbAdminAuth()

	if err != nil {
		log.Fatal("Error authenticating to PocketBase")
	}

	C = Config{
		PbAdmin: newHClient(token, os.Getenv("PB_URL")),
		BaseURL: os.Getenv("PB_URL"),
	}
}
