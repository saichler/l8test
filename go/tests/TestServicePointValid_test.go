package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/saichler/l8bus/go/overlay/health"
	"github.com/saichler/l8bus/go/overlay/protocol"
	"github.com/saichler/l8test/go/infra/t_resources"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestServiceValid(t *testing.T) {

	t_resources.CreateTestModelInstance(1)

	for vnetNum := 1; vnetNum <= 3; vnetNum++ {
		for vnicNum := 1; vnicNum <= 4; vnicNum++ {
			nic := topo.VnicByVnetNum(vnetNum, vnicNum)
			hp := health.All(nic.Resources())
			if len(hp) != 15 {
				t_resources.Log.Fail(t, "Expected ", nic.Resources().SysConfig().LocalAlias,
					" to have 15 heath points, but it has ", len(hp))
				return
			}
		}
	}

	fmt.Println("Before Total - ", protocol.MsgLog.Total())
	time.Sleep(time.Second * 10)
	protocol.MsgLog.Print()
}
