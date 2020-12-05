# Golang universal captcha resolve service

One package for all services that support /in.php and /resp.php api. Tested on cptch.net, 2captcha.com, XEvil and capmonster.cloud

Example usage with reCaptchaV2 and cptch.net:

```
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

Also usege with image captcha example:
```
client := captcha.Client(true, "cptch.net", "api-key")
resolve, err := client.RegularCaptcha("base64-string")
if err != nil{
    log.Fatal(err)
}
log.Println(resolve)
```

Example client for 2captcha:
```
captcha.Client(true, "2captcha.com", "api-key")
```

Or if you want, you can use this package with XEvil:
```
captcha.Client(false, "localhost", "api-key")
```