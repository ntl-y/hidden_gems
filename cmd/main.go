package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	server     = "api.vk.ru"
	methodName = "video.get"
	apiVersion = "5.199"
)

func gatherParams() url.Values {
	params := url.Values{}

	params.Add("q", "майнкрафт")
	// params.Add("sort", "0")
	// params.Add("adult", "1")
	params.Add("filters", "vk")
	// params.Add("search_own", "0")
	// params.Add("offset", "0") // смещение
	params.Add("count", "200")
	params.Add("longer", "5")
	// params.Add("shorter", "6000")
	params.Add("extended", "0")
	params.Add("v", apiVersion)

	return params

}

func createRequest() *http.Request {
	bearer := "Bearer " + os.Getenv("VK_TOKEN")
	params := gatherParams().Encode()

	url := fmt.Sprintf("https://%s/method/%s?%s", server, methodName, params)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/x-www-form-encoded")

	return req
}

func sendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := godotenv.Load(); err != nil {
		logrus.Fatalln(err)
	}

	client := &http.Client{}
	req := createRequest()
	resp, err := sendRequest(client, req)
	if err != nil {
		logrus.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Println(err)
	}

	err = os.WriteFile("response.json", body, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}

	fmt.Println("Response saved to response.json")
}
