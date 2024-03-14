package utils

import (
	"fmt"
	"math/rand"
	models "server/models/status"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDefaultString(c *gin.Context, key string, r string) string {
	k := c.GetString(key)
	if len(k) == 0 {
		return r
	}
	return k
}

func GetUtcNowRFC3339() string {
	t := time.Now().UTC().Format(time.RFC3339)
	return t
}

// 比较目标时间和当前时间的差值，返回秒数
func CompareNowTimeWithSeconds(utc_time_str string) (d float64, err error) {
	targetTime, err := time.Parse(time.RFC3339, utc_time_str)
	if err != nil {
		return
	}
	n := time.Now().UTC()
	d = n.Sub(targetTime).Seconds()
	return
}

func FormatUrls(urls []models.NewUrl) []string {
	n := make([]string, 0)
	for i := 0; i < len(urls); i++ {
		n = append(n, fmt.Sprintf(`[%s]%s`, urls[i].UrlKey, urls[i].Url))
	}
	return n
}

func UidsToEdges(uids []string) []map[string]string {
	n := make([]map[string]string, 0)
	for i := 0; i < len(uids); i++ {
		n = append(n, map[string]string{"uid": uids[i]})
	}
	return n
}

func UidsMap(uids []string) []map[string]string {
	n := make([]map[string]string, 0)
	for i := 0; i < len(uids); i++ {
		n = append(n, map[string]string{"uid": uids[i]})
	}
	return n
}

// 生成n位验证码
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
