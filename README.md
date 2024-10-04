# SIDT_Final_Practice

:globe_with_meridians:此仓库内容为大四上最后一个实习课：空间数据数字工程综合实习。

采用前后端分离的架构：

+ 前端：
  + Vue
  + JS
  + HTML + CSS
  + TypeScript
+ 后端
  + Golang

## 使用说明：

在 VS Code中打开项目，然后在终端中输入以下命令：

```shell
cd backend; go run main.go
```

运行前请确保 Go 语言环境已经被正确安装。

如果在终端看到如下 `CengKeHelper` 样的字符信息，则说明服务启动成功：

![](.\docs\pic\Snipaste_2024-10-04_19-07-05.png)

可以在浏览器中输入以下地址查看服务：

```http
http://localhost:8000
```

![](.\docs\pic\Snipaste_2024-10-04_19-17-00.png)

## 文件架构：

+ **backend**：后端 Go 语言代码存放。
  + **api**：存放后端 API 代码。
  + **logger**：实现了一个能够根据不同平台（Linux 或 Windows）输出彩色日志的系统，并同时将日志写入文件。
  + **logs**：后端运行时日志文件回存放于此。
  + **process**：存放数据处理部分的代码。
  + **setup**：存放配置文件。
  + **whu**：用于存放数据文件
+ **frontend**：存放前端代码。
+ **docs**：需求分析档案和图片。
