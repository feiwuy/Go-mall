package main

import (
	"encoding/json"
	"fmt"
	"gomall/config"
	"gomall/routers"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Username string	`json:"username"`
	Password string	`json:"password"`
}

func TestMain(m *testing.M) {
	// 抑制调试输出
	gin.SetMode(gin.TestMode)
	config.InitConfig()
	config.AppConfig.DataBase.Name = "go_mall_test_db" // 切测试库
	config.InitDB()
	os.Exit(m.Run())
}


// RegisterTest 注册模块
func RegisterTest(t *testing.T, username, password string, code int) {
	router := routers.SetUpRouter()
	body := fmt.Sprintf(`{"username":%q, "password":%q}`, username, password)

	// 注册逻辑
	req, _ := http.NewRequest("POST", "/api/users", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code, "body=%s", w.Body.String())
	var resp map[string]string
	switch code {
		case http.StatusCreated:
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "注册成功!", resp["msg"])
		case http.StatusConflict:
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "该用户名已被占用", resp["error"])
		case http.StatusBadRequest:
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			if username == "" {
				assert.Equal(t, "用户名不能为空", resp["error"])
			} else if len(password) < 6 {
				assert.Equal(t, "密码至少6位", resp["error"])
			} else {
				assert.Equal(t, "非法请求", resp["error"])
			}
	}
}

// LoginTest 登录模块
func LoginTest(t *testing.T, username, password string, code int) string {
	router := routers.SetUpRouter()

	// %q自带双引号
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	req, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)
	var resp map[string]string
	// 不同状态码有不同的测试流程
	switch code {
		case http.StatusOK:
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			if resp["token"] == "" {
				t.Fatal("登录成功但未返回token")
			} else {
				return resp["token"]
			}
		case http.StatusUnauthorized:
			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "用户名或密码错误", resp["error"])
		case http.StatusBadRequest:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "非法请求", resp["error"])
		case http.StatusInternalServerError:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "服务器错误,请稍后再试", resp["error"])
	}
	return ""
}

// 用户测试
// 包括: 注册成功, 注册错误, 登录错误, 成功登录
func TestUserOption(t *testing.T) {
	// 存放本轮测试使用的信息
	var user = User{
		Username: fmt.Sprintf("test_%d", time.Now().UnixNano()),
		Password: "123456",
	}

	// 注册测试
	RegisterTest(t, "", user.Password, http.StatusBadRequest)
	RegisterTest(t, user.Username, "1234", http.StatusBadRequest)
	// 成功注册
	RegisterTest(t, user.Username, user.Password, http.StatusCreated)
	// 用户名被占用
	RegisterTest(t, user.Username, user.Password, http.StatusConflict)

	// 登录测试
	LoginTest(t, user.Username+"test", user.Password, http.StatusUnauthorized)
	LoginTest(t, user.Username, user.Password+"test", http.StatusUnauthorized)
	// 成功登录
	LoginTest(t, user.Username, user.Password, http.StatusOK)
}

func ProfileGetTest(t *testing.T, username, password string, code int, token string) {
	router := routers.SetUpRouter()

	body := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	req, _ := http.NewRequest("GET", "/api/users", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)

	switch code {
		case http.StatusOK:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "查询成功", resp["msg"])
		case http.StatusUnauthorized:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			if token == "" {
				assert.Equal(t, "Missing Authorization Header", resp["error"])
			} else {
				assert.Equal(t, "登录失效,请重新登录", resp["error"])
			}
	}
}

// ProfileUpdateTest 修改模块
func ProfileUpdateTest(t *testing.T, username, password string, code int, token string) {
	router := routers.SetUpRouter()

	body := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	req, _ := http.NewRequest("PATCH", "/api/users", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)

	switch code {
		case http.StatusOK:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "修改成功", resp["msg"])
		case http.StatusConflict:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "该用户名已被占用", resp["error"])
		case http.StatusBadRequest:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			if username == "" {
				assert.Equal(t, "用户名不能为空", resp["error"])
			} else if len(password) < 6 {
				assert.Equal(t, "密码至少6位", resp["error"])
			} else {
				assert.Equal(t, "非法请求", resp["error"])
			}
		case http.StatusInternalServerError:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "服务器错误,请稍后再试", resp["error"])
	}
}

