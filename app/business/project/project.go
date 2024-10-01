package project

import (
	"encoding/json"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"strconv"
	"strings"
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

func (api *Project) Get_search_owner(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)

	account_nikename := c.DefaultQuery("accountNikeName", "")
	projectno := c.DefaultQuery("projectNo", "")
	cooperate_time := c.DefaultQuery("cooperateTime", "")

	platform := c.DefaultQuery("platform", "")
	fansCntStr := c.DefaultQuery("fansCnt", "")
	fansCntRange := strings.Split(fansCntStr, "-")
	priceStr := c.DefaultQuery("priceRange", "")
	priceRange := strings.Split(priceStr, "-")
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

	MDB.Where("b.uid", user.ID)
	MDBC.Where("b.uid", user.ID)

	if platform != "" && platform != "不限" {
		MDB.Where("a.platform", platform)
		MDBC.Where("a.platform", platform)
	}
	if len(fansCntRange) == 2 {
		MDB.WhereBetween("a.fanscnt", []interface{}{fansCntRange[0], fansCntRange[1]})
		MDBC.WhereBetween("a.fanscnt", []interface{}{fansCntRange[0], fansCntRange[1]})
	}

	if len(priceRange) == 2 {
		MDB.WhereBetween("a.platform_price", []interface{}{priceRange[0], priceRange[1]})
		MDBC.WhereBetween("a.platform_price", []interface{}{priceRange[0], priceRange[1]})
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
		MDB.Where("a.projectno", projectno)
		MDBC.Where("a.projectno", projectno)
	}
	if cooperate_time != "" && cooperate_time != "0" {
		t, _ := strconv.ParseInt(cooperate_time, 10, 64)
		tm := gf.NowBeforeTimestamp(t)
		MDB.Where("a.cooperate_time", ">=", tm)
		MDBC.Where("a.cooperate_time", ">=", tm)
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

func (api *Project) Get_search(c *gin.Context) {
	// getuser, _ := c.Get("user")
	// user := getuser.(*middleware.UserClaims)

	account_nikename := c.DefaultQuery("accountNikeName", "")
	projectno := c.DefaultQuery("projectNo", "")
	cooperate_time := c.DefaultQuery("cooperateTime", "")

	platform := c.DefaultQuery("platform", "")
	fansCntStr := c.DefaultQuery("fansCnt", "")
	fansCntRange := strings.Split(fansCntStr, "-")
	priceStr := c.DefaultQuery("priceRange", "")
	priceRange := strings.Split(priceStr, "-")
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
	if len(fansCntRange) == 2 {
		MDB.WhereBetween("a.fanscnt", []interface{}{fansCntRange[0], fansCntRange[1]})
		MDBC.WhereBetween("a.fanscnt", []interface{}{fansCntRange[0], fansCntRange[1]})
	}

	if len(priceRange) == 2 {
		MDB.WhereBetween("a.platform_price", []interface{}{priceRange[0], priceRange[1]})
		MDBC.WhereBetween("a.platform_price", []interface{}{priceRange[0], priceRange[1]})
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
		MDB.Where("a.projectno", projectno)
		MDBC.Where("a.projectno", projectno)
	}
	if cooperate_time != "" && cooperate_time != "0" {
		t, _ := strconv.ParseInt(cooperate_time, 10, 64)
		tm := gf.NowBeforeTimestamp(t)
		MDB.Where("a.cooperate_time", ">=", tm)
		MDBC.Where("a.cooperate_time", ">=", tm)
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

// 保存
func (api *Project) Save(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	delete(parameter, "catename")
	parameter["createtime"] = time.Now().Unix()
	res, err := model.DB().Table("common_picture").
		Data(parameter).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(c, "更新失败", err)
	} else {
		results.Success(c, "更新成功！", res, nil)
	}
}
