package middleware

import (
	"github.com/afex/hystrix-go/hystrix"
)

func NewServiceWrapper(name string) {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  10,
		RequestVolumeThreshold: 10,
		SleepWindow:            2000,
		ErrorPercentThreshold:  50,
	})
}
