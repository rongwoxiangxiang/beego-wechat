package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	wechatApi "github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/message"
	"gowe/models"
	"time"
)

var redis *cache.Redis

func init()  {
	opts := &cache.RedisOpts{
		Host: beego.AppConfig.String("redis_host"),
	}
	redis = cache.NewRedis(opts)
}

func Service(ctx *context.Context) {
	wechatConfig := config(ctx)

	//配置微信参数
	config := &wechatApi.Config{
		AppID:          wechatConfig["Appid"].(string),
		AppSecret:      wechatConfig["Appsecret"].(string),
		Token:          wechatConfig["Token"].(string),
		EncodingAESKey: wechatConfig["EncodingAesKey"].(string),
		Cache:			redis,
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

func config(ctx *context.Context) map[string]interface{} {
	var mp map[string]interface{}
	flag := ctx.Input.Query(":flag")
	wechatConfig := redis.Get(flag)
	if wechatConfig != nil {
		return wechatConfig.(map[string]interface{})
	}
	wechatStruct  := models.Wechat{Flag:flag}.Find()
	wechatJson, _ := json.Marshal(wechatStruct[0])
	json.Unmarshal([]byte(wechatJson), &mp)
	if err := redis.Set(flag, mp, 10 * time.Hour); err != nil {
		fmt.Println("cache: set wechat config error", err)
	}
	return mp
}