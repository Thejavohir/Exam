package repo

type NewRepo interface {
	SetWithTTL(key, value string, seconds int64) error
	Get(key string) (interface{}, error)
	Exists(key string) (interface{}, error)
}
