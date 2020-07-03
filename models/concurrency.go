package models

import "sync"

type Concurrencyy struct {
	Wg sync.WaitGroup
	Mutexx sync.Mutex
}
