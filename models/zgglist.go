package models

import (
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
)

type ZggList struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func GetCodeList() (error, []ZggList) {
	var zggList []ZggList

	if err := db.Table("zgg_list").Where("status=0").Find(&zggList).Error; err != nil {
		logging.Errorf("查询中概股列表失败,错误信息:%s", err)
		return err, nil
	}
	return nil, zggList

}
