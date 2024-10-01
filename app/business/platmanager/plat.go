package platmanager

import (
	"encoding/json"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/results"
	"io"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

// 用于自动注册路由
type Plat struct{}

func init() {
	fpath := Plat{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *Plat) Isplatexist(c *gin.Context) {
	//获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["plat_name"] != nil {
		res1, err := model.DB().Table("media_plat").Where("plat_name", parameter["plat_name"]).Value("id")
		if err != nil {
			results.Failed(c, "数据库操作失败", err)
		} else if res1 != nil {
			results.Failed(c, "平台名称已存在", err)
		} else {
			results.Success(c, "验证通过", res1, nil)
		}
	} else {
		results.Failed(c, "请求参数错误", nil)
	}
}

// 获取审批列表
func (api *Plat) Get_list(c *gin.Context) {
	list, _ := model.DB().Table("media_plat").Fields("id,plat_name,remark,order_id,createtime").Order("order_id asc").Get()
	if list == nil {
		list = make([]gform.Data, 0)
	}
	results.Success(c, "获取列表", list, nil)
}

// 保存
func (api *Plat) Save(c *gin.Context) {
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
		addId, err := model.DB().Table("media_plat").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(c, "添加失败", err)
		} else {
			if addId != 0 {
				model.DB().Table("media_plat").
					Data(map[string]interface{}{"order_id": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(c, "添加成功！", addId, nil)
		}
	} else {
		res, err := model.DB().Table("media_plat").
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
func (api *Plat) Del(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := model.DB().Table("media_plat").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(c, "删除失败", err)
	} else {
		results.Success(c, "删除成功！", res2, nil)
	}
}
