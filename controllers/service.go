package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"gowe/common"
	"strconv"
	"strings"
	wechatApi "github.com/silenceper/wechat"
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
	wechaBase := wechatApi.NewWechat(&wechatApi.Config{
		AppID:          wechatConfig["Appid"].(string),
		AppSecret:      wechatConfig["Appsecret"].(string),
		Token:          wechatConfig["Token"].(string),
		EncodingAESKey: wechatConfig["EncodingAesKey"].(string),
		Cache:			redis,
	})
	server := wechaBase.GetServer(ctx.Request, ctx.ResponseWriter)
	wxUser := getWechatUser(server.GetOpenID(), int64(wechatConfig["Id"].(float64)), wechaBase)
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		return responseEventText(msg, wxUser)
	})
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	server.Send()
}

func responseEventText(msg message.MixMessage, user models.WechatUser) (msgReply *message.Reply) {
	var reply models.Reply
	if msg.MsgType == message.MsgTypeText {
		reply = models.Reply{Wid:user.Wid, Alias:msg.Content}.FindOne()
	} else if msg.MsgType == message.MsgTypeEvent {
		reply = models.Reply{Wid:user.Wid, ClickKey:msg.EventKey}.FindOne()
	} else {
		reply = models.Reply{Wid:user.Wid, Alias:msg.EventKey}.FindOne()
	}
	if reply.Id == 0 {
		return
	}
	switch reply.Type {
	case models.REPLY_TYPE_TEXT:
		msgReply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(reply.Success),
		}
	case models.REPLY_TYPE_CODE:
		msgReply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(doReplyCode(reply, user)),
		}
	case models.REPLY_TYPE_LUCKY:
		msgReply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(doReplyLuck(reply, user)),
		}
	case models.REPLY_TYPE_CHECKIN:
		msgReply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(doReplyCheckin(reply, user)),
		}
	}
	return
}

func doReplyCode(reply models.Reply, wechatUser models.WechatUser) string {
	history := models.PrizeHistory{ActivityId:reply.ActivityId,Wuid:wechatUser.Id}.GetByActivityWuId()
	if len(history) > 0 {
		return strings.Replace(reply.Success, "%prize%", history[0].Prize, 1)
	}
	prize, err := models.Prize{ActivityId:reply.ActivityId, Level:int8(models.PRIZE_LEVEL_DEFAULT), Used:common.NO_VALUE}.FindOneUsedCode()
	if err == common.ErrDataUnExist {
		return reply.Fail
	}
	if prize.Code != "" {
		_, err = models.PrizeHistory{ActivityId:reply.ActivityId,Wuid:wechatUser.Id,Prize:prize.Code}.Insert()
		if err != nil {
			return reply.Fail
		}
		return strings.Replace(reply.Success, "%prize%", prize.Code, 1)
	}
	return models.PLEASE_TRY_AGAIN
}

func doReplyLuck(reply models.Reply, wechatUser models.WechatUser) string {
	history := models.PrizeHistory{ActivityId:reply.ActivityId,Wuid:wechatUser.Id}.GetByActivityWuId()
	if len(history) > 0 {
		return strings.Replace(reply.Success, "%prize%", history[0].Prize, 1)
	}
	//TODO，多次参与

	luck, err := models.Lottery{Wid:reply.Wid, ActivityId:reply.ActivityId}.Luck()
	if err == common.ErrLuckFinal {
		return common.ErrLuckFinal.Msg
	}
	if err != nil {
		return common.ErrLuckFail.Msg
	}
	if luck.Name != "" {
		_, err = models.PrizeHistory{ActivityId:reply.ActivityId,Wuid:wechatUser.Id,Prize:luck.Name,Level:luck.Level}.Insert()
	}
	return strings.Replace(reply.Success, "%prize%", luck.Name, 1)
}

func doReplyCheckin(reply models.Reply, wechatUser models.WechatUser) string {
	checkin := models.Checkin{ActivityId:reply.ActivityId,Wuid:wechatUser.Id,Wid:wechatUser.Wid}.GetCheckinByActivityWuid()
	if checkin.Id == 0 {
		return models.CHECK_FAIL
	}
	lastCheckinDate := checkin.Lastcheckin.Format("2006-01-02")
	if lastCheckinDate == time.Now().Format("2006-01-02") {
		return strings.
			NewReplacer("%liner%",  strconv.FormatInt(checkin.Liner, 10), "%total%", strconv.FormatInt(checkin.Total, 10)).
			Replace(reply.Success)
	}
	if lastCheckinDate == time.Now().Add(-24 * time.Hour).Format("2006-01-02"){//连续签到
		checkin.Liner = checkin.Liner + 1
	}
	checkin.Total = checkin.Total + 1
	checkin.Lastcheckin = time.Now()
	_, err := checkin.Update()
	if err != nil {
		return models.CHECK_FAIL
	}
	return strings.
		NewReplacer("%liner%", strconv.FormatInt(checkin.Liner, 10), "%total%", strconv.FormatInt(checkin.Total, 10)).
		Replace(reply.Success)
}

func getWechatUser(openId string, wid int64, wechatApi *wechatApi.Wechat) (user models.WechatUser) {
	user.Wid = wid
	user.Openid = openId
	user = user.GetByOpenid()
	go func(user models.WechatUser) {
		if user.Openid != "" && user.UpdatedAt.After(time.Now().Add(-25 * time.Hour)) {
			userInfo, err := wechatApi.GetUser().GetUserInfo(user.Openid)
			if err == nil {
				user.Nickname = userInfo.Nickname
				user.Sex = userInfo.Sex
				user.Province = userInfo.Province
				user.City = userInfo.City
				user.Country = userInfo.Country
				user.Language = userInfo.Language
				user.Headimgurl = userInfo.Headimgurl
				user.Update()
			}
		}
	}(user)
	return
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