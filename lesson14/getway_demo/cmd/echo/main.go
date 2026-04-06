package main

import "getway_demo/internal/common"

func main() {
	common.RunHertz(
		common.EnvOrDefault("SERVICE_NAME", "echo"),
		common.EnvOrDefault("INSTANCE_NAME", "echo-1"),
		common.EnvOrDefault("PORT", "8080"),
		true,
	)
}
