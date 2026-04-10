# MenusApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminMenuCreate**](#adminmenucreate) | **POST** /api/admin/menus | 新增菜单|
|[**adminMenuDelete**](#adminmenudelete) | **DELETE** /api/admin/menus/{id} | 删除菜单|
|[**adminMenuGet**](#adminmenuget) | **GET** /api/admin/menus/{id} | 获取菜单详情|
|[**adminMenuUpdate**](#adminmenuupdate) | **PUT** /api/admin/menus/{id} | 更新菜单|
|[**adminMenusList**](#adminmenuslist) | **GET** /api/admin/menus | 获取菜单列表|

# **adminMenuCreate**
> MenuMenuResponse adminMenuCreate(request)


### Example

```typescript
import {
    MenusApi,
    Configuration,
    MenuCreateMenuRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new MenusApi(configuration);

let request: MenuCreateMenuRequest; //菜单参数

const { status, data } = await apiInstance.adminMenuCreate(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **MenuCreateMenuRequest**| 菜单参数 | |


### Return type

**MenuMenuResponse**

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

# **adminMenuDelete**
> ResponseBody adminMenuDelete()


### Example

```typescript
import {
    MenusApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new MenusApi(configuration);

let id: number; //菜单ID (default to undefined)

const { status, data } = await apiInstance.adminMenuDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 菜单ID | defaults to undefined|


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

# **adminMenuGet**
> MenuMenuResponse adminMenuGet()


### Example

```typescript
import {
    MenusApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new MenusApi(configuration);

let id: number; //菜单ID (default to undefined)

const { status, data } = await apiInstance.adminMenuGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 菜单ID | defaults to undefined|


### Return type

**MenuMenuResponse**

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

# **adminMenuUpdate**
> MenuMenuResponse adminMenuUpdate(request)


### Example

```typescript
import {
    MenusApi,
    Configuration,
    MenuUpdateMenuRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new MenusApi(configuration);

let id: number; //菜单ID (default to undefined)
let request: MenuUpdateMenuRequest; //菜单参数

const { status, data } = await apiInstance.adminMenuUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **MenuUpdateMenuRequest**| 菜单参数 | |
| **id** | [**number**] | 菜单ID | defaults to undefined|


### Return type

**MenuMenuResponse**

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

# **adminMenusList**
> MenuMenuListResponse adminMenusList()


### Example

```typescript
import {
    MenusApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new MenusApi(configuration);

const { status, data } = await apiInstance.adminMenusList();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**MenuMenuListResponse**

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

