package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 简单初始化日志器
	// 使用zap的默认配置，控制台输出，人类可读格式
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // 确保所有日志都被写入
	//logger.Sync() 方法的作用是强制将所有缓冲中的日志消息同步写入到输出目标，确保没有任何日志消息丢失。它通常在程序的退出阶段被调用，以确保所有日志都被正确记录。

	fmt.Println("=== Zap日志库示例 ===")

	// 1. 基本日志记录
	fmt.Println("\n1. 基本日志记录:")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")

	// 2. 结构化日志记录
	fmt.Println("\n2. 结构化日志记录:")
	// 使用字段创建结构化日志，便于后续分析和查询
	logger.Info("用户登录事件",
		zap.String("用户ID", "123456"),
		zap.String("用户名", "张三"),
		zap.String("IP地址", "192.168.1.100"),
		zap.Bool("登录成功", true),
		zap.Time("登录时间", time.Now()),
	)

	// 3. 不同级别的日志
	fmt.Println("\n3. 不同级别的日志:")
	// Debug级别 - 仅在开发模式下显示
	logger.Debug("这是一条调试日志，用于开发调试")
	// Info级别 - 一般信息
	logger.Info("这是一条信息日志，记录正常操作")
	// Warn级别 - 需要注意的情况
	logger.Warn("这是一条警告日志，提示潜在问题")
	// Error级别 - 错误情况，但程序可以继续运行
	logger.Error("这是一条错误日志，记录错误情况")

	// 4. 带上下文的日志
	fmt.Println("\n4. 带上下文的日志:")
	// 创建一个带有固定上下文的日志器
	requestLogger := logger.With(
		zap.String("请求ID", "req-20240101-001"),
		zap.String("API路径", "/api/users"),
	)

	// 使用带上下文的日志器
	requestLogger.Info("开始处理请求")
	// 模拟处理时间
	time.Sleep(50 * time.Millisecond)
	requestLogger.Info("请求处理完成",
		zap.Int("状态码", 200),
		zap.Duration("处理时间", 50*time.Millisecond),
	)

	// 5. 错误处理日志
	fmt.Println("\n5. 错误处理日志:")
	// 模拟一个错误
	err := fmt.Errorf("数据库连接失败: %w", fmt.Errorf("网络超时"))
	// 记录错误，包含错误详情
	logger.Error("操作失败",
		zap.Error(err), // 自动记录错误信息
		zap.String("操作类型", "数据库查询"),
		zap.String("查询语句", "SELECT * FROM users"),
	)

	fmt.Println("\n=== 示例结束 ===")
}
