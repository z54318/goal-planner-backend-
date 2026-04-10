# UserUser


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** | CreatedAt 创建时间 | [optional] [default to undefined]
**email** | **string** | Email 邮箱 | [optional] [default to undefined]
**id** | **number** | ID 用户ID | [optional] [default to undefined]
**nickname** | **string** | Nickname 昵称 | [optional] [default to undefined]
**role_ids** | **Array&lt;number&gt;** | RoleIDs 已绑定角色ID | [optional] [default to undefined]
**roles** | [**Array&lt;UserRole&gt;**](UserRole.md) | Roles 已绑定角色列表 | [optional] [default to undefined]
**status** | **string** | Status 用户状态 | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt 更新时间 | [optional] [default to undefined]
**username** | **string** | Username 用户名 | [optional] [default to undefined]

## Example

```typescript
import { UserUser } from '@goal-planner/backend-sdk';

const instance: UserUser = {
    created_at,
    email,
    id,
    nickname,
    role_ids,
    roles,
    status,
    updated_at,
    username,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
