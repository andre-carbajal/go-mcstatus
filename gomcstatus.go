package gomcstat

type Server interface {
	Status() (StatusResponse, error)
	Ping() (int64, error)
}

type StatusResponse interface {
	GetLatency() int64
}
