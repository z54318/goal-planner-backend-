# RbacApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminPermissionCreate**](#adminpermissioncreate) | **POST** /api/admin/permissions | 新增权限|
|[**adminPermissionDelete**](#adminpermissiondelete) | **DELETE** /api/admin/permissions/{id} | 删除权限|
|[**adminPermissionUpdate**](#adminpermissionupdate) | **PUT** /api/admin/permissions/{id} | 更新权限|
|[**adminPermissionsList**](#adminpermissionslist) | **GET** /api/admin/permissions | 获取权限列表|
|[**adminRolePermissionsGet**](#adminrolepermissionsget) | **GET** /api/admin/roles/{id}/permissions | 获取角色权限绑定|
|[**adminRolePermissionsUpdate**](#adminrolepermissionsupdate) | **PUT** /api/admin/roles/{id}/permissions | 更新角色权限绑定|
|[**adminRolesList**](#adminroleslist) | **GET** /api/admin/roles | 获取角色列表|

# **adminPermissionCreate**
> RbacPermissionResponse adminPermissionCreate(request)


### Example

```typescript
import {
    RbacApi,
    Configuration,
    RbacCreatePermissionRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new RbacApi(configuration);

let request: RbacCreatePermissionRequest; //权限参数

const { status, data } = await apiInstance.adminPermissionCreate(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **RbacCreatePermissionRequest**| 权限参数 | |


### Return type

**RbacPermissionResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminPermissionDelete**
> ResponseBody adminPermissionDelete()


### Example

```typescript
import {
    RbacApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new RbacApi(configuration);

let id: number; //权限ID (default to undefined)

const { status, data } = await apiInstance.adminPermissionDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 权限ID | defaults to undefined|


### Return type

**ResponseBody**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminPermissionUpdate**
> RbacPermissionResponse adminPermissionUpdate(request)


### Example

```typescript
import {
    RbacApi,
    Configuration,
    RbacUpdatePermissionRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new RbacApi(configuration);

let id: number; //权限ID (default to undefined)
let request: RbacUpdatePermissionRequest; //权限参数

const { status, data } = await apiInstance.adminPermissionUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **RbacUpdatePermissionRequest**| 权限参数 | |
| **id** | [**number**] | 权限ID | defaults to undefined|


### Return type

**RbacPermissionResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminPermissionsList**
> RbacPermissionListResponse adminPermissionsList()


### Example

```typescript
import {
    RbacApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new RbacApi(configuration);

const { status, data } = await apiInstance.adminPermissionsList();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**RbacPermissionListResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminRolePermissionsGet**
> RbacRolePermissionIDsResponse adminRolePermissionsGet()


### Example

```typescript
import {
    RbacApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new RbacApi(configuration);

let id: number; //角色ID (default to undefined)

const { status, data } = await apiInstance.adminRolePermissionsGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 角色ID | defaults to undefined|


### Return type

**RbacRolePermissionIDsResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminRolePermissionsUpdate**
> RbacRolePermissionIDsResponse adminRolePermissionsUpdate(request)


### Example

```typescript
import {
    RbacApi,
    Configuration,
    RbacUpdateRolePermissionsRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new RbacApi(configuration);

let id: number; //角色ID (default to undefined)
let request: RbacUpdateRolePermissionsRequest; //角色权限参数

const { status, data } = await apiInstance.adminRolePermissionsUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **RbacUpdateRolePermissionsRequest**| 角色权限参数 | |
| **id** | [**number**] | 角色ID | defaults to undefined|


### Return type

**RbacRolePermissionIDsResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminRolesList**
> RbacRoleListResponse adminRolesList()


### Example

```typescript
import {
    RbacApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new RbacApi(configuration);

const { status, data } = await apiInstance.adminRolesList();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**RbacRoleListResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

