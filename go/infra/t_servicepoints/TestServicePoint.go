package t_servicepoints

import (
	"errors"
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/testtypes"
	"github.com/saichler/types/go/types"
	"sync/atomic"
)

type TestServicePointHandler struct {
	name             string
	postNumber       atomic.Int32
	putNumber        atomic.Int32
	patchNumber      atomic.Int32
	deleteNumber     atomic.Int32
	getNumber        atomic.Int32
	failedNumber     atomic.Int32
	tr               bool
	errorMode        bool
	replicationCount int
	replicationScore int
}

const (
	ServiceName = "Tests"
)

func NewTestServicePointHandler(name string) *TestServicePointHandler {
	tsp := &TestServicePointHandler{}
	tsp.name = name
	return tsp
}

func (this *TestServicePointHandler) Post(pb common.IMObjects, resourcs common.IResources) common.IMObjects {
	Log.Debug("Post -", this.name, "- Test callback")
	this.postNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Post - TestServicePointHandler Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointHandler) Put(pb common.IMObjects, resourcs common.IResources) common.IMObjects {
	Log.Debug("Put -", this.name, "- Test callback")
	this.putNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Put - TestServicePointHandler Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointHandler) Patch(pb common.IMObjects, resourcs common.IResources) common.IMObjects {
	Log.Debug("Patch -", this.name, "- Test callback")
	this.patchNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Patch - TestServicePointHandler Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointHandler) Delete(pb common.IMObjects, resourcs common.IResources) common.IMObjects {
	Log.Debug("Delete -", this.name, "- Test callback")
	this.deleteNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Delete - TestServicePointHandler Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointHandler) GetCopy(pb common.IMObjects, resourcs common.IResources) common.IMObjects {
	Log.Debug("Get -", this.name, "- Test callback")
	this.getNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("GetCopy - TestServicePointHandler Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointHandler) Get(pb common.IMObjects, resourcs common.IResources) common.IMObjects {
	Log.Debug("Get -", this.name, "- Test callback")
	this.getNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Get - TestServicePointHandler Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointHandler) Failed(pb common.IMObjects, resourcs common.IResources, info *types.Message) common.IMObjects {
	dest := "n/a"
	msg := "n/a"
	if info != nil {
		dest = info.Source
		msg = info.FailMsg
	}
	Log.Debug("Failed -", this.name, " to ", dest, "- Test callback")
	Log.Debug("Failed Reason is ", msg)
	this.failedNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Failed - TestServicePointHandler Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointHandler) EndPoint() string {
	return "/Tests"
}
func (this *TestServicePointHandler) ServiceName() string {
	return ServiceName
}
func (this *TestServicePointHandler) ServiceModel() common.IMObjects {
	return New(nil, &testtypes.TestProto{})
}
func (this *TestServicePointHandler) Transactional() bool {
	return this.tr
}
func (this *TestServicePointHandler) ReplicationCount() int {
	return this.replicationCount
}
func (this *TestServicePointHandler) ReplicationScore() int {
	return this.replicationScore
}
