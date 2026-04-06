package main

import "getway_demo/internal/common"

func main() {
	common.RunHertz(
		common.EnvOrDefault("SERVICE_NAME", "admin"),
		common.EnvOrDefault("INSTANCE_NAME", "admin-1"),
		common.EnvOrDefault("PORT", "8080"),
		false,
	)
}
