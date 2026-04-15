package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	// ErrNotConfigured 表示 AI 能力未配置。
	ErrNotConfigured = errors.New("ai service not configured")
	// ErrInvalidResponse 表示 AI 返回内容不符合预期。
	ErrInvalidResponse = errors.New("ai service returned invalid response")
)

// RequestError 表示 AI 服务请求失败。
type RequestError struct {
	StatusCode int
	Body       string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("ai service request failed: status=%d body=%s", e.StatusCode, e.Body)
}

// GoalInput 表示生成计划时需要的目标信息。
type GoalInput struct {
	Title          string
	Description    string
	Category       string
	TargetDeadline *time.Time
}

// PlanOutput 表示 AI 生成出的计划内容。
type PlanOutput struct {
	Title    string        `json:"title"`
	Overview string        `json:"overview"`
	Phases   []PhaseOutput `json:"phases"`
}

// PhaseOutput 表示 AI 生成出的阶段内容。
type PhaseOutput struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Order       int          `json:"order"`
	Tasks       []TaskOutput `json:"tasks"`
}

// TaskOutput 表示 AI 生成出的任务内容。
type TaskOutput struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	EstimatedDays int    `json:"estimated_days"`
	Deliverables  string `json:"deliverables"`
	Deadline      string `json:"deadline"`
	Priority      string `json:"priority"`
	Order         int    `json:"order"`
}

// NextStepSuggestion 表示 AI 给出的执行建议。
type NextStepSuggestion struct {
	Summary    string   `json:"summary"`
	NextAction string   `json:"next_action"`
	Reason     string   `json:"reason"`
	Checklist  []string `json:"checklist"`
	Risk       string   `json:"risk"`
}

// PlanSuggestionInput 表示计划执行建议所需上下文。
type PlanSuggestionInput struct {
	GoalTitle       string                `json:"goal_title"`
	GoalDescription string                `json:"goal_description"`
	GoalStatus      string                `json:"goal_status"`
	PlanTitle       string                `json:"plan_title"`
	PlanOverview    string                `json:"plan_overview"`
	Phases          []PlanSuggestionPhase `json:"phases"`
}

// PlanSuggestionPhase 表示计划执行建议中的阶段摘要。
type PlanSuggestionPhase struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	TodoCount       int    `json:"todo_count"`
	InProgressCount int    `json:"in_progress_count"`
	DoneCount       int    `json:"done_count"`
}

// PhaseSuggestionInput 表示阶段执行建议所需上下文。
type PhaseSuggestionInput struct {
	GoalTitle        string                 `json:"goal_title"`
	PlanTitle        string                 `json:"plan_title"`
	PhaseTitle       string                 `json:"phase_title"`
	PhaseDescription string                 `json:"phase_description"`
	Tasks            []SuggestionTaskDigest `json:"tasks"`
}

// TaskSuggestionInput 表示任务执行建议所需上下文。
type TaskSuggestionInput struct {
	GoalTitle    string                 `json:"goal_title"`
	PlanTitle    string                 `json:"plan_title"`
	PhaseTitle   string                 `json:"phase_title"`
	Task         SuggestionTaskDigest   `json:"task"`
	SiblingTasks []SuggestionTaskDigest `json:"sibling_tasks"`
}

// SuggestionTaskDigest 表示执行建议上下文中的任务摘要。
type SuggestionTaskDigest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Deadline    string `json:"deadline,omitempty"`
}

type planSkeletonOutput struct {
	Title    string        `json:"title"`
	Overview string        `json:"overview"`
	Phases   []PhaseOutput `json:"phases"`
}

type phaseTasksOutput struct {
	Tasks []TaskOutput `json:"tasks"`
}

