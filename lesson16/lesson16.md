# AI Agent 开发入门

---

## 一、核心概念

### 什么是 LLM？

大型语言模型（Large Language Model）。本质上是一个函数：

```
输入：[]{role, content} 消息列表  
输出：一条 assistant 消息
```

不要想太复杂，就是一个很强的"文字接龙"函数(生成不可控/生成过程去看隔壁py的课)

### 什么是 AI Agent？

**Agent 约等于 LLM + 循环 + 工具**

> 不是所有“调用了 Tool 的 LLM”都算完整 Agent
> 是否有“根据结果继续决策的循环”，是一个关键区别。

- 普通 LLM：你问一句，它答一句，结束
- Agent：它可以自己决定"我需要先查一下数据"，调工具，看结果，再继续思考，直到给出最终答案

类比：LLM 是大脑，Tool 是手，循环是思考过程。

### 三层递进

```
第一层：LLM 会说话
  chatModel.Generate(messages) → 回一条消息

第二层：LLM + Tool，有手了
LLM 决定是否调用工具 → 工具执行 → 结果返回给 LLM → 产出回答

第三层：Agent Loop，循环推理
think → act → observe → think → ...
```

其他概念（RAG / MCP / Skill 等等）都是在这三层上加东西，后面会讲。

---

## 二、LLM 基础使用

### 安装依赖

```bash
go get github.com/cloudwego/eino
go get github.com/cloudwego/eino-ext
go get github.com/cloudwego/eino-ext/components/model/ollama
```

### 示例：调用 LLM + 使用模板

```go
package main

import (
	"context"
	"github.com/cloudwego/eino/components/prompt"
	"io"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 使用模版创建messages
	log.Printf("===create messages===\n")
	messages := createMessagesFromTemplate()
	log.Printf("messages: %+v\n\n", messages)

	// 创建llm
	log.Printf("===create llm===\n")
	//cm := createOpenAIChatModel(ctx)
	cm := createOllamaChatModel(ctx)
	log.Printf("create llm success\n\n")

	log.Printf("===llm generate===\n")
	result := generate(ctx, cm, messages)
	log.Printf("result: %+v\n\n", result)

	log.Printf("===llm stream generate===\n")
	streamResult := stream(ctx, cm, messages)
	reportStream(streamResult)
}

func generate(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.Message {
	result, err := llm.Generate(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

func stream(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.StreamReader[*schema.Message] {
	result, err := llm.Stream(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

func createOllamaChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Model:   "deepseek-r1:8b",         // 模型名称
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}
func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()

	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		log.Printf("message[%d]: %+v\n", i, message)
		i++
	}
}

func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
}

func createMessagesFromTemplate() []*schema.Message {
	template := createTemplate()

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	return messages
}
```

- `ChatTemplate`：用变量填充 prompt，避免硬编码

## 三、Tool 调用

### 什么是 Tool？

Tool 本质上是“带参数 schema 的可调用能力”。
在 Go 里它可以由普通函数实现，框架会把它包装成 LLM 可理解的工具描述（如名称、参数、JSON schema）。
LLM 不执行代码，只负责生成“调用哪个工具、传什么参数”。

### 三种写法

```go
// 方式1：InferTool（最简单）
updateTool, err := utils.InferTool("update_todo", "更新一个待办事项", UpdateTodoFunc)

// 方式2：NewTool（手动定义参数 schema）
info := &schema.ToolInfo{
    Name: "add_todo",
    Desc: "添加一个待办事项",
    ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
        "content": {Desc: "内容", Type: schema.String, Required: true},
    }),
}
addTool := utils.NewTool(info, AddTodoFunc)

// 方式3：实现接口（最灵活）
type ListTodoTool struct{}
func (lt *ListTodoTool) Info(_ context.Context) (*schema.ToolInfo, error) { ... }
func (lt *ListTodoTool) InvokableRun(_ context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) { ... }
```

### 带Tool 示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"

	"github.com/cloudwego/eino/schema"
)

