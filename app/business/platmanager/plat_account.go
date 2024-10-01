package platmanager

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
type Account struct{}

func init() {
	fpath := Account{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *Account) Isaccounttypeexist(c *gin.Context) {
	//获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["account_type"] != nil && parameter["plat_id"] != nil {
		res1, err := model.DB().Table("media_plat_account").Where("account_type", parameter["account_type"]).Where("plat_id", parameter["plat_id"]).Value("id")
		if err != nil {
			results.Failed(c, "数据库操作失败", err)
		} else if res1 != nil {
			results.Failed(c, "类型已存在", err)
		} else {
			results.Success(c, "验证通过", res1, nil)
		}
	} else {
		results.Failed(c, "请求参数错误", nil)
	}
}

// 获取列表
func (api *Account) Get_list(c *gin.Context) {
	plat_id := c.DefaultQuery("plat_id", "")
	account_type := c.DefaultQuery("account_type", "")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := model.DB().Table("media_plat_account").Where("plat_id", plat_id)
	MDBC := model.DB().Table("media_plat_account").Where("plat_id", plat_id)
	if account_type != "" {
		MDB.Where("account_type", "like", "%"+account_type+"%")
		MDBC.Where("account_type", "like", "%"+account_type+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("order_id asc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		var totalCount int64
		totalCount, _ = MDBC.Count("*")
		results.Success(c, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}, nil)
	}
}

// 获取list
func (api *Account) Get_list_nopage(c *gin.Context) {
	plat_id := c.DefaultQuery("plat_id", "")
	plat_name := c.DefaultQuery("plat_name", "")
	MDB := model.DB().Table("media_plat_account a").LeftJoin("media_plat b", "a.plat_id", "=", "b.id")
	if plat_id == "" && plat_name == "" {
		results.Failed(c, "请求参数错误,参数为空", nil)
	}
	if plat_id != "" {
		MDB.Where("a.plat_id", plat_id)
	}
	if plat_name != "" {
		MDB.Where("b.plat_name", plat_name)
	}
	list, err := MDB.Fields("a.*,b.plat_name").Order("a.order_id asc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		results.Success(c, "获取全部列表", map[string]interface{}{
			"items": list}, nil)
	}
}

// 保存
func (api *Account) Save(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	if f_id == 0 {
		delete(parameter, "id")
		parameter["createtime"] = time.Now().Unix()
		addId, err := model.DB().Table("media_plat_account").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(c, "添加失败", err)
		} else {
			if addId != 0 {
				model.DB().Table("media_plat_account").
					Data(map[string]interface{}{"order_id": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(c, "添加成功！", addId, nil)
		}
	} else {
		res, err := model.DB().Table("media_plat_account").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(c, "更新失败", err)
		} else {
			results.Success(c, "更新成功！", res, nil)
		}
	}
}

// 删除
func (api *Account) Del(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	plat_id := parameter["plat_id"]
	res2, err := model.DB().Table("media_plat_account").Where("plat_id", plat_id).WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(c, "删除失败", err)
	} else {
		results.Success(c, "删除成功！", res2, nil)
	}
}
