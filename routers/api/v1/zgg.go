package v1

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/service/zgg_service"
	_ "github.com/EDDYCJY/go-gin-example/service/zgg_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Get multiple articles
// @Produce  json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [get]
func GetZggStockValue(c *gin.Context) {
	appG := app.Gin{C: c}
	err, values := models.GetZggStockValue()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ZGG_DATA_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = values
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func GetZggStockValueCurrent(c *gin.Context) {
	appG := app.Gin{C: c}
	zggJob := &zgg_service.ZggJob{}
	err, value := zggJob.GetZggStockValueBySohu()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ZGG_DATA_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = value
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
