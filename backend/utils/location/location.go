package location

import (
	"CengkeHelper/logger"
	"CengkeHelper/setup"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func IpToLocation(ip string) string {

	type webIpLocationJson struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Ip       string `json:"ip"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			AdInfo struct {
				Nation   string `json:"nation"`
				Province string `json:"province"`
				City     string `json:"city"`
				District string `json:"district"`
				Adcode   int    `json:"adcode"`
			} `json:"ad_info"`
		} `json:"result"`
	}

	client := resty.New()

	// 计算签名
	str := fmt.Sprintf("/ws/location/v1/ip?ip=%v&key=%v%v",
		ip, setup.Config.LocationKey, setup.Config.SecretKey)
	hash := md5.Sum([]byte(str))
	md5Str := hex.EncodeToString(hash[:])

	// 请求
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"key": setup.Config.LocationKey,
			"ip":  ip,
			"sig": md5Str,
		}).
		SetResult(&webIpLocationJson{}).
		Get("https://apis.map.qq.com/ws/location/v1/ip")

	if err != nil {
		// 请求太频繁
		logger.Warning(err)
		return err.Error()
	}
	resJson := resp.Result().(*webIpLocationJson)

	if resJson.Status != 0 {
		logger.Debug(resp)
	}

	resStr := resJson.Result.AdInfo.Nation + "-" + resJson.Result.AdInfo.Province

	if resJson.Result.AdInfo.City != "" {
		resStr += "-" + resJson.Result.AdInfo.City
	}

	if resJson.Result.AdInfo.District != "" {
		resStr += "-" + resJson.Result.AdInfo.District
	}

	return resStr

}
