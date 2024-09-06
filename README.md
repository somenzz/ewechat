# golang 发送企业微信通知的包

## 更新说明

- access_token 失效之前，避免重复调用接口获取 access_token


使用样例：

```golang
package main

import (
	"fmt"

	"github.com/somenzz/ewechat"
)

func main() {
    var ewechat = ewechat.EWechat{
        CorpID:     "your corpid",
        CorpSecret: "your corpsecret",
        AgentID:    your agentid,
    }
    
    msg := fmt.Sprintf("your message")
    ewechat.SendMessage(msg, "your enterprise wechat account, for more receiver, use like this UserID1|UserID2|UserID3")
}
```
