package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	_ "github.com/shjtmy/fpgasim/docs"

	"github.com/shjtmy/fpgasim/internal/repositry"
	"github.com/shjtmy/fpgasim/pkg/model"
	"github.com/shjtmy/fpgasim/pkg/dao"
)

// @title FPGA APIドキュメントのタイトル
// @version バージョン(1.0)
// @description 仕様書に関する内容説明

// @contact.name APIサポーター
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:33333
// @BasePath /
func main() {
	r := gin.New()
	url := ginSwagger.URL("http://localhost:33333/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/register",getRegister)
	r.GET("/register/:id",getRegisterID)
	r.PUT("/register", updateRegister)
	r.Run(":33333")
}
// @description テスト用APIの詳細
// @version 1.0
// @accept application/x-json-stream
// @param none query string false "必須ではありません。"
// @Success 200 {object} model.Property 
// @router /getRegister/id [get]
func getRegister(c *gin.Context){
	properties := repositry.Get()
	c.JSON(http.StatusOK, properties)
}

func getRegisterID(c *gin.Context){
	idstr := c.Param("id")
	property := repositry.GetByID(idstr)
	c.JSON(http.StatusOK, properties)
}

// @description テスト用APIの詳細
// @version 1.0
// @accept application/x-json-stream
// @param none query string false "必須ではありません。"
// @Success 200 {object} string 
// @failuer 403 {object} string  
// @failuer 500 {object} string  
// @router /updateRegister/ [put]
func updateRegister(c *gin.Context){
	err := repositry.Update(property)
	c.JSON(http.StatusOK, properties)
}
