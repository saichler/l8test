package tests

import (
	"fmt"
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/layer8/go/overlay/health"
	"github.com/saichler/layer8/go/overlay/protocol"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestServicePointValid(t *testing.T) {

	t_resources.CreateTestModelInstance(1)

	for vnetNum := 1; vnetNum <= 3; vnetNum++ {
		for vnicNum := 1; vnicNum <= 4; vnicNum++ {
			nic := topo.VnicByVnetNum(vnetNum, vnicNum)
			hc := health.Health(nic.Resources())
			hp := hc.All()
			if len(hp) != 15 {
				t_resources.Log.Fail(t, "Expected ", nic.Resources().SysConfig().LocalAlias,
					" to have 15 heath points, but it has ", len(hp))
				return
			}
		}
	}

	time.Sleep(time.Second * 5)
	fmt.Println("Messages created before:", protocol.MessagesCreated())
	time.Sleep(time.Second * 5)
	fmt.Println("Messages created after:", protocol.MessagesCreated())
}
