package main

import (
	"net/http"
	"time"
)

type Response struct {
	Origin         string      `json:"origin"`
	Method         string      `json:"method"`
	Status         int         `json:"status"`
	ContentType    string      `json:"content_type"`
	ContentLength  uint        `json:"content_length"`
	Host           string      `json:"host"`
	Port           string      `json:"port"`
	URL            string      `json:"url"`
	Scheme         string      `json:"scheme"`
	Path           string      `json:"path"`
	Extension      string      `json:"ext"`
	ResponseHeader http.Header `json:"response_header,omitempty"`
	ResponseBody   string      `json:"response_body,omitempty"`
	RequestHeader  http.Header `json:"request_header,omitempty"`
	RequestBody    string      `json:"request_body,omitempty"`
	DateStart      time.Time   `json:"date_start"`
	DateEnd        time.Time   `json:"date_end"`
}
