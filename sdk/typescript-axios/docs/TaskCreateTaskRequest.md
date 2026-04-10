# TaskCreateTaskRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**deliverables** | **string** | Deliverables 交付物 | [optional] [default to undefined]
**description** | **string** | Description 任务描述 | [optional] [default to undefined]
**estimated_days** | **number** | EstimatedDays 预估天数 | [optional] [default to undefined]
**phase_id** | **number** | PhaseID 所属阶段ID | [optional] [default to undefined]
**priority** | [**TaskTaskPriority**](TaskTaskPriority.md) | Priority 任务优先级：high 高，medium 中，low 低 | [optional] [default to undefined]
**sort_order** | **number** | SortOrder 任务顺序 | [optional] [default to undefined]
**title** | **string** | Title 任务标题 | [optional] [default to undefined]

## Example

```typescript
import { TaskCreateTaskRequest } from '@goal-planner/backend-sdk';

const instance: TaskCreateTaskRequest = {
    deliverables,
    description,
    estimated_days,
    phase_id,
    priority,
    sort_order,
    title,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
