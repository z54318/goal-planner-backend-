# PlanPlan


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** | CreatedAt 创建时间 | [optional] [default to undefined]
**goal_id** | **number** | GoalID 目标ID | [optional] [default to undefined]
**id** | **number** | ID 计划ID | [optional] [default to undefined]
**overview** | **string** | Overview 计划概述 | [optional] [default to undefined]
**phases** | [**Array&lt;PlanPhase&gt;**](PlanPhase.md) | Phases 阶段列表 | [optional] [default to undefined]
**title** | **string** | Title 计划标题 | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt 更新时间 | [optional] [default to undefined]
**user_id** | **number** | UserID 用户ID | [optional] [default to undefined]

## Example

```typescript
import { PlanPlan } from '@goal-planner/backend-sdk';

const instance: PlanPlan = {
    created_at,
    goal_id,
    id,
    overview,
    phases,
    title,
    updated_at,
    user_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
