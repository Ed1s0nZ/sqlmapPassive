package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
)

type ReqJSONData struct {
	Origin         string              `json:"origin"`
	Method         string              `json:"method"`
	Status         int                 `json:"status"`
	ContentType    string              `json:"content_type"`
	ContentLength  int                 `json:"content_length"`
	Host           string              `json:"host"`
	Port           string              `json:"port"`
	URL            string              `json:"url"`
	Scheme         string              `json:"scheme"`
	Path           string              `json:"path"`
	Ext            string              `json:"ext"`
	ResponseHeader map[string][]string `json:"response_header"`
	ResponseBody   string              `json:"response_body"`
	RequestHeader  map[string][]string `json:"request_header"`
	RequestBody    string              `json:"request_body"`
	DateStart      string              `json:"date_start"`
	DateEnd        string              `json:"date_end"`
}

func formatRequest(jsonStr string) (string, error) {
	// fmt.Println(jsonStr)
	var data ReqJSONData
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	// 写请求行
	fmt.Fprintf(&sb, "%s %s HTTP/1.1\n", data.Method, data.URL)

	// 写Host头
	fmt.Fprintf(&sb, "Host: %s\n", data.Host)

	// 写其他请求头
	for key, values := range data.RequestHeader {
		for _, value := range values {
			fmt.Fprintf(&sb, "%s: %s\n", strings.ToLower(key), value)
		}
	}

	// 添加连接头
	fmt.Fprintf(&sb, "Connection: close\n")

	// 添加Content-Type和Content-Length头，并添加请求体
	if data.RequestBody != "" {
		// fmt.Fprintf(&sb, "Content-Type: %s\n", data.ContentType)
		// fmt.Fprintf(&sb, "Content-Length: %d\n\n", len(data.RequestBody))
		fmt.Fprintf(&sb, "\n%s", data.RequestBody)
	} else {
		fmt.Fprint(&sb, "\n")
	}

	return sb.String(), nil
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func handleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	reqbody, err := RequestBody(req)
	checkErr(err)
	RequestBodyMap.Store(ctx.Session, reqbody)
	return req, nil
}

func RequestBody(res *http.Request) ([]byte, error) {

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = io.NopCloser(bytes.NewReader(buf))
	return buf, nil
}

// json.Marshal方法优化，不对html做转义处理
func MarshalHTML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ResponseBody(res *http.Response) ([]byte, error) {
	if res != nil {
		defer res.Body.Close()
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = io.NopCloser(bytes.NewReader(buf))
	return buf, nil
}

func New(resp *http.Response, reqbody []byte, respbody []byte) *ParserHTTP {
	return &ParserHTTP{r: resp, reqbody: reqbody, respbody: respbody, s: time.Now()}
}

func NewResType(ext string, ctype string) *ResType {
	var mtype string
	if ctype != "" {
		mtype = strings.Split(ctype, "/")[0]
	}
	return &ResType{ext, ctype, mtype}
}

type ParserHTTP struct {
	r        *http.Response
	reqbody  []byte
	respbody []byte
	s        time.Time
}

type ResType struct {
	ext   string
	ctype string
	mtype string
}

func (parser *ParserHTTP) Parser() Response {

	var (
		ctype   string
		clength int
		StrHost string
		StrPort string
	)

	if len(parser.r.Header["Content-Type"]) >= 1 {
		ctype = GetContentType(parser.r.Header["Content-Type"][0])
	}

	if len(parser.r.Header["Content-Length"]) >= 1 {
		clength, _ = strconv.Atoi(parser.r.Header["Content-Length"][0])
	}

	SliceHost := strings.Split(parser.r.Request.URL.Host, ":")
	if len(SliceHost) > 1 {
		StrHost, StrPort = SliceHost[0], SliceHost[1]
	} else {
		StrHost = SliceHost[0]
		if parser.r.Request.URL.Scheme == "https" {
			StrPort = "443"
		} else {
			StrPort = "80"
		}
	}

	now := time.Now()

	r := Response{
		Origin:         parser.r.Request.RemoteAddr,
		Method:         parser.r.Request.Method,
		Status:         parser.r.StatusCode,
		ContentType:    string(ctype),
		ContentLength:  uint(clength),
		Host:           StrHost,
		Port:           StrPort,
		URL:            parser.r.Request.URL.String(),
		Scheme:         parser.r.Request.URL.Scheme,
		Path:           parser.r.Request.URL.Path,
		Extension:      GetExtension(parser.r.Request.URL.Path),
		ResponseHeader: parser.r.Header,
		ResponseBody:   string(parser.respbody),
		RequestHeader:  parser.r.Request.Header,
		RequestBody:    string(parser.reqbody),
		DateStart:      parser.s,
		DateEnd:        now,
	}

	return r
}

func (r *ResType) isStatic() bool {
	if ContainsString(static_ext, r.ext) {
		return true
	} else if ContainsString(static_types, r.ctype) {
		return true
	} else if ContainsString(media_types, r.mtype) {
		return true
	}
	return false
}

func GetContentType(HeradeCT string) string {
	ct := strings.Split(HeradeCT, "; ")[0]
	return ct
}

func GetExtension(path string) string {
	SlicePath := strings.Split(path, ".")
	if len(SlicePath) > 1 {
		return SlicePath[len(SlicePath)-1]
	}
	return ""
}

func ContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func handleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	// Getting the Body
	reqbody, ok := RequestBodyMap.Load(ctx.Session)
	RequestBodyMap.Delete(ctx.Session)
	if ok && resp != nil {
		respbody, err := ResponseBody(resp)
		checkErr(err)
		// Attaching capture tool.
		if respbody != nil {
			RespCapture := New(resp, reqbody.([]byte), respbody).Parser()

			static := NewResType(
				RespCapture.Extension,
				RespCapture.ContentType).isStatic()
			if !static {
				jsonStr, err := MarshalHTML(RespCapture)
				if err != nil {
					log.Fatal()
				}
				go func() {
					a, err := formatRequest(string(jsonStr))
					if err != nil {
						fmt.Println(err)
					}
					// fmt.Println(a)
					err1 := calculateAndSaveMD5(a)
					if err1 != nil {
						fmt.Println(err1)
					}
				}()

			}
		}
	}
	return resp
}
func proxy() {
	fmt.Println("Proxy start")
	// 定义代理日志目录
	_dir := "log"
	exist, err := PathExists(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	if exist {
		fmt.Printf("Proxy log dir -> [%v]\n", _dir)
	} else {
		fmt.Printf("No proxy log dir -> [%v]\n", _dir)
		// 创建代理目录
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Mkdir proxy log failed![%v]\n", err)
		} else {
			fmt.Printf("Mkdir proxy log success!\n")
		}
	}
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("l", ":3234", "on which address should the proxy listen")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose
	log.Printf("Listening %s \n", *addr)
	log.Printf("proxy Start success... \n")
	log.Println(goproxy.ReqHostMatches())
	go proxy.OnRequest().DoFunc(handleRequest)
	go proxy.OnResponse().DoFunc(handleResponse)
	proxy.OnResponse().DoFunc(handleResponse)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
