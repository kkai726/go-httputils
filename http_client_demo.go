package httputils

import (
	"fmt"
	"os"
)

func main() {
	Debugger = true
	var params Params
	params.Data = Data{
		"id":     "admin",
		"name":   "admin",
		"gender": "man",
	}

	params.Cookies = Cookies{
		{Name: "id", Value: "admin"},
		{Name: "name", Value: "name"},
	}

	params.Headers = Headers{
		"Content-Type":    "multipart/form-data",
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "keep-alive",
	}
	r, e := Request("GET", "http://localhost/test", params)
	if e != nil {
		fmt.Println(e)
		os.Exit(0)
	}
	fmt.Println(r.Body)
}
