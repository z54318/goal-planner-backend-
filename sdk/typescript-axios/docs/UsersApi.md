# UsersApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminUserRolesUpdate**](#adminuserrolesupdate) | **PUT** /api/admin/users/{id}/roles | 更新用户角色|
|[**adminUsersList**](#adminuserslist) | **GET** /api/admin/users | 获取用户列表|

# **adminUserRolesUpdate**
> UserUserResponse adminUserRolesUpdate(request)


### Example

```typescript
import {
    UsersApi,
    Configuration,
    UserUpdateUserRolesRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let id: number; //用户ID (default to undefined)
let request: UserUpdateUserRolesRequest; //用户角色参数

const { status, data } = await apiInstance.adminUserRolesUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **UserUpdateUserRolesRequest**| 用户角色参数 | |
| **id** | [**number**] | 用户ID | defaults to undefined|


### Return type

**UserUserResponse**

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

# **adminUsersList**
> UserUserListResponse adminUsersList()


### Example

```typescript
import {
    UsersApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

const { status, data } = await apiInstance.adminUsersList();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**UserUserListResponse**

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

