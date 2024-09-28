package project

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 用于自动注册路由
type Upfile struct{}

func init() {
	fpath := Upfile{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

// 删除记录
func (api *Upfile) Del(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res, err := model.DB().Table("business_project_files").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(c, "删除失败", err)
	} else {
		_, err := model.DB().Table("business_project").WhereIn("file_id", ids.([]interface{})).Data(map[string]interface{}{"status": -1}).Update()
		if err != nil {
			log.Printf("[error] 删除文件明细失败！ %v \n", err)
		}
		results.Success(c, "删除成功！", res, nil)
	}
}

// 获取获取上传列表
func (api *Upfile) Get_list(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	file_name := c.DefaultQuery("file_name", "")
	approve_status := c.DefaultQuery("approve_status", "")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)

	MDB := model.DB().Table("business_project_files")
	MDBC := model.DB().Table("business_project_files")

	MDB.Where("uid", user.ID)
	MDBC.Where("uid", user.ID)

	if approve_status != "" && approve_status != "*" {
		MDB.Where("approve_status", approve_status)
		MDBC.Where("approve_status", approve_status)
	}

	if file_name != "" {
		MDB.Where("file_name", "like", "%"+file_name+"%")
		MDBC.Where("file_name", "like", "%"+file_name+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		for _, val := range list {
			val["create_time"] = time.Unix(val["create_time"].(int64), 0).Format("2006-01-02 01:01:00")
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

// 上传文件
func (api *Upfile) UploadExcel(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	file, err := c.FormFile("file")
	if err != nil {
		results.Failed(c, "文件上传失败", err)
		return
	}
	//判断文件是否已经传过
	fileContent, err := file.Open()
	if err != nil {
		results.Failed(c, "文件破损，打开失败", err)
		return
	}
	projectDatas, err := parseExcel(fileContent)
	if err != nil {
		results.Failed(c, "文件解析失败", err)
		return
	}
	m_d5 := md5.New()
	if _, err := io.Copy(m_d5, fileContent); err != nil {
		results.Failed(c, "文件签名失败", err)
		return
	}
	sha1_str := hex.EncodeToString(m_d5.Sum(nil))
	nowTime := time.Now().Unix() //当前时间
	//查找该用户是否传过
	fileInfo, _ := model.DB().Table("business_project_files").Where("uid", user.ID).
		Where("sha1", sha1_str).Fields("file_name").First()
	if fileInfo != nil { //文件是否已经存在
		results.Success(c, "文件已上传", fileInfo, nil)
		return
	} else {
		file_path := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
		//如果没有filepath文件目录就创建一个
		if _, err := os.Stat(file_path); err != nil {
			if !os.IsExist(err) {
				os.MkdirAll(file_path, os.ModePerm)
			}
		}
		//上传到的路径
		filename_arr := strings.Split(file.Filename, ".")
		//重新名片-lunix系统不支持中文
		name_str := md5Str(fmt.Sprintf("%v%s", nowTime, filename_arr[0]))      //组装文件保存名字
		file_Filename := fmt.Sprintf("%s%s%s", name_str, ".", filename_arr[1]) //文件加.后缀
		path := file_path + file_Filename
		// 上传文件到指定的目录
		err = c.SaveUploadedFile(file, path)
		if err != nil { //上传失败
			c.JSON(200, gin.H{
				"uid":      sha1_str,
				"name":     file.Filename,
				"status":   "error",
				"response": "上传失败",
				"time":     nowTime,
			})
		} else { //上传成功
			//保存数据
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			Insertdata := map[string]interface{}{
				"user_name":      user.Username,
				"uid":            user.ID,
				"file_name":      file.Filename,
				"sha1":           sha1_str,
				"file_size":      file.Size,
				"url":            path,
				"project_count":  len(projectDatas),
				"storage":        dir + "/" + path,
				"create_time":    nowTime,
				"approve_status": 0,
			}
			//保存数据
			fileId, err := model.DB().Table("business_project_files").Data(Insertdata).InsertGetId()
			if err != nil {
				results.Failed(c, "数据库存储失败", err)
				return
			}
			insertProjects(fileId, projectDatas)
			//返回数据
			getdata, err := model.DB().Table("business_project_files").Where("id", fileId).Fields("id, file_name,file_size, approve_status, create_time").First()
			if err != nil {
				results.Failed(c, "查询更新失败", err)
				return
			}
			results.Success(c, "上传成功", getdata, nil)
		}
	}
}

// md5加密
func md5Str(origin string) string {
	m := md5.New()
	m.Write([]byte(origin))
	return hex.EncodeToString(m.Sum(nil))
}

func parseExcel(file io.Reader) ([]map[string]interface{}, error) {
	ex, err := excelize.OpenReader(file)
	if err != nil {
		log.Printf("[error] %v \n", err)
		return nil, err
	}
	idx := ex.GetActiveSheetIndex()
	name := ex.GetSheetName(idx)
	rows, err := ex.GetRows(name)

	if err != nil {
		log.Printf("[error] %v \n", err)
		return nil, err
	}
	insertDatas := make([]map[string]interface{}, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		for _, colCell := range row {
			fmt.Print(colCell, "|")

		}
		fmt.Print("\n")
		if len(row) < 18 {
			log.Printf("[error] sheet name[%s] 列数%d 小于 18 \n", name, len(row))
			return nil, fmt.Errorf("Excel模版错误")
		}
		data := map[string]interface{}{
			"platform":         row[1],
			"cooperate_time":   gf.StringTimestamp2(row[2]),
			"account_type":     row[3],
			"account_nikename": row[4],
			"fanscnt":          row[5],
			"publish_link":     row[6],
			"cooperate_type":   row[7],
			"platform_price":   row[8],
			"actual_price":     row[9],
			"discount_note":    row[10],
			"tax_rate":         row[11],
			"department":       row[12],
			"projectno":        row[13],
			"project_name":     row[14],
			"payno":            row[15],
			"supply_name":      row[16],
			"contact":          row[17],
			"create_time":      time.Now().Unix(),
		}
		insertDatas = append(insertDatas, data)
	}
	return insertDatas, nil
}

func insertProjects(fileId int64, datas []map[string]interface{}) (failed int) {
	for _, data := range datas {
		data["file_id"] = fileId
		_, err := model.DB().Table("business_project").Data(data).InsertGetId()
		if err != nil {
			log.Printf("[error] mysql insert error %v", err)
			failed += 1
			continue
		}
	}
	return
}
