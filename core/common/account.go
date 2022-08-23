package common

type AccountData struct {
	Balance     int64
	Nonce       int64
	CodeHash    Hash
	StorageHash Hash
}
