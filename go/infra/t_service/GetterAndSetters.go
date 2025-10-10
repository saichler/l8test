package t_service

import "sync/atomic"

func (this *TestServiceBase) Reset() {
	this.postNumber.Add(this.postNumber.Load() * -1)
	this.putNumber.Add(this.putNumber.Load() * -1)
	this.patchNumber.Add(this.patchNumber.Load() * -1)
	this.getNumber.Add(this.getNumber.Load() * -1)
	this.deleteNumber.Add(this.deleteNumber.Load() * -1)
	this.errorMode = false
}

func (this *TestServiceBase) Name() string {
	return this.name
}

func (this *TestServiceBase) ErrorMode() bool {
	return this.errorMode
}

func (this *TestServiceBase) SetErrorMode(b bool) {
	this.errorMode = b
}

func (this *TestServiceBase) PostN() int {
	return int(this.postNumber.Load())
}

func (this *TestServiceBase) PostNReplica() map[string]atomic.Int32 {
	return this.postReplica
}

func (this *TestServiceBase) PutN() int {
	return int(this.putNumber.Load())
}

func (this *TestServiceBase) PatchN() int {
	return int(this.patchNumber.Load())
}

func (this *TestServiceBase) GetN() int {
	return int(this.getNumber.Load())
}

func (this *TestServiceBase) GetNReplica() map[string]atomic.Int32 {
	return this.getReplica
}

func (this *TestServiceBase) DeleteN() int {
	return int(this.deleteNumber.Load())
}

func (this *TestServiceBase) FailedN() int {
	return int(this.failedNumber.Load())
}
