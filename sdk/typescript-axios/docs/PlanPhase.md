# PlanPhase


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** | CreatedAt 创建时间 | [optional] [default to undefined]
**description** | **string** | Description 阶段描述 | [optional] [default to undefined]
**id** | **number** | ID 阶段ID | [optional] [default to undefined]
**plan_id** | **number** | PlanID 所属计划ID | [optional] [default to undefined]
**sort_order** | **number** | SortOrder 阶段顺序 | [optional] [default to undefined]
**tasks** | [**Array&lt;PlanTask&gt;**](PlanTask.md) | Tasks 任务列表 | [optional] [default to undefined]
**title** | **string** | Title 阶段标题 | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt 更新时间 | [optional] [default to undefined]

## Example

```typescript
import { PlanPhase } from '@goal-planner/backend-sdk';

const instance: PlanPhase = {
    created_at,
    description,
    id,
    plan_id,
    sort_order,
    tasks,
    title,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
