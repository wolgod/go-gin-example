package zgg_service

import (
	"encoding/json"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/shopspring/decimal"
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
	logging.Info("拉取中概股数据")
	beg := time.Now()
	err, zggStockValues := t.GetZggStockValueBySohu()
	if err == nil {
		insertValues := make([]models.ZggStockValue, 0)
		for _, val := range zggStockValues {
			insertValues = append(insertValues, val.(models.ZggStockValue))
		}
		value := models.ZggStockValue{}
		value.InsertStockValue(insertValues)
	} else {
		logging.Errorf("获取数据失败,错误信息:%s", err.Error())
	}
	logging.Infof("time consumed: %f", time.Now().Sub(beg).Seconds())

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

//从搜狐财经获取数据
func (t ZggJob) GetZggStockValueBySohu() (error, []interface{}) {

	bucket := &util.Bucket{}
	beg := time.Now()
	logging.Info("请求地址:", setting.AppSetting.SohuUrl)
	body, err := util.SendHttpRequest(http.MethodGet, setting.AppSetting.SohuUrl, "", nil)
	if err != nil {
		logging.Error("请求接口失败")
	} else {
		m := make(map[string]interface{})
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(body)
		str := string(decodeBytes)
		result := str[11 : len(str)-1]
		err := json.Unmarshal([]byte(result), &m)
		if err != nil {
			logging.Error("接口数据解析失败")
			return err, nil
		}
		chgrateTop := m["chgrate_top"].([]interface{})
		for _, chgrate := range chgrateTop {
			data := chgrate.([]interface{})
			code := data[0].(string)
			name := data[1].(string)
			decimal.DivisionPrecision = 2
			currentPrice := decimal.NewFromFloat(data[3].(float64))
			stockValueUs := decimal.NewFromFloat(data[8].(float64)).Div(decimal.NewFromFloat(100000000))
			zggStockvalue := models.ZggStockValue{
				Name:         name,
				Code:         code,
				StockValueUs: stockValueUs.InexactFloat64(),
				Date:         time.Now(),
				CurrentPrice: currentPrice.InexactFloat64(),
			}
			bucket.Slice = append(bucket.Slice, zggStockvalue)
		}
		logging.Infof("time consumed: %f", time.Now().Sub(beg).Seconds())
	}
	bucket.By = func(a, b interface{}) bool {
		return a.(models.ZggStockValue).StockValueUs > b.(models.ZggStockValue).StockValueUs
	}
	sort.Sort(bucket)
	return nil, bucket.Slice
}