func main() {

	ctx := context.Background()

	updateTool, err := utils.InferTool("update_todo", "Update a todo item, eg: content,deadline...", UpdateTodoFunc)
	if err != nil {
		fmt.Printf("InferTool failed, err=%v\n", err)
		return
	}

	// 创建 Google Search 工具
	searchTool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{
		MaxResults: 3,
		Region:     ddgsearch.RegionCN,
		DDGConfig: &ddgsearch.Config{
			Timeout:    10 * time.Second,
			Cache:      true,
			MaxRetries: 5,
		},
	})

	if err != nil {
		fmt.Printf("NewDuckDuckGoTool failed, err=%v", err)
		return
	}

	// 初始化 tools
	todoTools := []tool.BaseTool{
		getAddTodoTool(), // 使用 NewTool 方式
		updateTool,       // 使用 InferTool 方式
		&ListTodoTool{},  // 使用结构体实现方式, 此处未实现底层逻辑
		searchTool,
	}

	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{})
	if err != nil {
		fmt.Printf("NewChatModel failed, err=%v", err)
		return
	}

	// 获取工具信息, 用于绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(todoTools))
	var info *schema.ToolInfo
	for _, todoTool := range todoTools {
		info, err = todoTool.Info(ctx)
		if err != nil {
			fmt.Printf("get ToolInfo failed, err=%v", err)
			return
		}
		toolInfos = append(toolInfos, info)
	}
	// 将 tools 绑定到 ChatModel
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		fmt.Printf("BindTools failed, err=%v", err)
		return
	}

	// 创建 tools 节点
	todoToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: todoTools,
	})
	if err != nil {
		fmt.Printf("NewToolNode failed, err=%v", err)
		return
	}

	// 构建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		fmt.Printf("chain.Compile failed, err=%v", err)
		return
	}

	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "搜索cloudwego的信息",
		},
	})
	if err != nil {
		fmt.Printf("agent.Invoke failed, err=%v", err)
		return
	}

	// 输出结果
	for idx, msg := range resp {
		fmt.Printf("message %d: %s: %s\n", idx, msg.Role, msg.Content)
	}
}

// 获取添加 todo 工具
// 使用 utils.NewTool 创建工具
func getAddTodoTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "add_todo",
		Desc: "Add a todo item",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Desc:     "The content of the todo item",
				Type:     schema.String,
				Required: true,
			},
			"started_at": {
				Desc: "The started time of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
			"deadline": {
				Desc: "The deadline of the todo item, in unix timestamp",
				Type: schema.Integer,
			},
		}),
	}

	return utils.NewTool(info, AddTodoFunc)
}

// ListTodoTool
// 获取列出 todo 工具
// 自行实现 InvokableTool 接口
type ListTodoTool struct{}

func (lt *ListTodoTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "list_todo",
		Desc: "List all todo items",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"finished": {
				Desc:     "filter todo items if finished",
				Type:     schema.Boolean,
				Required: false,
			},
		}),
	}, nil
}

