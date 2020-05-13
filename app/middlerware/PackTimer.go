package middlerware

import (
	"math/rand"
	"time"

	"github.com/wmyi/gn/gn"
)

type PackTimer struct {
}

func (t *PackTimer) Before(pack gn.IPack) {
	reqId := rand.Intn(1 << 10)
	nowTime := time.Now()
	pack.SetContextValue("reqId", reqId)
	pack.SetContextValue("inTime", nowTime)
	pack.GetLogger().Infof("Before reqId: %d    time %v  ", reqId, nowTime)
}

func (t *PackTimer) After(pack gn.IPack) {
	reqId := pack.GetContextValue("reqId").(int)
	nowTime := pack.GetContextValue("inTime").(time.Time)
	pack.GetLogger().Infof("After reqId: %d   diff time %v  ", reqId, time.Now().Sub(nowTime))
}
