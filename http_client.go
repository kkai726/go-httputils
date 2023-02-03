package httputils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Headers map[string]string

type Cookies []*http.Cookie

type Data map[string]interface{}

// Params
type Params struct {
	Data    Data
	Headers Headers
	Cookies Cookies
}

// Respones
type responses struct {
	Response *http.Response
	Body     string
}

// Debugger
var Debugger bool = false

func request(method string, url string, params Params) (r responses, e error) {
	// print the debugger to initiate the request
	printSendDebugger(method, url, params)
	r, e = do(method, url, params)
	if e != nil {
		return
	}
	// print debugger return request
	printResultDebugger(r)
	return
}

func do(method string, url string, params Params) (r responses, e error) {
	var data *bytes.Buffer

	// params.Data cannot be empty, otherwise, a null pointer will be reported.
	if params.Data == nil {
		params.Data = Data{
			"placeholder": " ",
		}
	}

	jsonStr, _ := json.Marshal(params.Data)
	data = bytes.NewBuffer(jsonStr)

	request, e := http.NewRequest(method, url, data)
	// add header
	addHeaders(request, params.Headers)
	// add Cookies
	addCookies(request, params.Cookies)
	// add Data
	addData(request, params.Data)
	r.Response, e = http.DefaultClient.Do(request)
	if e != nil {
		return
	}
	defer r.Response.Body.Close()
	unCoding(&r)
	return
}

// add headers
func addHeaders(request *http.Request, headers Headers) {
	for k, v := range headers {
		request.Header.Add(k, v)
	}
}

// add Cookies
func addCookies(request *http.Request, cookies Cookies) {
	for _, v := range cookies {
		request.AddCookie(v)
	}
}

// add Data
func addData(request *http.Request, data Data) {
	query := request.URL.Query()
	for k, v := range data {
		query.Add(k, fmt.Sprint(v))
	}
	request.URL.RawQuery = query.Encode()

}

func unCoding(r *responses) {
	if r.Response.StatusCode == http.StatusOK {
		switch r.Response.Header.Get("Contend-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(r.Response.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)
				if err != nil && err != io.EOF {
					return
				}
				if n == 0 {
					break
				}
				r.Body += string(buf)
			}
		default:
			bodyByte, _ := io.ReadAll(r.Response.Body)
			r.Body = string(bodyByte)
		}
	} else {
		bodyByte, _ := io.ReadAll(r.Response.Body)
		r.Body = string(bodyByte)
	}
}

// debugger make a request
func printSendDebugger(method string, url string, params Params) {
	if Debugger {
		log.Println("debug log start ------------------------------------------------------")
		fmt.Println("Method", method)
		fmt.Println("Host", ":", url)
		for k, v := range params.Headers {
			fmt.Println(k, ":", v)
		}
		fmt.Println("-----------------------------------------------------------------------")
	}
}

// debugger request result
func printResultDebugger(r responses) {
	if Debugger {
		fmt.Println("Status", ":", r.Response.Status)
		for key, val := range r.Response.Header {
			fmt.Println(key, ":", val[0])
		}
		log.Println("debug log end --------------------------------------------------------")
	}
}

// request
func Request(method string, url string, params Params) (responses, error) {
	return request(method, url, params)
}

// get
func GetRequest(url string, params Params) (responses, error) {
	return request("GET", url, params)
}

// post
func PostRequest(url string, params Params) (responses, error) {
	return request("POST", url, params)
}

// put
func PutRequest(url string, params Params) (responses, error) {
	return request("PUT", url, params)
}

// delete
func DeleteRequest(url string, params Params) (responses, error) {
	return request("DELETE", url, params)
}
