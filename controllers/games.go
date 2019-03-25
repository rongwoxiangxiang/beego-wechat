package controllers

import (
	"encoding/json"
	"gowe/common"
	"gowe/models"
)

type GameController struct {
	BaseController
}

/**
 * @List 游戏列表
 * @Description 游戏列表
 * @Param page int
 * @Param size int
 * @Success 200 models.games
 * @router /game [get]
 */
func (this *GameController) List(){
	page, _  := this.GetInt("page",1)
	size, _  := this.GetInt("size",20)
	games := models.Game{}.LimitList(size, page)
	if len(games) <= 0 {
		games = []models.Game{}
	}
	this.Response(common.SUCCESS, games)
}

/**
 * @List 游戏详情
 * @Description 游戏详情
 * @Param id int
 * @Success 200 models.games
 * @router /game/:id [get]
 */
func (this *GameController) View(){
	id, _ := this.GetInt64(":id")
	game, err := models.Game{Id:id}.Get()
	if err != nil {
		this.Response(common.ErrDataUnExist.Code,common.ErrDataUnExist.Msg)
		return
	}
	this.Response(common.SUCCESS, game)
}

/**
 * @List 删除游戏
 * @Description 删除游戏
 * @Param id int
 * @Success 200 success
 * @router /game/:id [delete]
 */
func (this *GameController) DeleteGame() {
	id, _  := this.GetInt64(":id")
	game   := models.Game{Id:id}
	_, err := game.Get()
	if err != nil {
		this.Response(common.ErrDataUnExist.Code,common.ErrDataUnExist.Msg)
		return
	}
	result := models.Game.DeleteById(game)
	if result != true {
		this.Response(common.ErrDataUpdate.Code,common.ErrDataUpdate.Msg)
		return
	}
	this.Response(common.Success.Code, common.Success.Msg)
}

/**
 * @List 更新游戏
 * @Description 更新游戏
 * @Param id int
 * @Param params models.games
 * @Success 200 {int} models.games
 * @router /game/:id [put]
 */
func (this *GameController) UpdateGame() {
	id, _  := this.GetInt64(":id")
	game, err := models.Game{Id:id}.Get()
	if err != nil {
		this.Response(common.ErrClientParams.Code, common.ErrClientParams.Msg)
	}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &game)
	if err != nil {
		this.Response(common.ErrClientParams.Code, common.ErrClientParams.Msg)
	}
	game, err = game.Update()
	if err != nil {
		this.Response(common.ErrDataCreate.Code, common.ErrDataCreate.Msg)
	}
	this.Response(common.Success.Code, game)
	return
}

/**
 * @List 插入游戏
 * @Description 插入游戏
 * @Param params models.games
 * @Success 200 {int} models.games
 * @router /game/:id [post]
 */
func (this *GameController) InsertGame() {
	var g models.Game
	resp := this.Ctx.Input.RequestBody
	errJson := json.Unmarshal(resp, &g)
	if errJson != nil {
		this.Response(common.ErrClientParams.Code, common.ErrClientParams.Msg)
	}
	game, err := g.Insert()
	if err != nil {
		this.Response(common.ErrDataCreate.Code, common.ErrDataCreate.Msg)
	}
	this.Response(common.Success.Code, game)
	return
}