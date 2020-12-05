package captcha

import (
	"errors"
	"encoding/json"
	"net/http"
	"time"
	"net/url"
	"strings"
)

type CaptchaServiceClient struct {
	ApiKey 		string
	Client 		*http.Client
}

var baseURL *url.URL

// Preparation client to use package
func Client(secure bool, provider string, apikey string) *CaptchaServiceClient{
	switch secure{
		case false:
			baseURL = &url.URL{Host: provider, Scheme: "http", Path: "/"}
		default:
			baseURL = &url.URL{Host: provider, Scheme: "https", Path: "/"}
	}

	return &CaptchaServiceClient{
		ApiKey: apikey,
		Client: http.DefaultClient,
	}
}

// A method to check answer for job
func (c *CaptchaServiceClient) CheckResult(captchaId string) (string, bool, error){
	// Preparation form data
	data := url.Values{
		"key":		{c.ApiKey},
		"action":	{"get"},
		"id":		{captchaId},
		"json":		{"1"},
	}

	// Formalize the POST request
	resp, err := http.Post(baseURL.ResolveReference(&url.URL{Path: "/res.php"}).String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", true, err
	}
	defer resp.Body.Close()

	// Decode the response
	responseBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseBody)

	// Interpret the answer
	switch responseBody["request"].(string){
		case "CAPCHA_NOT_READY":
			return "", false, nil
		default:
			if responseBody["status"].(float64) != 1{
				return "", true, errors.New(responseBody["request"].(string))
			}
	}

	return responseBody["request"].(string), true, nil
}

// A method that creates job with Google reCaptcha and gives you the answer
func (c *CaptchaServiceClient) ReCaptchaV2(siteURL string, siteKey string) (string, error){
	// Preparation form data
	data := url.Values{
		"key":		{c.ApiKey},
		"method":	{"userrecaptcha"},
		"googlekey":	{siteKey},
		"pageurl":	{siteURL},
		"json":		{"1"},
	}

	// Formalize the POST request
	resp, err := http.Post(baseURL.ResolveReference(&url.URL{Path: "/in.php"}).String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	responseBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseBody)

	// Interpret the answer
	if _, ok := responseBody["status"]; !ok {
		return "", errors.New("Wrong resopnse struct")
	}

	if responseBody["status"].(float64) != 1{
			return "", errors.New(responseBody["request"].(string))
	}
	
	// Check the status of job every 5 seconds
	for t := 0; t <= 20; t++{
		answer, ready, err := c.CheckResult(responseBody["request"].(string))
		if err != nil{
			return "", err
		}

		if ready{
			return answer, nil
		}

		time.Sleep(time.Second * 5)
	}

	// If the job takes too long
	return "", errors.New("Too many tries")
}

// A method that creates job with regular image captcha in base64 and gives you the answer
func (c *CaptchaServiceClient) RegularCaptcha(base64Image string) (string, error){
	// Preparation form data
	data := url.Values{
		"key":		{c.ApiKey},
		"method":	{"base64"},
		"body":		{base64Image},
		"json":		{"1"},
	}

	// Formalize the POST request
	resp, err := http.Post(baseURL.ResolveReference(&url.URL{Path: "/in.php"}).String(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	responseBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseBody)

	// Interpret the answer
	if _, ok := responseBody["status"]; !ok {
		return "", errors.New("Wrong resopnse struct")
	}
	if responseBody["status"].(float64) != 1{
		return "", errors.New(responseBody["request"].(string))
	}

	// Check the status of job every 5 seconds
	for t := 0; t <= 20; t++{
		answer, ready, err := c.CheckResult(responseBody["request"].(string))
		if err != nil{
			return "", err
		}

		if ready{
			return answer, nil
		}

		time.Sleep(time.Second * 5)
	}

	// If the job takes too long
	return "", errors.New("Too many tries")
}
