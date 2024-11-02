package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-study/webook/internal/service"
	svcmocks "go-study/webook/internal/service/mocks"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestUserHandler_SignUp(t *testing.T) {
//	testCases := []struct {
//		name     string
//		mock     func(ctrl *gomock.Controller) service.UserServiceInterface
//		reqBody  string
//		wantCode int
//		wantBody string
//	}{
//		{
//			name: "注册成功",
//			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
//				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
//				usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
//				return usersvc
//			},
//			reqBody: `{
//				"phone": "15930989104",
//				"confirmPassword": "helloweo1",
//				"password": "helloweo1",
//				"verification":"469897"
//			}`,
//			wantCode: http.StatusOK,
//			wantBody: "success",
//		},
//		{
//			name: "手机号不正确",
//			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
//				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
//				return usersvc
//			},
//			reqBody: `{
//				"phone": "159309891041",
//				"confirmPassword": "helloweo1",
//				"password": "helloweo1",
//				"verification":"469897"
//			}`,
//			wantCode: http.StatusOK,
//			wantBody: "手机号不正确",
//		},
//		{
//			name: "两次输入的密码不相同",
//			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
//				usersvc := svcmocks.NewMockUserServiceInterface(ctrl)
//				return usersvc
//			},
//			reqBody: `{
//				"phone": "15930989104",
//				"confirmPassword": "helloweo2",
//				"password": "helloweo1",
//				"verification":"469897"
//			}`,
//			wantCode: http.StatusOK,
//			wantBody: "两次输入的密码不相同",
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			server := gin.Default()
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			h := NewUserHandler(tc.mock(ctrl))
//			h.RegisterUserRoutes(server)
//			//生成一个请求
//			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
//			require.NoError(t, err)
//			//数据是json
//			req.Header.Set("Content-Type", "application/json")
//			//接受response响应
//			resp := httptest.NewRecorder()
//			//HTTP请求进去gin框架的入口
//			server.ServeHTTP(resp, req)
//			assert.Equal(t, tc.wantCode, resp.Code)
//			assert.Equal(t, tc.wantBody, resp.Body.String())
//		})
//	}
//}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServiceInterface
		wantCode int
		wantData string
		reqBody  string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				svc := svcmocks.NewMockUserServiceInterface(ctrl)
				svc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return svc
			},
			reqBody: `{
				"phone": "15930989104",
				"confirmPassword": "helloweo1",
				"password": "helloweo1",
				"verification":"469897"		
			}`,
			wantCode: http.StatusOK,
			wantData: "success",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := tc.mock(ctrl)
			u := NewUserHandler(svc)
			u.RegisterUserRoutes(server)
			req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantData, resp.Body.String())
		})
	}
}
