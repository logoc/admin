package project

import (
	"encoding/json"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 用于自动注册路由
type UpConfig struct{}

func init() {
	fpath := UpConfig{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

// 删除记录
func (api *UpConfig) Del(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	id := parameter["id"].(float64)
	del_id, err := model.DB().Table("media_upconfig").Where("id", id).Delete()
	if err != nil {
		results.Failed(c, "删除失败", err)
	} else {
		results.Success(c, "删除成功！", del_id, nil)
	}
}

// 新增导入配置
func (api *UpConfig) Save(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)

	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	log.Printf("[info] user:%s upconfig save\n", user.Username)
	cate_type := parameter["cate_type"]
	cate := parameter["cate"]
	id := parameter["id"].(float64)
	if cate_type == "" || cate == "" {
		results.Failed(c, "参数错误", "参数错误")
		return
	}
	// 更新
	if id != 0 {
		res, err := model.DB().Table("media_upconfig").Where("id", id).Data(map[string]interface{}{"cate_type": cate_type, "cate": cate}).Update()
		if err != nil {
			results.Failed(c, "数据库访问失败", err)
		} else {
			results.Success(c, "更新成功", res, nil)
		}
		return
	}

	list, err := model.DB().Table("media_upconfig").Where("cate_type", cate_type).Where("cate", cate).Get()
	if err != nil {
		results.Failed(c, "数据库访问失败", err)
		return
	}
	if len(list) > 0 {
		results.Failed(c, "已经存在", "已经存在")
		return
	}

	Insertdata := map[string]interface{}{
		"cate":       cate,
		"cate_type":  cate_type,
		"cate_id":    0,
		"type_id":    0,
		"createtime": time.Now().Unix(),
	}
	newId, err := model.DB().Table("media_upconfig").Data(Insertdata).InsertGetId()
	if err != nil {
		results.Failed(c, "数据库访问失败", err)
		return
	}
	results.Success(c, "新增成功", newId, nil)
}

// 获取导入配置列表
func (api *UpConfig) Get_list(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	log.Printf("[info] user:%s upconfig get_list\n", user.Username)
	t := c.DefaultQuery("cate_type", "")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)

	MDB := model.DB().Table("media_upconfig")
	MDBC := model.DB().Table("media_upconfig")
	if t != "" {
		MDB.Where("cate_type", t)
		MDBC.Where("cate_type", t)
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("cate_type,cate_id asc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		totalCount, _ := MDBC.Count("*")
		results.Success(c, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}, nil)
	}
}

// 获取导入配置列表
func (api *UpConfig) Get_list_nopage(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	log.Printf("[info] user:%s upconfig Get_list_nopage\n", user.Username)
	t := c.DefaultQuery("cate_type", "")
	MDB := model.DB().Table("media_upconfig")
	if t != "" {
		MDB.Where("cate_type", t)
	}
	list, err := MDB.Order("cate_type, cate_id asc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		results.Success(c, "获取全部列表", map[string]interface{}{
			"items": list}, nil)
	}
}
