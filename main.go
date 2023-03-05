package main

import (
	"fmt"
	"log"

	"github.com/chi2nagisa/onebot-plus-plugin-chatbot/db"
	"github.com/dezhishen/onebot-plus/pkg/plugin"
	"github.com/dezhishen/onebot-sdk/pkg/api"
	"github.com/dezhishen/onebot-sdk/pkg/model"
)

func getLoginUId(onebotApi api.OnebotApiClientInterface) string {
	var uId string
	db.Get("loginInfo", &uId)
	if uId != "" {
		return uId
	}
	loginInfo, err := onebotApi.GetLoginInfo()
	if err != nil {
		log.Println(err)
	}
	uId = fmt.Sprintf("%v", loginInfo.Data.UserId)
	db.Set("loginInfo", uId)
	return uId
}

func main() {
	// 初始化存储
	db.Init("./data")
	plugin.OnebotPluginBuilder().
		Id("chatbot").
		Name("chatbot").
		Help("@本账号 .chatbot 开启聊天机器人\n之后@本账号开始对话模式").
		Init(func(cli api.OnebotApiClientInterface) error {
			log.Println("init")
			return nil
		}).
		BeforeExit(func(cli api.OnebotApiClientInterface) error {
			log.Println("before exit")
			return nil
		}).
		HandleMessageGroup(func(data *model.EventMessageGroup, onebotApi api.OnebotApiClientInterface) error {
			loginUId := getLoginUId(onebotApi)
			hasAt := false
			for _, msg := range data.Message {
				if msg.Type == "at" && msg.Data.(*model.MessageElementAt).Qq == loginUId {
					hasAt = true
					break
				}
			}
			if !hasAt {
				return nil
			}
			chatFlagKey := fmt.Sprintf("chatFlag:%v:%v", data.GroupId, data.Sender.UserId)
			var chatFlagValue bool
			db.Get(chatFlagKey, &chatFlagValue)
			if !chatFlagValue {
				onebotApi.SendGroupMsg(
					&model.GroupMsg{
						GroupId: data.GroupId,
						Message: []*model.MessageSegment{
							{Type: "text", Data: &model.MessageElementText{Text: "你好，我是聊天机器人\n @我并且输入.chatbot\n开启聊天模式"}},
						},
					},
				)
				return nil
			}
			// todo call chatgpt
			return nil
		}).
		Build().
		Start()
}
