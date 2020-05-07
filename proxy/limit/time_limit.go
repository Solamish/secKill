package limit

type Limit struct {
	secLimit TimeLimit
}

var (
	lockChan = make(chan struct{},1)
)

type TimeLimit interface {
	Count(nowTime int64) (curCount int)
	Check(nowTime int64) int
}

type SecLimit struct {
	count   int
	curTime int64
}

// 记录每秒用户访问频率
func (m *SecLimit) Count(nowTime int64) (curCount int) {

	// 超过一秒钟，重新计数
	lockChan <- struct{}{}
	if nowTime-m.curTime > 60 {
		m.curTime = nowTime
		m.count = 1
		curCount = m.count
		<- lockChan
		return
	}

	m.count++
	curCount = m.count
	<- lockChan
	return
}

func (m *SecLimit) Check(nowTime int64) int {
	if nowTime-m.curTime > 60 {
		return 0
	}
	return m.count
}
