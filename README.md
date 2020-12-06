# Golang universal captcha resolve service

[![GoDoc](https://godoc.org/github.com/xta/okrun?status.svg)](https://pkg.go.dev/github.com/Flagon00/CaptchaResolveApiClient)

One package for all services that support /in.php and /resp.php api. Tested on cptch.net, 2captcha.com, XEvil and capmonster.cloud

Setup:
```go get -u github.com/Flagon00/CaptchaResolveApiClient```

Example usage with reCaptchaV2 and cptch.net:
```go
package main

import (
    "log"

    "github.com/Flagon00/CaptchaResolveApiClient"
)

func main() {
    client := captcha.Client(true, "cptch.net", "api-key")
    resolve, err := client.ReCaptchaV2("https://www.google.com/recaptcha/api2/demo",  "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-")
    if err != nil{
        log.Fatal(err)
    }
    log.Println(resolve)
}
```

Also usage with image captcha example:
```go
client := captcha.Client(true, "cptch.net", "api-key")
resolve, err := client.RegularCaptcha("base64-string")
if err != nil{
    log.Fatal(err)
}
log.Println(resolve)
```

Example client for 2captcha:
```go
captcha.Client(true, "2captcha.com", "api-key")
```

Or if you want, you can use this package with XEvil:
```go
captcha.Client(false, "localhost", "api-key")
```
