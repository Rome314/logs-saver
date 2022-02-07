package eventEntities

type IpInfoManager interface {
	IpInfoProvider
	SetIpInfo(ip string, info *IpInfo) (id int32, err error)
}

type IpInfoProvider interface {
	GetIpInfo(ip string) (info *IpInfo, err error)
}

type BufferRepo interface {
	Store(event *Event) (bufferSize uint64, err error)
	StoreToErrorStorage(events []*Event) (err error)
	PopAll() (events []*Event, err error)
	Status() error
}

type Repository interface {
	Store(event *Event) (err error)
	StoreMany(events ...*Event) (insertedCount int64, err error)
	Status() error
}
