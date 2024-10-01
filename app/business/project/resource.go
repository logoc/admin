package project

import (
	"gofly/utils/gf"
	"gofly/utils/results"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Resource struct {
}

func init() {
	gf.Register(&Resource{}, reflect.TypeOf(Resource{}).PkgPath())
}

func (api *Resource) Get_excelfile(c *gin.Context) {
	file_path := c.DefaultQuery("file_path", "")
	if file_path == "" {
		results.BadRequest(c, "请求参数错误", nil)
		return
	}
	if _, err := os.Stat(file_path); err != nil {
		results.Failed(c, "文件不存在", nil)
		return
	}
	content := gf.ReaderFileBystring(file_path)
	c.Header("Cache-Control", "no-cache")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("response-type", "blob") // 以流的形式下载必须设置这一项，否则前端下载下来的文件会出现格式不正确或已损坏的问题
	c.Header("Content-Disposition", "attachment; filename=test.xls")
	c.Data(http.StatusOK, "application/vnd.ms-excel", []byte(content))

}