// TestProfile
// 包括: 查看资料, 未登录访问, 修改用户名和密码, 旧密码失效, 新密码成功
func TestProfile(t *testing.T) {
	// 存放本轮测试使用的信息
	var user = User{
		Username: fmt.Sprintf("test_%d", time.Now().UnixNano()),
		Password: "123456",
	}
	var other = User {
		Username: fmt.Sprintf("test_%d", time.Now().UnixNano()),
		Password: "123456",
	}

	// user注册
	RegisterTest(t, user.Username, user.Password, http.StatusCreated)
	token := LoginTest(t, user.Username, user.Password, http.StatusOK)

	// 查看资料
	ProfileGetTest(t, user.Username, user.Username, http.StatusUnauthorized, "")
	// 过期token测试
	tokenTest := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "whatever",
		"userID":   uint(1),
		"exp":      time.Now().Add(-time.Hour).Unix(),
	})
	signed, err := tokenTest.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}
	tokenNew := "Bearer " + signed
	ProfileGetTest(t, user.Username, user.Username, http.StatusUnauthorized, tokenNew)
	ProfileGetTest(t, user.Username, user.Username, http.StatusOK, token)

	// 新用户名和密码
	usernameNew := fmt.Sprintf(`{"username":"test_%d"}`, time.Now().UnixNano())
	passwordNew := "test_new"

	// other注册
	RegisterTest(t, other.Username, other.Password, http.StatusCreated)

	// 更新账号密码
	ProfileUpdateTest(t, "", passwordNew, http.StatusBadRequest, token)
	ProfileUpdateTest(t, usernameNew, "1234", http.StatusBadRequest, token)
	ProfileUpdateTest(t, "", "1234", http.StatusBadRequest, token)
	ProfileUpdateTest(t, other.Username, user.Password, http.StatusConflict, token)
	ProfileUpdateTest(t, usernameNew, passwordNew, http.StatusOK, token)

	// 旧密码失效
	LoginTest(t, user.Username, user.Password, http.StatusUnauthorized)
	// 新密码成功
	LoginTest(t, usernameNew, passwordNew, http.StatusOK)
}

type Todo struct {
	Title 		string
	Content 	string
	Completed	bool
}

type TodoCreateResp struct {
	Msg string	`json:"msg"`
	ID	uint	`json:"id"`
}

func TodoCreateTest(t *testing.T, token string, code int, todo Todo) string {
	router := routers.SetUpRouter()

	body := fmt.Sprintf(`{"title":%q,"content":%q,"completed":%t}`, todo.Title, todo.Content, todo.Completed)
	req, _ := http.NewRequest("POST", "/api/todos", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)

	switch code {
		case http.StatusCreated:
			var resp TodoCreateResp
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "创建成功", resp.Msg)
			return fmt.Sprintf("%d", resp.ID)
	case http.StatusBadRequest:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "非法请求", resp["error"])
	case http.StatusUnauthorized:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "登录失效,请重新登录", resp["error"])
	case http.StatusInternalServerError:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "服务器错误,请稍后再试", resp["error"])
	}
	return ""
}

func TodoDeleteTest(t *testing.T, token, id string, code int) {
	router := routers.SetUpRouter()

	url := fmt.Sprintf("/api/todos/%s", id)
	req, _ := http.NewRequest("DELETE", url, strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)

	switch code {
		case http.StatusOK:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "删除成功", resp["msg"])
		case http.StatusNotFound:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "不存在或无权访问", resp["error"])
		case http.StatusInternalServerError:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "服务器错误,请稍后再试", resp["error"])
	}
}

func TodoUpdateTest(t *testing.T, token, id string, code int, todo Todo) {
	router := routers.SetUpRouter()

	url := fmt.Sprintf("/api/todos/%s", id)
	body := fmt.Sprintf(`{"title":%q,"content":%q,"completed":%t}`, todo.Title, todo.Content, todo.Completed)
	req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)

	switch code {
		case http.StatusOK:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "更新成功", resp["msg"])
		case http.StatusNotFound:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "不存在或无权访问", resp["error"])
		case http.StatusBadRequest:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "非法请求", resp["error"])
		case http.StatusInternalServerError:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "服务器错误,请稍后再试", resp["error"])
	}
}

func TodoReadTest(t *testing.T, id, token string, code int) {
	router := routers.SetUpRouter()

	url := fmt.Sprintf("/api/todos/%s", id)
	req, _ := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)

	switch code {
		case http.StatusCreated:
			var resp map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "获取成功", resp["msg"])
		case http.StatusNotFound:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "不存在或无权访问", resp["error"])
		case http.StatusInternalServerError:
			var resp map[string]string
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, "服务器错误,请稍后再试", resp["error"])
	}
}

// TestTodo
// 包括: 
func TestTodo(t *testing.T) {
	// 存放本轮测试使用的信息
	var user = User{
		Username: fmt.Sprintf("test_%d", time.Now().UnixNano()),
		Password: "123456",
	}
	RegisterTest(t, user.Username, user.Password, http.StatusCreated)
	token := LoginTest(t, user.Username, user.Password, http.StatusOK)

	// 创建初始化todo
	var todo = Todo {
		Title: "test",
		Content: "test_content",
		Completed: false,
	}

	// 创建todo
	TodoCreateTest(t, "1", http.StatusUnauthorized, todo)
	id := TodoCreateTest(t, token, http.StatusCreated, todo)
	assert.NotEqual(t, "", id)

	// 读取todo
	TodoReadTest(t, id, token, http.StatusOK)

	// 更新后的todo
	var todoNew = Todo {
		Title: "test2",
		Content: "test_content2",
		Completed: true,
	}
	// 更新todo
	TodoUpdateTest(t, token, id+"1", http.StatusNotFound, todoNew)
	TodoUpdateTest(t, token+"1", id, http.StatusUnauthorized, todoNew)
	TodoUpdateTest(t, token, id, http.StatusOK, todoNew)

	// 删除todo
	TodoDeleteTest(t, token, id+"1", http.StatusNotFound)
	TodoDeleteTest(t, token+"1", id, http.StatusUnauthorized)
	TodoDeleteTest(t, token, id, http.StatusOK)
}
