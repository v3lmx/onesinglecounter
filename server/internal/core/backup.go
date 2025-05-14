package core

type Backup interface {
	Backup(current uint64, best Best) error
	Recover() (uint64, Best, error)
}
