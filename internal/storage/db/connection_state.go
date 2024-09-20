package db

const (
	NotConnected ConnectionState = iota
	InvalidConnectionString
	RefusedConnection
	Connected
)

type ConnectionState int
