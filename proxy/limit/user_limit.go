package limit

import (
	"math/rand"
	"time"
)

// 用户抽取一个随机数字
// 如果抽到的数字大于幸运数字
// 用户的请求将被拦截
func Probability(luckyNum float64) bool{
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())
	randNumber := rand.Float64()
	if randNumber > luckyNum {
		return false
	}
	return true
}

 

