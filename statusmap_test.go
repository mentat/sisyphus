package main

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {

	sm := NewThreadSafeMap()

	sem := make(chan bool, 100)

	for i := int(0); i < cap(sem); i++ {
		sem <- true

		go func(index int) {
			defer func() { <-sem }()
			d := Status{
				Running: true,
			}
			sm.Set(fmt.Sprintf("key%d", index), d)

			for j := int(0); j < cap(sem); j++ {
				data, ok := sm.Get(fmt.Sprintf("key%d", index))
				if !ok || !data.Running {
					panic("Data bad.")
				}
			}
		}(i)

	}

	for i := int(0); i < cap(sem); i++ {
		sem <- true
	}

}
