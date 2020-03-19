package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

const adminAuthKey = "admin"
const userAuthKey = "admin"

type User struct {
	Id   int64
	Name string
}

func DummyUserRepository() (User, error) {
	user := User{
		Id:   1,
		Name: "Name",
	}
	return user, nil
}

func DummyMiddleware(authKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if funk.IndexOf(authKeys, adminAuthKey) != -1 {
			user, err := DummyUserRepository()
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Set(gin.AuthUserKey, user)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func GetDummyEndpoint(c *gin.Context) {
	user, _ := c.Get(gin.AuthUserKey)
	resp := map[string]string{"Name": user.(User).Name, "Id": cast.ToString(user.(User).Id)}
	c.JSON(200, resp)
}

func main() {
	api := gin.Default()
	v1 := api.Group("/v1")
	adminRoute := v1.Group("/admin")
	adminRoutenMiddleware := DummyMiddleware(adminAuthKey, userAuthKey)
	adminRoute.Use(adminRoutenMiddleware)
	{
		adminRoute.GET("/", GetDummyEndpoint)
	}
	api.Run(":5000")
}
