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
	ApiKey	string
	Client	*http.Client
}

type Response struct{
	Request	string
	Status	float64
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

// A method create job and return the ID
func (c *CaptchaServiceClient) CreatTask(dataForm url.Values) (string, error){
	// Formalize the POST request
	resp, err := http.Post(baseURL.ResolveReference(&url.URL{Path: "/in.php"}).String(), "application/x-www-form-urlencoded", strings.NewReader(dataForm.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var responseBody Response
	json.NewDecoder(resp.Body).Decode(&responseBody)

	if responseBody.Status != 1{
			return "", errors.New(responseBody.Request)
	}

	return responseBody.Request, nil
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
	var responseBody Response
	json.NewDecoder(resp.Body).Decode(&responseBody)

	// Interpret the answer
	switch responseBody.Request{
		case "CAPCHA_NOT_READY":
			return "", false, nil
		default:
			if responseBody.Status != 1{
				return "", true, errors.New(responseBody.Request)
			}
	}

	return responseBody.Request, true, nil
}

// A method that creates job with Google reCaptcha and gives you the answer
func (c *CaptchaServiceClient) ReCaptchaV2(siteURL string, siteKey string, timeout time.Duration) (string, error){
	// Preparation form data
	data := url.Values{
		"key":			{c.ApiKey},
		"method":		{"userrecaptcha"},
		"googlekey":		{siteKey},
		"pageurl":		{siteURL},
		"json":			{"1"},
	}

	// Create task to resolve
	jobID, err := c.CreatTask(data)
	if err != nil{
		return "", err
	}

	// Check the status of job every 5 seconds
	ping := time.NewTicker(5 * time.Second)
	timeoutBreak := time.NewTimer(timeout * time.Second)

	for {
		select {
		case <-ping.C:
			answer, ready, err := c.CheckResult(jobID)
			if err != nil{
				return "", err
			}

			if ready{
				return answer, nil
			}
		case <-timeoutBreak.C:
			return "", errors.New("Waiting for captcha result timeout")
		}
	}

	// If the job takes too long
	return "", errors.New("Waiting for captcha result timeout")
}

// A method that creates job with regular image captcha in base64 and gives you the answer
func (c *CaptchaServiceClient) RegularCaptcha(base64Image string, timeout time.Duration) (string, error){
	// Preparation form data
	data := url.Values{
		"key":		{c.ApiKey},
		"method":	{"base64"},
		"body":		{base64Image},
		"json":		{"1"},
	}

	// Create task to solve
	jobID, err := c.CreatTask(data)
	if err != nil{
		return "", err
	} 

	// Check the status of job every 5 seconds
	ping := time.NewTicker(5 * time.Second)
	timeoutBreak := time.NewTimer(timeout * time.Second)

	for {
		select {
		case <-ping.C:
			answer, ready, err := c.CheckResult(jobID)
			if err != nil{
				return "", err
			}

			if ready{
				return answer, nil
			}
		case <-timeoutBreak.C:
			return "", errors.New("Waiting for captcha result timeout")
		}
	}

	// If the job takes too long
	return "", errors.New("Waiting for captcha result timeout")
}
