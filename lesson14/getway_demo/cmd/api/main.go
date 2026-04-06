package main

import "getway_demo/internal/common"

func main() {
	common.RunHertz(
		common.EnvOrDefault("SERVICE_NAME", "api"),
		common.EnvOrDefault("INSTANCE_NAME", "api-1"),
		common.EnvOrDefault("PORT", "8080"),
		false,
	)
}