type TodoUpdateParams struct {
	ID        string  `json:"id" jsonschema:"description=id of the todo"`
	Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
	StartedAt *int64  `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
	Deadline  *int64  `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
	Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

type TodoAddParams struct {
	Content  string `json:"content"`
	StartAt  *int64 `json:"started_at,omitempty"` // 开始时间
	Deadline *int64 `json:"deadline,omitempty"`
}

func (lt *ListTodoTool) InvokableRun(_ context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	fmt.Printf("invoke tool list_todo: %s", argumentsInJSON)

	// Tool处理代码
	// ...

	return `{"todos": [{"id": "1", "content": "在2024年12月10日之前完成Eino项目演示文稿的准备工作", "started_at": 1717401600, "deadline": 1717488000, "done": false}]}`, nil
}

func AddTodoFunc(_ context.Context, params *TodoAddParams) (string, error) {
	fmt.Printf("invoke tool add_todo: %+v", params)

	// Tool处理代码
	// ...

	return `{"msg": "add todo success"}`, nil
}

func UpdateTodoFunc(_ context.Context, params *TodoUpdateParams) (string, error) {
	fmt.Printf("invoke tool update_todo: %+v", params)

	// Tool处理代码
	// ...

	return `{"msg": "update todo success"}`, nil
}

```

---

## 四、ReactAgent：循环推理（Loop）

### 和普通 Agent 的区别

普通 Workflow：`model → tools`，只走一遍

ReactAgent：
```
think（LLM 决定下一步）
  → act（调工具）
  → observe（看结果）
  → think（继续推理）
  → ...
  → 最终回答（LLM 决定不再调工具）
```

适合需要多步推理的任务，比如"先查用户信息，再根据信息查课程，再生成推荐"。

### 示例

```go
func main() {
    ctx := context.Background()

    toolableChatModel, err := openai.NewChatModel(...)

    tools := compose.ToolsNodeConfig{
        Tools: []tool.BaseTool{myTool, ...},
    }

    agent, err := react.NewAgent(ctx, &react.AgentConfig{
        ToolCallingModel: toolableChatModel,
        ToolsConfig:      tools,
        MaxStep:          25, // 最多循环25步，防止死循环
    })

    resp, err := agent.Generate(ctx, []*schema.Message{
        schema.UserMessage("帮我查一下用户123的信息，然后推荐适合他的课程"),
    })
}

/*
ReactAgent 什么时候停？
1. LLM 认为已经可以直接回答
2. 达到最大步数 MaxStep
3. 工具调用失败或流程中断
*/
```

---

## 五、RAG：给 Agent 加知识库

### 为什么需要 RAG？

LLM 的知识有截止日期，也不知道你的私有数据（公司文档、课程资料等）。

RAG（检索增强生成）= 先从知识库里找相关内容，再把内容塞进 prompt，让 LLM 参考着回答。

### 流程

```
【数据准备】
文档 → 分块 → embedding model（向量化）→ 存入向量数据库

【查询时】
用户问题 → embedding model→ 向量数据库检索相似内容 → 注入 prompt → LLM 回答
```

### 关键概念

- **Embedding**：把文字转成一串数字（向量），语义相近的文字，向量也相近
- **向量数据库**：存向量、查相似，常用：Milvus、Chroma、ES、Qdrant等等
- **召回率**：检索到的相关内容占全部相关内容的比例，越高越好
- **Rerank**：对召回结果二次排序，提升精度

---

## 六、Memory：给 Agent 加记忆

### 为什么需要 Memory？

LLM 本身无状态，每次请求都是全新的。Memory 让 Agent 记住用户说过的话。

### 三种记忆

| 类型 | 说明 | 实现方式 |
|---|---|---|
| 短期记忆 | 当前对话历史 | 把历史消息拼进 prompt |
| 长期记忆 | 跨会话的用户信息 | 存数据库，按需查询 |
| 用户画像 | 对用户的总结性描述 | 定期从长期记忆提炼 |

### 短期记忆示例

```go
// 每次请求时，把历史消息一起发给 LLM
messages := []*schema.Message{}
messages = append(messages, history...)  // 历史消息
messages = append(messages, schema.UserMessage(userInput))  // 当前消息

resp, _ := chatModel.Generate(ctx, messages)

// 把这轮对话存起来
history = append(history, schema.UserMessage(userInput))
history = append(history, resp)
```

---

## 七、MCP：接入现成工具生态

### 什么是 MCP？

MCP（Model Context Protocol）是标准化的“远程能力暴露协议”（MCP 不只是 Tool 协议，还可以承载更多上下文与能力。但大家最常用到的还是“通过 MCP 接入工具”这一部分）

简单理解：**别人写好了一堆工具（搜索、数据库、GitHub、Slack...），你通过 MCP Server直接接入用，不用自己实现。**

### 和 InferTool 的区别

| | InferTool | MCP |
|---|---|---|
| 适合 | 自己写的业务逻辑 | 复用第三方工具 |
| tool 在哪 | 同进程 Go 函数 | 独立进程，通过协议通信 |
| 复杂度 | 低 | 稍高 |

### MCP 工具资源

- https://github.com/punkpeye/awesome-mcp-servers — GitHub 精选列表
- https://glama.ai/mcp/servers — 精选 MCP 服务集合
- https://mcp.composio.dev/ — 可组合的 MCP 服务平台

---

## 八、Skill：给 Agent 划定能力边界

### 什么是 Skill？

Skill 是对一组相关 Tool 和 Prompt 的封装，代表 Agent 的一项"专项能力"。

类比：Tool 是一把锤子，Skill 是"装修"这件事——它知道什么时候用锤子、什么时候用钻头、先做什么后做什么。

### 和 Tool 的区别

| | Tool | Skill |
|---|---|---|
| 粒度 | 单个函数 | 一组工具 + 流程 |
| 职责 | 执行一个动作 | 完成一类任务 |
| 例子 | `search(query)` | "信息检索"（搜索 + 过滤 + 摘要） |

### 为什么需要 Skill？

当 Agent 的 Tool 越来越多，直接把几十个 Tool 全塞给 LLM 会有两个问题：

1. **Token 浪费**：每次请求都要把所有 Tool 的描述发给 LLM
2. **选择困难**：Tool 太多，LLM 容易选错

Skill 的做法：先判断用户意图属于哪个 Skill，再只把这个 Skill 下的 Tool 暴露给 LLM。

### 在架构中的位置

```
用户消息
    ↓
意图识别 → 匹配 Skill
    ↓
Skill（内含相关 Tools + Prompt）
    ↓
ReactAgent 执行
    ↓
最终回答
```

> Skill 更多是一种架构思想，而非具体 API

---

## 九、完整架构总览

```
用户消息
    ↓
Advisor（意图分类，选择处理流程）
    ↓
ReactAgent
    ├── 短期记忆（对话历史）
    ├── 长期记忆（用户画像、历史事实）
    ├── RAG（课程资料检索）
    └── Tools
         ├── 业务 Tool（InferTool，Go 函数）
         └── 第三方 Tool（MCP）
    ↓
最终回答
```

---

## 十、作业

**基础：** 把第二节和第三节的示例改成 HTTP 服务

**进阶：** 实现一个带记忆的问答 Agent：

1. 支持多轮对话（短期记忆）
2. 注册至少一个自定义 Tool（比如查天气、查时间）

**挑战：** 加入 RAG，让 Agent 能回答基于本地文档的问题。

---

## 十一、延伸阅读

- [Eino 官方文档](https://www.cloudwego.io/zh/docs/eino/)

- [Eino Components](https://www.cloudwego.io/zh/docs/eino/core_modules/components/)

- [learn-claude-code（推荐）](https://github.com/shareAI-lab/learn-claude-code)

  