// Client 表示一个兼容 OpenAI Chat Completions 的 AI 客户端。
type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewClient 创建 AI 客户端。
func NewClient(apiKey string, baseURL string, model string) *Client {
	return &Client{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		model:   strings.TrimSpace(model),
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// GeneratePlan 根据目标生成一份计划。
func (c *Client) GeneratePlan(ctx context.Context, goal GoalInput) (PlanOutput, error) {
	if c.apiKey == "" || c.baseURL == "" || c.model == "" {
		return PlanOutput{}, ErrNotConfigured
	}

	deadline := "未设置"
	if goal.TargetDeadline != nil {
		deadline = goal.TargetDeadline.Format(time.RFC3339)
	}

	skeletonPrompt := fmt.Sprintf(
		"目标标题：%s\n目标描述：%s\n目标分类：%s\n截止时间：%s\n请先输出计划骨架 JSON。",
		goal.Title,
		valueOrDefault(goal.Description, "未填写"),
		valueOrDefault(goal.Category, "未分类"),
		deadline,
	)

	skeletonContent, err := c.chatCompletion(
		ctx,
		"你是一个计划拆解助手。请根据用户目标先生成计划骨架。你只能返回 JSON，且必须包含 title、overview、phases 三个字段。限制：phases 最多 3 个；phase.title 不超过 24 个汉字；phase.description 不超过 80 个汉字；phase 只包含 title、description、order，不要包含 tasks。overview 使用中文，控制在 100 字内。阶段要体现明显的先后顺序。不要返回任何解释文字或 Markdown。",
		skeletonPrompt,
	)
	if err != nil {
		return PlanOutput{}, err
	}

	var skeleton planSkeletonOutput
	if err := json.Unmarshal([]byte(skeletonContent), &skeleton); err != nil {
		return PlanOutput{}, ErrInvalidResponse
	}

	output := PlanOutput{
		Title:    strings.TrimSpace(skeleton.Title),
		Overview: strings.TrimSpace(skeleton.Overview),
		Phases:   skeleton.Phases,
	}
	if output.Title == "" || output.Overview == "" || len(output.Phases) == 0 {
		return PlanOutput{}, ErrInvalidResponse
	}

	output.Title = limitRunes(output.Title, 50)
	output.Overview = limitRunes(output.Overview, 120)
	if len(output.Phases) > 3 {
		output.Phases = output.Phases[:3]
	}

	for i := range output.Phases {
		output.Phases[i].Title = strings.TrimSpace(output.Phases[i].Title)
		output.Phases[i].Description = strings.TrimSpace(output.Phases[i].Description)
		output.Phases[i].Tasks = nil
		if output.Phases[i].Title == "" {
			return PlanOutput{}, ErrInvalidResponse
		}
		output.Phases[i].Title = limitRunes(output.Phases[i].Title, 24)
		output.Phases[i].Description = limitRunes(output.Phases[i].Description, 80)
		if output.Phases[i].Order <= 0 {
			output.Phases[i].Order = i + 1
		}

		phaseTaskPrompt := fmt.Sprintf(
			"目标标题：%s\n阶段标题：%s\n阶段描述：%s\n请仅为这个阶段生成任务 JSON。",
			output.Title,
			output.Phases[i].Title,
			valueOrDefault(output.Phases[i].Description, "未填写"),
		)
		taskContent, err := c.chatCompletion(
			ctx,
			"你是一个任务拆解助手。请只为单个阶段生成任务。你只能返回 JSON，且必须包含 tasks 字段。tasks 最多 3 个。每个 task 必须包含 title、description、estimated_days、deliverables、deadline、priority、order。限制：task.title 不超过 24 个汉字；task.description 不超过 60 个汉字；deliverables 不超过 24 个汉字；estimated_days 取 1 到 7 的整数；deadline 使用 RFC3339 时间字符串；priority 只能是 high、medium、low 之一。任务要具体、可操作、可验证。不要返回任何解释文字或 Markdown。",
			phaseTaskPrompt,
		)
		if err != nil {
			return PlanOutput{}, err
		}

		var taskOutput phaseTasksOutput
		if err := json.Unmarshal([]byte(taskContent), &taskOutput); err != nil {
			return PlanOutput{}, ErrInvalidResponse
		}
		if len(taskOutput.Tasks) == 0 {
			return PlanOutput{}, ErrInvalidResponse
		}
		if len(taskOutput.Tasks) > 3 {
			taskOutput.Tasks = taskOutput.Tasks[:3]
		}

		for j := range taskOutput.Tasks {
			taskOutput.Tasks[j].Title = strings.TrimSpace(taskOutput.Tasks[j].Title)
			taskOutput.Tasks[j].Description = strings.TrimSpace(taskOutput.Tasks[j].Description)
			taskOutput.Tasks[j].Deliverables = strings.TrimSpace(taskOutput.Tasks[j].Deliverables)
			taskOutput.Tasks[j].Priority = strings.TrimSpace(strings.ToLower(taskOutput.Tasks[j].Priority))
			if taskOutput.Tasks[j].Title == "" {
				return PlanOutput{}, ErrInvalidResponse
			}
			taskOutput.Tasks[j].Title = limitRunes(taskOutput.Tasks[j].Title, 24)
			taskOutput.Tasks[j].Description = limitRunes(taskOutput.Tasks[j].Description, 60)
			taskOutput.Tasks[j].Deliverables = limitRunes(taskOutput.Tasks[j].Deliverables, 24)
			if taskOutput.Tasks[j].Order <= 0 {
				taskOutput.Tasks[j].Order = j + 1
			}
			if taskOutput.Tasks[j].EstimatedDays <= 0 {
				taskOutput.Tasks[j].EstimatedDays = 1
			}
			if taskOutput.Tasks[j].EstimatedDays > 7 {
				taskOutput.Tasks[j].EstimatedDays = 7
			}
			if _, err := time.Parse(time.RFC3339, taskOutput.Tasks[j].Deadline); err != nil {
				taskOutput.Tasks[j].Deadline = ""
			}
			if !isValidTaskPriority(taskOutput.Tasks[j].Priority) {
				taskOutput.Tasks[j].Priority = "medium"
			}
		}

		applyDefaultTaskDeadlines(taskOutput.Tasks, goal.TargetDeadline)
		output.Phases[i].Tasks = taskOutput.Tasks
	}

	return output, nil
}

// SuggestNextStepForPlan 根据计划上下文生成执行建议。
func (c *Client) SuggestNextStepForPlan(ctx context.Context, input PlanSuggestionInput) (NextStepSuggestion, error) {
	return c.suggestNextStep(
		ctx,
		"你是一个项目执行助手。请根据用户当前的计划状态给出执行建议。你只能返回 JSON，且必须包含 summary、next_action、reason、checklist、risk 五个字段。summary 用一句中文概括当前最值得执行的重点；next_action 给出最具体的一步执行动作；reason 说明为什么现在应该执行这一步；checklist 返回 2 到 4 个可执行检查项；risk 用一句中文说明如果不执行这一步的主要风险。不要返回任何解释文字或 Markdown。",
		input,
	)
}

// SuggestNextStepForPhase 根据阶段上下文生成执行建议。
func (c *Client) SuggestNextStepForPhase(ctx context.Context, input PhaseSuggestionInput) (NextStepSuggestion, error) {
	return c.suggestNextStep(
		ctx,
		"你是一个阶段执行助手。请根据用户当前阶段和阶段内任务状态给出执行建议。你只能返回 JSON，且必须包含 summary、next_action、reason、checklist、risk 五个字段。summary 用一句中文概括当前阶段最该推进的重点；next_action 给出最具体的一步执行动作；reason 说明为什么现在应该执行这一步；checklist 返回 2 到 4 个可执行检查项；risk 用一句中文说明如果不执行这一步的主要风险。不要返回任何解释文字或 Markdown。",
		input,
	)
}

// SuggestNextStepForTask 根据任务上下文生成执行建议。
func (c *Client) SuggestNextStepForTask(ctx context.Context, input TaskSuggestionInput) (NextStepSuggestion, error) {
	return c.suggestNextStep(
		ctx,
		"你是一个任务执行助手。请根据当前任务和同阶段上下文给出执行建议。你只能返回 JSON，且必须包含 summary、next_action、reason、checklist、risk 五个字段。summary 用一句中文概括当前任务最该执行的重点；next_action 给出最具体的一步执行动作；reason 说明为什么现在应该执行这一步；checklist 返回 2 到 4 个可执行检查项；risk 用一句中文说明如果不执行这一步的主要风险。不要返回任何解释文字或 Markdown。",
		input,
	)
}

func isValidTaskPriority(priority string) bool {
	switch priority {
	case "high", "medium", "low":
		return true
	default:
		return false
	}
}

func (c *Client) suggestNextStep(ctx context.Context, systemPrompt string, input any) (NextStepSuggestion, error) {
	if c.apiKey == "" || c.baseURL == "" || c.model == "" {
		return NextStepSuggestion{}, ErrNotConfigured
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return NextStepSuggestion{}, err
	}

	content, err := c.chatCompletion(
		ctx,
		systemPrompt,
		"请基于以下 JSON 上下文给出建议：\n"+string(payload),
	)
	if err != nil {
		return NextStepSuggestion{}, err
	}

	var suggestion NextStepSuggestion
	if err := json.Unmarshal([]byte(content), &suggestion); err != nil {
		return NextStepSuggestion{}, ErrInvalidResponse
	}

	normalizeSuggestion(&suggestion)
	if suggestion.Summary == "" || suggestion.NextAction == "" || suggestion.Reason == "" {
		return NextStepSuggestion{}, ErrInvalidResponse
	}

	return suggestion, nil
}

func (c *Client) chatCompletion(ctx context.Context, systemPrompt string, userPrompt string) (string, error) {
	payload := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &RequestError{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(respBody)),
		}
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", ErrInvalidResponse
	}
	if len(chatResp.Choices) == 0 {
		return "", ErrInvalidResponse
	}

	return extractJSON(strings.TrimSpace(chatResp.Choices[0].Message.Content)), nil
}

func valueOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func extractJSON(content string) string {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}

func limitRunes(value string, max int) string {
	value = strings.TrimSpace(value)
	if max <= 0 || utf8.RuneCountInString(value) <= max {
		return value
	}

	runes := []rune(value)
	return string(runes[:max])
}

func applyDefaultTaskDeadlines(tasks []TaskOutput, goalDeadline *time.Time) {
	now := time.Now()
	current := now.Add(24 * time.Hour)
	for i := range tasks {
		if tasks[i].Deadline != "" {
			if parsed, err := time.Parse(time.RFC3339, tasks[i].Deadline); err == nil {
				if parsed.Before(now) {
					tasks[i].Deadline = ""
				} else {
					if goalDeadline != nil && parsed.After(*goalDeadline) {
						parsed = *goalDeadline
					}
					if parsed.Before(current) {
						tasks[i].Deadline = ""
					} else {
						current = parsed
						tasks[i].Deadline = current.Format(time.RFC3339)
						continue
					}
				}
			}
		}

		current = current.Add(time.Duration(tasks[i].EstimatedDays) * 24 * time.Hour)
		if goalDeadline != nil && current.After(*goalDeadline) {
			current = *goalDeadline
		}
		tasks[i].Deadline = current.Format(time.RFC3339)
	}
}

func normalizeSuggestion(suggestion *NextStepSuggestion) {
	suggestion.Summary = limitRunes(strings.TrimSpace(suggestion.Summary), 80)
	suggestion.NextAction = limitRunes(strings.TrimSpace(suggestion.NextAction), 120)
	suggestion.Reason = limitRunes(strings.TrimSpace(suggestion.Reason), 180)
	suggestion.Risk = limitRunes(strings.TrimSpace(suggestion.Risk), 120)

	if len(suggestion.Checklist) > 4 {
		suggestion.Checklist = suggestion.Checklist[:4]
	}

	items := make([]string, 0, len(suggestion.Checklist))
	for _, item := range suggestion.Checklist {
		item = limitRunes(strings.TrimSpace(item), 80)
		if item == "" {
			continue
		}
		items = append(items, item)
	}
	suggestion.Checklist = items
}
