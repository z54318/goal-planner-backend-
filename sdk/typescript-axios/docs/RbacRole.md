# RbacRole


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**code** | **string** | Code 角色编码 | [optional] [default to undefined]
**created_at** | **string** | CreatedAt 创建时间 | [optional] [default to undefined]
**id** | **number** | ID 角色ID | [optional] [default to undefined]
**name** | **string** | Name 角色名称 | [optional] [default to undefined]
**permission_ids** | **Array&lt;number&gt;** | PermissionIDs 当前角色绑定的权限ID列表 | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt 更新时间 | [optional] [default to undefined]

## Example

```typescript
import { RbacRole } from '@goal-planner/backend-sdk';

const instance: RbacRole = {
    code,
    created_at,
    id,
    name,
    permission_ids,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
