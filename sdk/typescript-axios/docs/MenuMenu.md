# MenuMenu


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**children** | [**Array&lt;MenuMenu&gt;**](MenuMenu.md) | Children 子菜单 | [optional] [default to undefined]
**component** | **string** | Component 组件路径 | [optional] [default to undefined]
**created_at** | **string** | CreatedAt 创建时间 | [optional] [default to undefined]
**hidden** | **boolean** | Hidden 是否隐藏 | [optional] [default to undefined]
**icon** | **string** | Icon 菜单图标 | [optional] [default to undefined]
**id** | **number** | ID 菜单ID | [optional] [default to undefined]
**name** | **string** | Name 菜单名称 | [optional] [default to undefined]
**parent_id** | **number** | ParentID 父菜单ID | [optional] [default to undefined]
**path** | **string** | Path 路由路径 | [optional] [default to undefined]
**permission_code** | **string** | PermissionCode 权限码 | [optional] [default to undefined]
**sort_order** | **number** | SortOrder 排序值 | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt 更新时间 | [optional] [default to undefined]

## Example

```typescript
import { MenuMenu } from '@goal-planner/backend-sdk';

const instance: MenuMenu = {
    children,
    component,
    created_at,
    hidden,
    icon,
    id,
    name,
    parent_id,
    path,
    permission_code,
    sort_order,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
