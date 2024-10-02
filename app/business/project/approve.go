package project

import (
	"encoding/json"
	"fmt"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 用于自动注册路由
type Approve struct{}

func init() {
	fpath := Approve{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

// 获取列表
func (api *Approve) Save(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	id := parameter["id"]
	note := parameter["note"]
	approveStatus := parameter["status"].([]interface{})
	data := map[string]interface{}{
		"approve_status": approveStatus[0],
		"approve_uid":    user.Name,
		"approve_note":   note,
		"approve_time":   time.Now().Unix(),
		"update_time":    time.Now().Unix(),
	}

	res, err := model.DB().Table("business_project_files").Where("id", id).Data(data).Update()
	if err != nil {
		results.Failed(c, "审批失败", err)
	} else {
		results.Success(c, "审批成功", res, nil)
	}
}

// 获取审批列表
func (api *Approve) Get_list(c *gin.Context) {
	// getuser, _ := c.Get("user")
	// user := getuser.(*middleware.UserClaims)

	approveStatus := c.DefaultQuery("status", "")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)

	MDB := model.DB().Table("business_project_files a").LeftJoin("business_account b", "a.uid", "=", "b.id")
	MDBC := model.DB().Table("business_project_files")

	if approveStatus != "" && approveStatus != "*" {
		MDB.Where("a.approve_status", approveStatus)
		MDBC.Where("approve_status", approveStatus)
	} else {
		MDB.Where("a.approve_status", ">", -2)
		MDBC.Where("approve_status", ">", -2)
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Fields("a.*,b.name").Order("a.id desc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		for _, val := range list {
			if val["name"] == nil {
				val["name"] = fmt.Sprintf("%s%s", val["user_name"], "[已删除]")
			}
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
