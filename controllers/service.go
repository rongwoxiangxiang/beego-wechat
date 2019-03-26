package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
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
	flag := ctx.Input.Query(":flag")
	res := responseText(flag,wechatConfig)
	fmt.Println(res)
	//wechatConfig := config(ctx)
	//server := wechatApi.NewWechat(&wechatApi.Config{
	//	AppID:          wechatConfig["Appid"].(string),
	//	AppSecret:      wechatConfig["Appsecret"].(string),
	//	Token:          wechatConfig["Token"].(string),
	//	EncodingAESKey: wechatConfig["EncodingAesKey"].(string),
	//	Cache:			redis,
	//}).GetServer(ctx.Request, ctx.ResponseWriter)
	//server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
	//	var msgType message.MsgType
	//	switch msg.MsgType {
	//	case message.MsgTypeText:
	//		return responseText(msg.Content,wechatConfig)
	//	case message.MsgTypeEvent:
	//		return responseEvent(msg.EventKey,wechatConfig)
	//	default:
	//		return &message.Reply{MsgType: msgType, MsgData: message.NewText(message.MsgTypeVoice)}
	//	}
	//})
	//
	////处理消息接收以及回复
	//err := server.Serve()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	////发送回复的消息
	//server.Send()
}

func responseText(msg string,conf map[string]interface{}) *message.Reply {
	reply := models.Reply{Wid:int64(conf["Id"].(float64)), Alias:msg}.FindOne()
	return replyActivity(reply)
}

func responseEvent(msg string,conf map[string]interface{}) *message.Reply {
	reply := models.Reply{Wid:int64(conf["Id"].(float64)), ClickKey:msg}.FindOne()
	return replyActivity(reply)
}

func replyActivity(reply models.Reply) *message.Reply {
	if reply.Id > 0 {
		switch reply.Type {
		case models.REPLY_TYPE_TEXT:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(reply.Success)}
		case models.REPLY_TYPE_CODE:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(reply.Success)}
		case models.REPLY_TYPE_LUCKY:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(reply.Success)}
		case models.REPLY_TYPE_CHECKIN:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(reply.Success)}
		default:
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(reply.Success)}
		}
	}
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(reply.Success)}
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