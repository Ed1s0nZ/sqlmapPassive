# sqlmapPassive
sqlmap被动代理版

## 用法
1. 下载原文件并打包；
2. 将自己的sqlmap文件夹放到该目录下；
3. 修改sqlmap源码lib/core/option.py 中的 init() 方法，将 conf.stdinPipe=None；
4. burp或者浏览器挂代理到该工具的监听端口，即可通过被动代理的形式进行扫描。
