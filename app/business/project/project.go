package project

import (
	"encoding/json"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/results"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

// 数据导出
func (api *Project) Get_export_owner(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	downloadAll := c.DefaultQuery("downloadAll", "false")
	all, _ := strconv.ParseBool(downloadAll)
	MDB := model.DB().Table("business_project")
	MDB.Where("uid", user.ID)
	LOGMDB := model.DB().Table("media_download_log")
	ids := make([]interface{}, 0)
	if !all {
		selectedKeys := c.QueryArray("selectedKeys[]")
		for _, id := range selectedKeys {
			if id == "" {
				continue
			}
			ids = append(ids, id)
		}
		MDB.WhereIn("id", ids)
	} else {
		platform := c.DefaultQuery("platform", "")
		MDB.Where("platform", platform)
	}
	list, err := MDB.Get()
	if err != nil {
		results.BadRequest(c, err.Error(), nil)
		return
	} else {
		var buf []byte
		var err error
		if user.Rolename == "普通用户" {
			buf, err = generatePartExcelFile(list)
		} else {
			buf, err = generateAllExcelFile(list)
		}
		if err != nil {
			results.Failed(c, err.Error(), nil)
			return
		}
		filename := "files-" + time.Now().Format("20060102-150405") + ".xlsx"
		c.Header("Cache-Control", "no-cache")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
		c.Header("response-type", "blob") // 以流的形式下载必须设置这一项，否则前端下载下来的文件会出现格式不正确或已损坏的问题
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.ms-excel", buf)
		Insertdata := map[string]interface{}{
			"name":         user.Name,
			"username":     user.Username,
			"uid":          user.ID,
			"rolename":     user.Rolename,
			"download_cnt": len(list),
			"createtime":   time.Now().Unix(),
		}
		_, err = LOGMDB.Data(Insertdata).InsertGetId()
		if err != nil {
			log.Printf("export log insert failed %s", err.Error())
		}
	}
}

// 数据导出
func (api *Project) Get_export(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	downloadAll := c.DefaultQuery("downloadAll", "false")
	all, _ := strconv.ParseBool(downloadAll)
	MDB := model.DB().Table("business_project")

	LOGMDB := model.DB().Table("media_download_log")
	ids := make([]interface{}, 0)
	if !all {
		selectedKeys := c.QueryArray("selectedKeys[]")
		for _, id := range selectedKeys {
			if id == "" {
				continue
			}
			ids = append(ids, id)
		}
		MDB.WhereIn("id", ids)
	} else {
		platform := c.DefaultQuery("platform", "")
		MDB.Where("platform", platform)
	}
	if user.Rolename == "普通用户" {
		MDB.Limit(500)
	}
	list, err := MDB.Get()
	if err != nil {
		results.BadRequest(c, err.Error(), nil)
		return
	} else {
		var buf []byte
		var err error
		if user.Rolename == "普通用户" {
			buf, err = generatePartExcelFile(list)
		} else {
			buf, err = generateAllExcelFile(list)
		}
		if err != nil {
			results.Failed(c, err.Error(), nil)
			return
		}
		filename := "files-" + time.Now().Format("20060102-150405") + ".xlsx"
		c.Header("Cache-Control", "no-cache")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
		c.Header("response-type", "blob") // 以流的形式下载必须设置这一项，否则前端下载下来的文件会出现格式不正确或已损坏的问题
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.ms-excel", buf)
		Insertdata := map[string]interface{}{
			"name":         user.Name,
			"username":     user.Username,
			"uid":          user.ID,
			"rolename":     user.Rolename,
			"download_cnt": len(list),
			"createtime":   time.Now().Unix(),
		}
		_, err = LOGMDB.Data(Insertdata).InsertGetId()
		if err != nil {
			log.Printf("export log insert failed %s", err.Error())
		}
	}
}

func generateAllExcelFile(list []gform.Data) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	var err error
	ss := func(sheet string, idx int, slice interface{}) int {
		if idx == 1 {
			f.NewSheet(sheet)
		}
		err = f.SetSheetRow(sheet, "A"+strconv.Itoa(idx), slice)
		if err != nil {
			log.Printf("[error] %s \n", err.Error())
		}
		return idx + 1
	}

	idx := 1
	idx = ss("Sheet1", idx, &[]interface{}{
		"编号", "合作平台", "合作时间", "账号类型", "账号昵称", "粉丝数（万）", "发布链接接",
		"合作方式", "平台价/刊例", "执行价（含税）", "折扣说明",
		"税率", "事业部", "项目号", "项目名称", "支付单号", "联系方式",
	})

	for _, data := range list {
		cooperateTime := data["cooperate_time"]
		c := cooperateTime.(int64)
		t := time.Unix(c, 0)
		idx = ss("Sheet1", idx, &[]interface{}{
			data["id"],
			data["platform"],
			t.Format(time.DateOnly),
			data["account_type"],
			data["account_nikename"],
			data["fanscnt"],
			data["publish_link"],
			data["cooperate_type"],
			data["platform_price"],
			data["actual_price"],
			data["discount_note"],
			data["tax_rate"],
			data["department"],
			data["projectno"],
			data["project_name"],
			data["payno"],
			data["supply_name"],
			data["contact"],
		})
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func generatePartExcelFile(list []gform.Data) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	var err error
	ss := func(sheet string, idx int, slice interface{}) int {
		if idx == 1 {
			f.NewSheet(sheet)
		}
		err = f.SetSheetRow(sheet, "A"+strconv.Itoa(idx), slice)
		if err != nil {
			log.Printf("[error] %s \n", err.Error())
		}
		return idx + 1
	}

	idx := 1
	idx = ss("Sheet1", idx, &[]interface{}{
		"编号", "合作平台", "合作时间", "账号类型", "账号昵称", "粉丝数（万）", "发布链接接",
		"合作方式", "平台价/刊例", "执行价（含税）", "折扣说明", "税率",
	})

	for _, data := range list {
		cooperateTime := data["cooperate_time"]
		c := cooperateTime.(int64)
		t := time.Unix(c, 0)
		idx = ss("Sheet1", idx, &[]interface{}{
			data["id"],
			data["platform"],
			t.Format(time.DateOnly),
			data["account_type"],
			data["account_nikename"],
			data["fanscnt"],
			data["publish_link"],
			data["cooperate_type"],
			data["platform_price"],
			data["actual_price"],
			data["discount_note"],
			data["tax_rate"],
		})
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 查询导出日志
func (api *Project) Get_exportlog(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	log.Printf("Get_exportlog %s ", user.Username)
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)

	name := c.DefaultQuery("name", "")
	createdTime := c.QueryArray("createdtime[]")

	MDB := model.DB().Table("media_download_log")
	MDBC := model.DB().Table("media_download_log")

	if name != "" {
		MDB.Where("name", name)
		MDBC.Where("name", name)
	}
	if len(createdTime) == 2 {
		star_time := gf.StringTimestamp(createdTime[0]+" 00:00", "datetime")
		end_time := gf.StringTimestamp(createdTime[1]+" 23:59", "datetime")
		MDB.WhereBetween("createtime", []interface{}{star_time, end_time})
		MDBC.WhereBetween("createtime", []interface{}{star_time, end_time})
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Fields("*").Order("id asc").Get()
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
