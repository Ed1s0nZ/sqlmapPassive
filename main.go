package main

const (
	dir        = "./sqlmapReq"
	numWorkers = 3   // 调节线程数 （启动3个线程调用sqlmap进行扫描）
	numTxt     = 100 // 一次最多存储多少个要扫描的txt在sqlmapReq文件夹下
	proxyPort  = "3231" // 此工具用作被动代理的监听端口
)

func main() {
	go sqlmap()
	go proxy()
	go delCsv()
	select {}
}
