package controllers

import (
	"encoding/json"
	"gowe/common"
	"gowe/models"
)

var wechat models.Wechat

type WechatController struct {
	BaseController
}

/**
 * @List 公众号列表
 * @Description 公众号列表
 * @Param page int
 * @Param size int
 * @Success 200 models.wechat
 * @router /wechat [get]
 */
func (this *WechatController) List(){
	page, _  := this.GetInt("page",1)
	size, _  := this.GetInt("size",20)
	wechats  := wechat.LimitList(size, page)
	if len(wechats) <= 0 {
		wechats = []models.Wechat{}
	}
	this.Response(common.SUCCESS, wechats)
}

/**
 * @List 公众号详情
 * @Description 公众号详情
 * @Param id int
 * @Success 200 models.wechats
 * @router /wechat/:id [get]
 */
func (this *WechatController) View(){
	wechat.Id, _ = this.GetInt64(":id")
	wechat, err := wechat.Get()
	if err != nil {
		this.Response(common.ErrDataUnExist.Code,common.ErrDataUnExist.Msg)
		return
	}
	this.Response(common.SUCCESS, wechat)
}

/**
 * @List 删除公众号
 * @Description 删除公众号
 * @Param id int
 * @Success 200 success
 * @router /wechat/:id [delete]
 */
func (this *WechatController) DeleteWechat() {
	wechat.Id, _ = this.GetInt64(":id")
	wechat, err := wechat.Get()
	if err != nil {
		this.Response(common.ErrDataUnExist.Code,common.ErrDataUnExist.Msg)
		return
	}
	result := models.Wechat.DeleteById(wechat)
	if result != true {
		this.Response(common.ErrDataUpdate.Code,common.ErrDataUpdate.Msg)
		return
	}
	this.Response(common.Success.Code, common.Success.Msg)
}

/**
 * @List 更新公众号
 * @Description 更新公众号
 * @Param id int
 * @Param params models.wechats
 * @Success 200 {int} models.wechats
 * @router /wechat/:id [put]
 */
func (this *WechatController) UpdateWechat() {
	wechat.Id, _ = this.GetInt64(":id")
	wechat, err := wechat.Get()
	if err != nil {
		this.Response(common.ErrClientParams.Code, common.ErrClientParams.Msg)
	}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &wechat)
	if err != nil {
		this.Response(common.ErrClientParams.Code, common.ErrClientParams.Msg)
	}
	wechat, err = wechat.Update()
	if err != nil {
		this.Response(common.ErrDataCreate.Code, common.ErrDataCreate.Msg)
	}
	this.Response(common.Success.Code, wechat)
	return
}

/**
 * @List 新增公众号
 * @Description 新增公众号
 * @Param params models.wechats
 * @Success 200 {int} models.wechats
 * @router /wechat/:id [post]
 */
func (this *WechatController) InsertWechat() {
	resp := this.Ctx.Input.RequestBody
	errJson := json.Unmarshal(resp, &wechat)
	if errJson != nil {
		this.Response(common.ErrClientParams.Code, common.ErrClientParams.Msg)
	}
	wechat, err := wechat.Insert()
	if err != nil {
		this.Response(common.ErrDataCreate.Code, common.ErrDataCreate.Msg)
	}
	this.Response(common.Success.Code, wechat)
	return
}