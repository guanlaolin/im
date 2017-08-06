# 基于WebSocket协议的即时通信系统

<strong>本程序在实现之初由于实在是赶时间，各种考虑不周，各种小bug，各种不优雅的代码，现正进行重构，所以几乎每天代码都可能会变</strong><br />

<strong>本项目已切换到分支v2中继续进行开发，在v2开发完成之前，master只进行bug修复，代码不再修改</strong> <br />
v2主要修改如下：<br />
1、使用RESTful风格 <br/>
2、路由使用github.com/gorilla/mux <br />
3、改cookie为session，使用github.com/gorilla/sessions <br />
4、数据库使用github.com/jinzhu/gorm <br />

<strong>说明：</strong> <br />
这是本人的毕业设计题目，获得了优秀毕业设计，在这里共享代码仅供参考。

<strong>介绍：</strong> <br />
随着网络技术爆炸式发展，B/S模式以其简单、方便的优势，越来越受欢迎，其功能也越来越强大。
但由于HTTP协议的限制，实时通信一直是B/S模式的短板。传统解决方案如轮询（polling）和Comet技术，
存在着实时性低、开发复杂度和资源消耗高等缺点，不是理想的解决方案。
为了解决这些问题，系统使用Golang语言表述WebSocket协议实现即时通信，实现了真正意义上的B/S模式的即时通信，且开发复杂度和资源消耗都相对较低。
系统还使用WebRTC实现视频通信，实现了B/S模式的点对点视频通信。
此外，为了提高用户体验，系统提供了用户注册与登录、用户好友管理等功能以及友好的用户界面，实现了用户信息的自我管理以及好友等自我管理。
经过详细、周密的测试，本系统用户友好、高效率、高容错、高安全性、高并发。

<strong>功能：</strong><br />
1、基本的聊天功能；
2、表情；
3、文件传输；
4、视频通信。
<br />

文件传输使用了HTML5的FileReader接口，进行分片传送，保证了大文件传输浏览器端不卡死；后期可能会改用WebRTC传输，保证效率。 <br />
视频通信使用的是WebRTC，目前使用的是公用ICE实现外网穿透。

<strong> 补充：</strong><br />
代码细节之处还需优化，正在不断改进。

由于最近在忙于找工作，没啥时间写。
