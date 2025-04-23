package t_servicepoints

func (this *TestServicePointBase) Reset() {
	this.postNumber.Add(this.postNumber.Load() * -1)
	this.putNumber.Add(this.putNumber.Load() * -1)
	this.patchNumber.Add(this.patchNumber.Load() * -1)
	this.getNumber.Add(this.getNumber.Load() * -1)
	this.deleteNumber.Add(this.deleteNumber.Load() * -1)
	this.errorMode = false
}

func (this *TestServicePointBase) Name() string {
	return this.name
}

func (this *TestServicePointBase) ErrorMode() bool {
	return this.errorMode
}

func (this *TestServicePointBase) SetErrorMode(b bool) {
	this.errorMode = b
}

func (this *TestServicePointBase) PostN() int {
	return int(this.postNumber.Load())
}

func (this *TestServicePointBase) PutN() int {
	return int(this.putNumber.Load())
}

func (this *TestServicePointBase) PatchN() int {
	return int(this.patchNumber.Load())
}

func (this *TestServicePointBase) GetN() int {
	return int(this.getNumber.Load())
}

func (this *TestServicePointBase) DeleteN() int {
	return int(this.deleteNumber.Load())
}

func (this *TestServicePointBase) FailedN() int {
	return int(this.failedNumber.Load())
}
