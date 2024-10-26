package project

import (
	"encoding/json"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 用于自动注册路由
type Project struct{}

func init() {
	fpath := Project{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

// 获取列表
func (api *Project) Get_list(c *gin.Context) {
	fileId := c.DefaultQuery("id", "0")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := model.DB().Table("business_project")
	MDBC := model.DB().Table("business_project")
	MDB.Where("file_id", fileId)
	MDBC.Where("file_id", fileId)

	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		for _, val := range list {
			val["cooperate_time"] = time.Unix(val["cooperate_time"].(int64), 0).Format("2006-01-02")
		}
		var totalCount int64
		totalCount, _ = MDBC.Count("*")
		results.Success(c, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}, nil)
	}
}

func (api *Project) Get_search(c *gin.Context) {
	// getuser, _ := c.Get("user")
	// user := getuser.(*middleware.UserClaims)

	account_nikename := c.DefaultQuery("accountNikeName", "")
	projectno := c.DefaultQuery("projectNo", "")
	cooperate_time := c.QueryArray("cooperateTime[]")

	supplyName := c.DefaultQuery("supplyName", "")
	platform := c.DefaultQuery("platform", "")
	fansCntMin := c.DefaultQuery("fansCntMin", "")
	fansCntMax := c.DefaultQuery("fansCntMax", "")
	priceMin := c.DefaultQuery("priceRangeMin", "")
	priceMax := c.DefaultQuery("priceRangeMax", "")
	accountType := c.QueryArray("accountType[]")
	// downloadAll := c.DefaultQuery("downloadAll", "false")

	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := model.DB().Table("business_project a").LeftJoin("business_project_files b", "a.file_id", "=", "b.id")
	MDBC := model.DB().Table("business_project a").LeftJoin("business_project_files b", "a.file_id", "=", "b.id")

	MDB.Where("b.approve_status", 1)
	MDBC.Where("b.approve_status", 1)

	if platform != "" && platform != "不限" {
		MDB.Where("a.platform", platform)
		MDBC.Where("a.platform", platform)
	}
	if supplyName != "" {
		MDB.Where("a.supply_name", supplyName)
		MDBC.Where("a.supply_name", supplyName)
	}
	if fansCntMin != "" && fansCntMax != "" {
		MDB.WhereBetween("a.fanscnt", []interface{}{fansCntMin, fansCntMax})
		MDBC.WhereBetween("a.fanscnt", []interface{}{fansCntMin, fansCntMax})
	}
	if priceMin != "" && priceMax != "" {
		MDB.WhereBetween("a.platform_price", []interface{}{priceMin, priceMax})
		MDBC.WhereBetween("a.platform_price", []interface{}{priceMin, priceMax})
	}
	if len(accountType) > 0 && accountType[0] != "不限" {
		account := make([]interface{}, 0)
		for _, a := range accountType {
			account = append(account, a)
		}
		MDB.WhereIn("a.account_type", account)
		MDBC.WhereIn("a.account_type", account)
	}

	if account_nikename != "" {
		MDB.Where("a.account_nikename", "like", "%"+account_nikename+"%")
		MDBC.Where("a.account_nikename", "like", "%"+account_nikename+"%")
	}
	if projectno != "" {
		MDB.Where("a.projectno", "like", "%"+projectno+"%")
		MDBC.Where("a.projectno", "like", "%"+projectno+"%")
	}
	if len(cooperate_time) > 0 {
		star_time := gf.StringTimestamp(cooperate_time[0]+" 00:00", "datetime")
		end_time := gf.StringTimestamp(cooperate_time[1]+" 23:59", "datetime")
		MDB.WhereBetween("a.cooperate_time", []interface{}{star_time, end_time})
		MDBC.WhereBetween("a.cooperate_time", []interface{}{star_time, end_time})
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Fields("a.*,b.approve_status, b.uid, b.user_name").Order("a.id desc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		// for _, val := range list {
		// 	val["a.cooperate_time"] = time.Unix(val["a.cooperate_time"].(int64), 0).Format("2006-01-02")
		// }
		var totalCount int64
		totalCount, _ = MDBC.Count("*")
		results.Success(c, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}, nil)
	}
}

// 更新发文链接
func (api *Project) Update_publishlink(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["id"] == nil || parameter["publish_link"] == nil {
		results.Failed(c, "参数错误", nil)
		return
	}
	data := map[string]interface{}{
		"publish_link": parameter["publish_link"],
	}
	res, err := model.DB().Table("business_project").
		Data(data).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(c, "更新失败", err)
	} else {
		results.Success(c, "更新成功！", res, nil)
	}
}

func (api *Project) Del(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["ids"] == nil {
		results.Failed(c, "参数错误", nil)
		return
	}
	ids := parameter["ids"].([]interface{})
	del_id, err := model.DB().Table("business_project").WhereIn("id", ids).Delete()
	if err != nil {
		results.Failed(c, "删除失败", err)
	} else {
		results.Success(c, "删除成功！", del_id, nil)
	}
}
