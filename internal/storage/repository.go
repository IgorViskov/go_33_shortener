package storage

type Repository[tid comparable, tentity any] interface {
	Get(id tid) (*tentity, error)
	Insert(entity *tentity) (*tentity, error)
	Update(entity *tentity) (*tentity, error)
	Delete(id tid) error
	Find(search string) (*tentity, error)
}
