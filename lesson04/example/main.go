package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

var sentence string

//	func main() {
//		fmt.Print("请输入一句话：") /*1234567890
//		  123456789*/
//		_, err := fmt.Scanln(&sentence)
//		if err != nil {
//			fmt.Print(err)
//		}
//		fmt.Println(sentence)
//		fmt.Scanln(&sentence)
//	}

//func main() { //go run main.go 乌鲁鲁
//	var args []string = os.Args
//	if len(args) < 2 {
//		fmt.Println("你没有提供名字哦")
//		return
//	}
//	name := args[1] // 获取第一个命令行参数
//
//	fmt.Printf("你好, %s!\n", name)
//	if name == "乌鲁鲁" {
//		fmt.Println("堵桥来！")
//	}
//}

func main() {
	// 1. 创建节点（参数为机器ID，范围0-1023）
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	// 2. 生成ID
	id := node.Generate()

	// 3. 多种格式输出
	fmt.Printf("Int64  ID: %d\n", id)          // 整数形式
	fmt.Printf("String ID: %s\n", id)          // 字符串形式
	fmt.Printf("Base2  ID: %s\n", id.Base2())  // 二进制形式
	fmt.Printf("Base64 ID: %s\n", id.Base64()) // Base64编码

	// 4. 解析ID信息
	fmt.Printf("时间戳: %d\n", id.Time())  // 毫秒时间戳
	fmt.Printf("节点ID: %d\n", id.Node()) // 机器ID
	fmt.Printf("序列号: %d\n", id.Step())  // 序列号
}
