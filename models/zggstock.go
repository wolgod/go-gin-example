package models

import (
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"time"
)

type ZggStockValue struct {
	Name         string
	Code         string
	Date         time.Time
	StockValueUs float64
	CurrentPrice float64
}

func (t *ZggStockValue) InsertStockValue(value []ZggStockValue) {

	err := db.Table("zgg_stock_value").Create(value).Error
	if err != nil {
		logging.Errorf("insert fail,err:%s", err.Error())
	}
}

func GetZggStockValue() (error, []ZggStockValue) {
	var zggList []ZggStockValue

	date := time.Now().Format("2006-01-02")
	if err := db.Table("zgg_stock_value").Where("date=?", date).Order("stock_value_us desc").Find(&zggList).Error; err != nil {
		logging.Errorf("查询中概股列表市值列表失败,错误信息:%s", err)
		return err, nil
	}
	return nil, zggList

}
