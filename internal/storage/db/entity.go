package db

type Entity interface {
	Deconstruct() []interface{}
}
