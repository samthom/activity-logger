package main

import (
	_ "github.com/lib/pq"
	"github.com/samthom/activity-logger/handlers"
)

func main() {
	handlers.HandleRequests()
}
