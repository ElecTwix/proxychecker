
# Proxy Checker

## Get

```
go get github.com/ElecTwix/proxychecker
```

## Example Usage

``` go

import "github.com/ElecTwix/proxychecker"


var strarr []string
secure := false

arr, err := proxychecker.checkall(&strarr, "http://github.com/", time.Second * 5, &secure)
if err != nil {
    panic(err.Error())
}


fmt.Println("Working proxies", arr)

 ```



