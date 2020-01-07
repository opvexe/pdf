### 1.使用说明 [json转excel]

```shell
Mac： 使用说明
$ ./main -i 输入文件目录 -o 输出文件目录
Liunx：
$ main -i 输入文件目录 -o 输出文件目录
Window:
$ main.exe -i 输入文件目录 -o 输出文件目录
```

#### 1.2使用说明 [pdf转jpg]

```shell
Mac： 使用说明
$ ./main -i 输入文件目录 -o 输出文件目录
Liunx：
$ main -i 输入文件目录 -o 输出文件目录
Window:
$ main.exe -i 输入文件目录 -o 输出文件目录
```

==**注意事项**==

```shell
# 相对路径 相对于main执行文件的路径
./main -i ../pdf/中文.pdf -o ../pdf/out  
```

#### 1.3 批量处理说明

```shell
Mac： 使用说明
$ ./main -f 文件目录 -o e:输入excel，默认json
Liunx：
$ main -f 文件目录 -o e:输入excel，默认json
Window:
$ main.exe -f 文件目录 -o e:输入excel，默认json

#使用案例
$ main -f ./input -o e  #[./input相对于可执行文件目录]
```

#### 1.4 window下gcc安装

```shell
# 第一步：官网[https://sourceforge.net/projects/tdm-gcc/]下载gcc
# 第二步：解压压缩包 [MinGW64]
# 第三步：将解压所得的压缩包，放置于D:\Tools\mingw64\MinGW64
# 第四步：配置系统环境变量 
# Path=D:\Tools\mingw64\MinGW64\bin 【此bin下一定要含有gcc.exe】
# 重启电脑 cmd --> gcc -v 查看
```

#### 1.5 window 下安装golang

```shell
# 第一步: 安装Goland
# 第二步：Goland --> setting ---> GoROOT --->download
# 第三步：设置环境变量【此电脑 --->属性 --->高级系统设置--->环境变量】
# GOROOT=D:\go\go1.13.5
# GOPATH=D:\workspace\go
# Path=%GOROOT%\bin
```

