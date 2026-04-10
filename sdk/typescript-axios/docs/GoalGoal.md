# GoalGoal


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**aggregate_status** | [**GoalGoalStatus**](GoalGoalStatus.md) | AggregateStatus 按任务聚合出来的目标状态 | [optional] [default to undefined]
**category** | **string** | Category 目标分类 | [optional] [default to undefined]
**created_at** | **string** | CreatedAt 创建时间 | [optional] [default to undefined]
**description** | **string** | Description 目标描述 | [optional] [default to undefined]
**id** | **number** | ID 目标ID | [optional] [default to undefined]
**status** | [**GoalGoalStatus**](GoalGoalStatus.md) | Status 目标状态 | [optional] [default to undefined]
**target_deadline** | **string** | TargetDeadline 截止时间 | [optional] [default to undefined]
**title** | **string** | Title 目标标题 | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt 更新时间 | [optional] [default to undefined]
**user_id** | **number** | UserID 用户ID | [optional] [default to undefined]

## Example

```typescript
import { GoalGoal } from '@goal-planner/backend-sdk';

const instance: GoalGoal = {
    aggregate_status,
    category,
    created_at,
    description,
    id,
    status,
    target_deadline,
    title,
    updated_at,
    user_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
