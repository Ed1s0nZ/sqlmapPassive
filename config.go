package main

import "sync"

var (
	RequestBodyMap sync.Map

	// http static resource file extension
	static_ext []string = []string{
		"js",
		"css",
		"ico",
		"woff",
		"ttf",
		"map",
		"woff2",
	}

	// media resource files type
	media_types []string = []string{
		"image",
		"video",
		"audio",
	}

	// http static resource files
	static_types []string = []string{
		"application/vnd.google.octet-stream-compressible",
		"font/woff",
		"font/woff2",
		"text/css",
		"text/javascript",
		"baiduApp/json",
		"application/javascript",
		"application/x-javascript",
		"application/msword",
		"application/vnd.ms-excel",
		"application/vnd.ms-powerpoint",
		"application/x-ms-wmd",
		"application/x-shockwave-flash",
	}
)
