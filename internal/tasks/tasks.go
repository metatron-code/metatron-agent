package tasks

type Task interface {
	Run() ([]byte, error)
}
