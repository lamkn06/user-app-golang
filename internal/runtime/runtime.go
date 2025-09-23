package runtime

import (
	"fmt"
	"os"
	"reflect"

	"github.com/caarlos0/env"
)

func LoadConfigs[C any](configs []C) {
	for _, config := range configs {
		LoadConfig(config)
	}
}

func LoadConfig[C any](c C) {
	FailOnError(
		env.Parse(c),
		fmt.Sprintf("could not load %s from env", reflect.TypeOf(c).String()),
	)
}

func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("FATAL: %s (err=%s)\n", err, msg)
		os.Exit(1)
	}
}
