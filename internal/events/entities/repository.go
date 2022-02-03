package eventEntities

type Repository interface {
	StoreMany(events ...*Event) (inserted int64, err error)
}
