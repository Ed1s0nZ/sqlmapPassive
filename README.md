# sqlmapPassive
sqlmap被动代理小工具

## 用法
1. 下载原文件并解压；
2. 将自己的sqlmap文件夹放到该目录下；
3. 修改sqlmap源码lib/core/option.py 中的 init() 方法，设置`conf.stdinPipe=None`；
   <img src="https://github.com/Ed1s0nZ/sqlmapPassive/blob/main/sqlmap.png" alt="sqlmap配置" width="400"/>
5. 主要配置在main.go，可配置线程数、每次最多存储多少个要扫描的txt在sqlmapReq文件夹下和代理端口等；
6. 配置完毕后编译并运行，编译：`go build`，运行：`./PassiveSqlmap`；
7. burp或者浏览器挂代理到该工具的监听端口，即可通过被动代理的形式进行扫描。

## 效果
开启三个线程进行扫描的效果图：   
<img src="https://github.com/Ed1s0nZ/sqlmapPassive/blob/main/xiaoguo.png" alt="效果图" width="700"/>
   
（自己用起来还可以，目前暂未发现其他问题。还有很多要优化的点，比如发现问题后的告警机器人等）
