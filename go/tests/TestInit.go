package tests

import (
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8test/go/infra/t_topology"
	"github.com/saichler/shared/go/share/logger"
	. "github.com/saichler/types/go/common"
)

var topo *TestTopology
var FLog = logger.NewLoggerDirectImpl(logger.NewFileLogMethod("test.log"))

func init() {
	Log.SetLogLevel(Trace_Level)
}

func setup() {
	setupTopology()
}

func tear() {
	shutdownTopology()
}

func reset(name string) {
	Log.Info("*** ", name, " end ***")
	topo.ResetHandlers()
}

func setupTopology() {
	topo = NewTestTopology(4, []int{20000, 30000, 40000}, Trace_Level)
}

func shutdownTopology() {
	topo.Shutdown()
}
