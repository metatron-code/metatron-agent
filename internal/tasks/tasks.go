package tasks

import "time"

type Task interface {
	Run(timeout time.Duration) ([]byte, error)
}
