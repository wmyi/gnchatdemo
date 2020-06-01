package middlerware

import (
	"math/rand"
	"time"

	logger "github.com/wmyi/gn/glog"
	"github.com/wmyi/gn/gn"
)

type PackTimer struct {
}

func (t *PackTimer) Before(pack gn.IPack) {
	reqId := rand.Intn(1 << 10)
	nowTime := time.Now()
	pack.SetContextValue("reqId", reqId)
	pack.SetContextValue("inTime", nowTime)
	logger.Infof("Before reqId: %d    time %v  ", reqId, nowTime)
}

func (t *PackTimer) After(pack gn.IPack) {
	reqId := pack.GetContextValue("reqId").(int)
	nowTime := pack.GetContextValue("inTime").(time.Time)
	logger.Infof("After reqId: %d   diff time %v  ", reqId, time.Now().Sub(nowTime))
}
