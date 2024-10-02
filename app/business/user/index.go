package user

import (
	"encoding/json"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var Client *dysmsapi.Client

/**
*使用 Index 是省略路径中的index
*本路径为： /admin/user/login -省去了index
 */
func init() {
	gf.Register(&Index{}, reflect.TypeOf(Index{}).PkgPath())
	// Client, _ := CreateDysmsapiClient()

}

// 使用AK&SK初始化账号Client
func CreateDysmsapiClient() (_result *dysmsapi.Client, _err error) {
	id := ""
	secret := ""
	config := &openapi.Config{
		AccessKeyId:     &id,
		AccessKeySecret: &secret,
	}
	endPoint := "dysmsapi.aliyuncs.com"
	config.Endpoint = &endPoint
	_result = &dysmsapi.Client{}
	_result, _err = dysmsapi.NewClient(config)
	return _result, _err
}

type Index struct {
}

/**
*1.《登录》
 */
func (api *Index) Login(c *gin.Context) {
	//获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["username"] == nil || parameter["password"] == nil {
		results.Failed(c, "请提交用户账号或密码！", nil)
		return
	}
	username := parameter["username"].(string)
	password := parameter["password"].(string)
	res, err := model.DB().Table("business_account").Fields("id,accountID,businessID,password,salt,name").Where("username", username).OrWhere("email", username).First()
	if res == nil || err != nil {
		results.Failed(c, "账号不存在！", nil)
		return
	}
	roleid, _ := model.DB().Table("business_auth_role_access").Where("uid", res["id"]).Pluck("role_id")
	role, _ := model.DB().Table("business_auth_role").WhereIn("id", roleid.([]interface{})).Fields("name").First()
	pass := gf.Md5(password + res["salt"].(string))
	if pass != res["password"] {
		results.Failed(c, "您输入的密码不正确！", pass)
		return
	}
	//token
	token := middleware.GenerateToken(&middleware.UserClaims{
		ID:             res["id"].(int64),
		Accountid:      res["accountID"].(int64),
		BusinessID:     res["businessID"].(int64),
		Name:           res["name"].(string),
		Username:       username,
		Rolename:       role["name"].(string),
		StandardClaims: jwt.StandardClaims{},
	})
	model.DB().Table("business_account").Where("id", res["id"]).Data(map[string]interface{}{"loginstatus": 1, "lastLoginTime": time.Now().Unix(), "lastLoginIp": gf.GetIp(c)}).Update()
	//登录日志
	model.DB().Table("login_logs").
		Data(map[string]interface{}{"type": 1, "uid": res["id"], "out_in": "in", "login_type": "password",
			"createtime": time.Now().Unix(), "loginIP": gf.GetIp(c)}).Insert()
	results.Success(c, "登录成功返回token！", token, nil)
}

/**
*1.《登录》
 */
func (api *Index) Login_msg(c *gin.Context) {
	//获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["mobile"] == nil || parameter["code"] == nil {
		results.Failed(c, "请提交手机号或密码！", nil)
		return
	}
	mobile := parameter["mobile"].(string)
	code := parameter["code"].(string)
	res, err := model.DB().Table("common_verify_code").Fields("*").Where("keyname", mobile).Where("code", code).First()
	if res == nil || err != nil {
		results.Failed(c, "验证失败，验证码错误", nil)
		return
	}
	// if res["createtime"].(int64)+600 < time.Now().Unix() {
	// 	results.Failed(c, "验证码已过期", nil)
	// 	return
	// }

	res, err = model.DB().Table("business_account").Fields("id,accountID,businessID,password,salt,name,username").Where("username", mobile).OrWhere("mobile", mobile).First()
	if res == nil || err != nil {
		results.Failed(c, "账号不存在！", nil)
		return
	}
	roleid, _ := model.DB().Table("business_auth_role_access").Where("uid", res["id"]).Pluck("role_id")
	role, _ := model.DB().Table("business_auth_role").WhereIn("id", roleid.([]interface{})).Fields("name").First()
	//token
	token := middleware.GenerateToken(&middleware.UserClaims{
		ID:             res["id"].(int64),
		Accountid:      res["accountID"].(int64),
		BusinessID:     res["businessID"].(int64),
		Name:           res["name"].(string),
		Username:       res["username"].(string),
		Rolename:       role["name"].(string),
		StandardClaims: jwt.StandardClaims{},
	})
	model.DB().Table("business_account").Where("id", res["id"]).Data(map[string]interface{}{"loginstatus": 1, "lastLoginTime": time.Now().Unix(), "lastLoginIp": gf.GetIp(c)}).Update()
	//登录日志
	model.DB().Table("login_logs").
		Data(map[string]interface{}{"type": 1, "uid": res["id"], "out_in": "in", "login_type": "mobile",
			"createtime": time.Now().Unix(), "loginIP": gf.GetIp(c)}).Insert()
	results.Success(c, "登录成功返回token！", token, nil)
}

/**
*2.注册
 */
