package main

import (
	"github.com/droptheplot/rand_chat/env"
)

func main() {
	db := env.Init()

	env.Migrate(db)
}
