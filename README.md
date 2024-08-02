# sqlmapPassive
sqlmap被动代理版

## 用法
1. 下载原文件并打包；
2. 将自己的sqlmap文件夹放到该目录下
3. 修改sqlmap的源码的lib/core/option.py 中的init()方法中设置conf.stdinPipe=None
4. burp或者浏览器挂代理到该工具的监听端口，即可自动化扫描。
