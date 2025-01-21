package utils

import "fmt"

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Errorf("panic: %v", err)
			}
		}()
		f()
	}()
}
