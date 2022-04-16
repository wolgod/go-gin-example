package zgg_service

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	routineCountTotal = 10 //限制线程数
	URL               = "http://qt.gtimg.cn/"
)

type ZggJob struct {
}

func (t ZggJob) Run() {

	zggStockValues := t.GetZggStockValue()
	insertValues := make([]models.ZggStockValue, 0)
	for _, val := range zggStockValues {
		insertValues = append(insertValues, val.(models.ZggStockValue))
	}
	value := models.ZggStockValue{}
	value.InsertStockValue(insertValues)
}

func (t ZggJob) GetZggStockValue() []interface{} {

	result := &util.Bucket{}
	err, list := models.GetCodeList()
	if err != nil {
		logging.Error("执行失败....")
	}
	if len(list) > 0 {
		g := util.NewG(routineCountTotal)
		wg := &sync.WaitGroup{}
		beg := time.Now()
		for _, value := range list {
			wg.Add(1)
			g.Run(func(args interface{}) {
				defer wg.Done()
				val := args.(models.ZggList)
				reqParam := fmt.Sprintf("q=s_us%s", val.Code)
				logging.Info("请求地址:", URL+"?"+reqParam)
				body, err := util.SendHttpRequest(http.MethodGet, URL+"?"+reqParam, "text/html; charset=GBK", nil)
				if err != nil {
					logging.Error("请求接口失败")
				} else {
					//v_s_usLI="200~理想汽车~LI.OQ~25.71~-0.32~-1.23~4646805~120581802~265.54772~";
					var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(body)
					str := string(decodeBytes)
					split := strings.Split(str, "=")
					//res:=split[1][1:len(split[1])-1]
					//获取市值
					zgg := strings.Split(split[1], "~")
					currentPrice, _ := strconv.ParseFloat(zgg[3], 64)
					stockValue, _ := strconv.ParseFloat(zgg[8], 64)
					zggStockvalue := models.ZggStockValue{
						Name:         val.Name,
						Code:         val.Code,
						StockValueUs: stockValue,
						Date:         time.Now(),
						CurrentPrice: currentPrice,
					}
					result.Slice = append(result.Slice, zggStockvalue)
				}
			}, value)
		}
		wg.Wait()
		logging.Infof("time consumed: %f", time.Now().Sub(beg).Seconds())
	}
	result.By = func(a, b interface{}) bool {
		return a.(models.ZggStockValue).StockValueUs > b.(models.ZggStockValue).StockValueUs
	}
	sort.Sort(result)
	return result.Slice
}
