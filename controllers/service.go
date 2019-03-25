package controllers

import (
	"fmt"
	"github.com/astaxie/beego/context"
	wechatApi "github.com/silenceper/wechat"
	//"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/message"
	"gowe/models"
)

func Service(ctx *context.Context) {
	//cache.NewRedis(&cache.RedisOpts{Host:"127.0.0.1:6379"})
	//
	flag := ctx.Input.Query(":flag")
	//wechatConfig := cache.Redis{}.Get(flag)

	wechatConfig := models.Wechat{Flag:flag}.Find()
	if len(wechatConfig) <= 0 {
		return
	}
	//配置微信参数
	config := &wechatApi.Config{
		AppID:          wechatConfig[0].Appid,
		AppSecret:      wechatConfig[0].Appsecret,
		Token:          wechatConfig[0].Token,
		EncodingAESKey: wechatConfig[0].EncodingAesKey,
	}
	wc := wechatApi.NewWechat(config)

	// 传入request和responseWriter
	server := wc.GetServer(ctx.Request, ctx.ResponseWriter)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {

		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}
