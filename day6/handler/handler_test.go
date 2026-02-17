package handler// package handler: 包名，用于组织代码，类似于Python中的模块

import (
	"net/http"
	"net/http/httptest"//httptest: 测试HTTP服务器
	"testing"//testing: 测试框架
)

// TestGetUser 表驱动测试：多组「路径 + 期望状态码」，循环调用 GetUser 检查。
//
// 和 Python 的对应关系（方便理解）：
//   - tests := []struct{...}{...}  相当于  tests = [{"path": "/api/users/1", "want_status": 200}, ...]
//   - for _, tt := range tests     相当于  for tt in tests:  （_ 是下标，这里不用）
//   - httptest.NewRequest(...)     造一个「假的」HTTP 请求，不经过网络
//   - httptest.NewRecorder()       造一个「录响应」的 Writer，用来接 Handler 写出的状态码和 Body
//   - GetUser(rec, req)            直接调 Handler，相当于  response = get_user(request)；结果写在 rec 里
//   - rec.Code                     相当于  response.status_code
//   - t.Errorf(...)                断言失败时打印信息并标记测试失败，相当于 assert 或 raise
func TestGetUser(t *testing.T) {
	tests := []struct {//[]: 数组,多组输入
		path       string // 请求路径
		wantStatus int    // 期望的 HTTP 状态码，200表示成功，400表示请求错误
	}{// 上面是生命一个结构体，下面是初始化一个结构体，包含多组path和wantStatus
		{"/api/users/1", http.StatusOK},//请求路径为/api/users/1，期望状态码为200
		{"/api/users/42", http.StatusOK},//请求路径为/api/users/42，期望状态码为200
		{"/api/users/0", http.StatusBadRequest},//请求路径为/api/users/0，期望状态码为400
		{"/api/users/abc", http.StatusBadRequest},//请求路径为/api/users/abc，期望状态码为400
	}
	for _, tt := range tests {//for _, tt := range tests: 循环遍历tests中的每个元素
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)//NewRequest: 创建一个HTTP请求, nil表示请求体为空	
		rec := httptest.NewRecorder()//NewRecorder: 创建一个HTTP响应记录器
		GetUser(rec, req) // 本文件也是 package handler，和 handler.go 同包，所以直接调 GetUser，无需 import
		if rec.Code != tt.wantStatus {//rec.Code: 响应记录器的HTTP状态码, tt.wantStatus: 期望的HTTP状态码
			t.Errorf("GetUser(%q) status = %d, want %d", tt.path, rec.Code, tt.wantStatus)//Errorf: 断言失败时打印信息并标记测试失败	
		}
	}//for循环结束
}//TestGetUser函数结束
