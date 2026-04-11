# PlanTask


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** | CreatedAt 创建时间 | [optional] [default to undefined]
**deadline** | **string** | Deadline 任务截止时间 | [optional] [default to undefined]
**deliverables** | **string** | Deliverables 交付物 | [optional] [default to undefined]
**description** | **string** | Description 任务描述 | [optional] [default to undefined]
**estimated_days** | **number** | EstimatedDays 预估天数 | [optional] [default to undefined]
**id** | **number** | ID 任务ID | [optional] [default to undefined]
**phase_id** | **number** | PhaseID 所属阶段ID | [optional] [default to undefined]
**priority** | [**PlanTaskPriority**](PlanTaskPriority.md) | Priority 任务优先级 | [optional] [default to undefined]
**sort_order** | **number** | SortOrder 任务顺序 | [optional] [default to undefined]
**status** | [**PlanTaskStatus**](PlanTaskStatus.md) | Status 任务状态 | [optional] [default to undefined]
**title** | **string** | Title 任务标题 | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt 更新时间 | [optional] [default to undefined]

## Example

```typescript
import { PlanTask } from '@goal-planner/backend-sdk';

const instance: PlanTask = {
    created_at,
    deadline,
    deliverables,
    description,
    estimated_days,
    id,
    phase_id,
    priority,
    sort_order,
    status,
    title,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
