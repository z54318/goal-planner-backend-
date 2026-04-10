# goal-planner-backend-

## 本地 AI 配置

当前项目默认按 Ollama 本地 OpenAI 兼容接口运行。

环境变量示例：

```env
AI_API_KEY=ollama
AI_BASE_URL=http://127.0.0.1:11434/v1
AI_MODEL=llama3.2
```

本地准备步骤：

1. 安装并启动 Ollama
2. 拉取模型

```powershell
ollama pull llama3.2
```

3. 启动 API

```powershell
go run ./cmd/api
```

计划生成接口：

```text
POST /api/goals/{id}/generate-plan
GET /api/goals/{id}/plan
```
