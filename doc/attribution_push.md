#DeepShare Attribution Push

深享系统可以提供app安装后用户事件和流量源头的关联分析。这些事件包括DeepShare内置事件安装和打开，也可以是用户自定义的事件，比如购买，分享。DeepShare归因分析可以针对渠道或者发送者，其中渠道的归因分析可以通过WEB来浏。而对于基于发生者的归因分析，由于数据量巨大，DeepShare提供自动Push功能，及时的将归因分析的结果Push到开发者提供的后台。该后台只需要提供一个HTTP/Restful API，能够接受如下POST：

```json
[
    {
        "sender_id": "Who is sending the url",
        "tag": "type of action, include ds/open, ds/install, and other user defined action",
        "value" : "the amount associated with tag",
        "timestamp" : "when the amount added"
    },
    {}
]
```

##Field 详细说明:
- senderID：DeepShare URL 的生成者ID, 此ID就是SDK中的API getSenderId() 的返回值，开发者可以绑定此ID和自己的用户系统
 
- tag & value: 
    ds/install, 表示通过此senderID带来的新安装数量，value为其数量
    ds/open, 表示通过此senderID带来的新打开数量，value为其数量
    自定义tag, 当分享接收方在SDK中调用attribute(tagToValue, DSFailListener)方法时，其value关联到此tag，归因在其发送方并返回。

该后台应该返回HTTP200.

##补充:
应该说明的是：DeepShare的归因分析是基于流的，会有分钟级别的延时。如果开发者后端不能及时接收，相应的数据会被延时再次Push（时间开发者可以定义）。如果仍然不行，就会被Attribution Push丢掉。