func (api *Index) RegisterUser(c *gin.Context) {
	//获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["username"] == nil || parameter["password"] == nil {
		results.Failed(c, "请提交用户账号或密码！", nil)
		return
	}
	password := parameter["password"].(string)
	userdata, _ := model.DB().Table("business_account").Fields("id").Where("username", parameter["username"]).First()
	if userdata != nil {
		results.Failed(c, "账号已存在！", nil)
		return
	}
	userdata2, _ := model.DB().Table("business_account").Fields("id").Where("email", parameter["email"]).First()
	if userdata2 != nil {
		results.Failed(c, "邮箱已存在！", nil)
		return
	}
	rnd := rand.New(rand.NewSource(6))
	salt := strconv.Itoa(rnd.Int())
	pass := gf.Md5(password + salt)
	userid, err := model.DB().Table("business_account").Data(map[string]interface{}{
		"salt":       salt,
		"username":   parameter["username"],
		"password":   pass,
		"email":      parameter["email"],
		"avatar":     "resource/staticfile/avatar.png",
		"createtime": time.Now().Unix(),
	}).InsertGetId()
	if err != nil {
		results.Failed(c, "注册失败", err)
	} else {
		results.Success(c, "注册成功", userid, nil)
	}
}

/**
* 3.《获取用户》
 */
func (api *Index) Get_userinfo(c *gin.Context) {
	getuser, _ := c.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*middleware.UserClaims)
	userdata, err := model.DB().Table("business_account").Fields("id,businessID,name,username,avatar").Where("id", user.ID).First()
	if err != nil {
		results.Failed(c, "查找用户数据！", err)
	} else {
		rooturl, _ := model.DB().Table("common_config").Where("keyname", "rooturl").Value("keyvalue")
		if userdata["avatar"] == nil {
			userdata["avatar"] = rooturl.(string) + "resource/staticfile/avatar.png"
		} else if !strings.Contains(userdata["avatar"].(string), "http") && rooturl != nil {
			userdata["avatar"] = rooturl.(string) + userdata["avatar"].(string)
		}
		results.Success(c, "获取用户信息", map[string]interface{}{
			"userId":       userdata["id"],
			"businessID":   userdata["businessID"],
			"username":     userdata["username"],
			"name":         userdata["name"],
			"avatar":       userdata["avatar"],
			"introduction": userdata["remark"],
			"rooturl":      rooturl, //图片
			// "role":         "admin", //权限
		}, nil)
	}
}

/**
* 4 刷新token
 */
func (api *Index) Refreshtoken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	newtoken := middleware.Refresh(token)
	results.Success(c, "刷新token", newtoken, nil)
}

/**
*  5退出登录
 */
func (api *Index) Logout(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != "" {
		middleware.Refresh(token)
		getuser, _ := c.Get("user") //取值 实现了跨中间件取值
		if getuser != nil {
			user := getuser.(*middleware.UserClaims)
			model.DB().Table("business_account").Where("id", user.ID).Data(map[string]interface{}{"loginstatus": 0}).Update()
		}
	}
	results.Success(c, "退出登录", true, nil)
}

/**
*  6获取验证码
 */
func (api *Index) Get_code(c *gin.Context) {
	mobile := c.DefaultQuery("mobile", "")
	if mobile == "" {
		results.Failed(c, "请填写电话号码", nil)
	} else {
		mobileConfig, _ := model.DB().Table("common_mobile").Where("data_from", "login").First()
		if mobileConfig == nil {
			results.Failed(c, "请联系开发人员，检查配置", nil)
		} else {
			code := gf.GenValidateCode(6)
			err, erro := model.DB().Table("common_verify_code").Data(map[string]interface{}{
				"keyname":    mobile,
				"code":       code,
				"createtime": time.Now().Unix(),
			}).Insert()
			results.Success(c, "获取验证码", err, erro)
		}
	}
}

/**
*7.重置密码
 */
func (api *Index) ResetPassword(c *gin.Context) {
	//获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["code"] == nil || parameter["password"] == nil {
		results.Failed(c, "请提交验证码或密码！", nil)
		return
	}
	password := parameter["password"].(string)
	userdata2, _ := model.DB().Table("business_account").Where("email", parameter["email"]).Fields("id").First()
	if userdata2 == nil {
		results.Failed(c, "邮箱不存在！", nil)
		return
	}
	code, _ := model.DB().Table("common_verify_code").Where("keyname", parameter["email"]).Value("code")
	if code == nil || code != parameter["code"] {
		results.Failed(c, "验证码无效", nil)
		return
	}
	rnd := rand.New(rand.NewSource(6))
	salt := strconv.Itoa(rnd.Int())
	pass := gf.Md5(password + salt)
	res, err := model.DB().Table("business_account").Where("id", userdata2["id"]).Data(map[string]interface{}{
		"salt":     salt,
		"password": pass,
	}).Update()
	if err != nil {
		results.Failed(c, "重置密码失败", err)
	} else {
		results.Success(c, "重置密码成功", res, nil)
	}
}

/**
*  8 获取登录页面信息
 */
func (api *Index) Get_logininfo(c *gin.Context) {
	res2, err := model.DB().Table("common_logininfo").Where("type", "business").OrWhere("type", "common").Fields("title,des,image").Order("weigh asc,id desc").Get()
	if err != nil {
		results.Failed(c, "获取登录页面失败", err)
	} else {
		results.Success(c, "获取登录页面成功！", res2, nil)
	}
}
