# GoalsApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**goalCreate**](#goalcreate) | **POST** /api/goals | 创建目标|
|[**goalDelete**](#goaldelete) | **DELETE** /api/goals/{id} | 删除目标|
|[**goalGet**](#goalget) | **GET** /api/goals/{id} | 获取目标详情|
|[**goalUpdate**](#goalupdate) | **PUT** /api/goals/{id} | 更新目标|
|[**goalUpdateStatus**](#goalupdatestatus) | **PATCH** /api/goals/{id}/status | 更新目标状态|
|[**goalsList**](#goalslist) | **GET** /api/goals | 获取目标列表|

# **goalCreate**
> GoalGoalResponse goalCreate(request)


### Example

```typescript
import {
    GoalsApi,
    Configuration,
    GoalCreateGoalRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new GoalsApi(configuration);

let request: GoalCreateGoalRequest; //目标参数

const { status, data } = await apiInstance.goalCreate(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **GoalCreateGoalRequest**| 目标参数 | |


### Return type

**GoalGoalResponse**

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

# **goalDelete**
> ResponseBody goalDelete()


### Example

```typescript
import {
    GoalsApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new GoalsApi(configuration);

let id: number; //目标ID (default to undefined)

const { status, data } = await apiInstance.goalDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 目标ID | defaults to undefined|


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

# **goalGet**
> GoalGoalResponse goalGet()


### Example

```typescript
import {
    GoalsApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new GoalsApi(configuration);

let id: number; //目标ID (default to undefined)

const { status, data } = await apiInstance.goalGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 目标ID | defaults to undefined|


### Return type

**GoalGoalResponse**

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

# **goalUpdate**
> GoalGoalResponse goalUpdate(request)


### Example

```typescript
import {
    GoalsApi,
    Configuration,
    GoalUpdateGoalRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new GoalsApi(configuration);

let id: number; //目标ID (default to undefined)
let request: GoalUpdateGoalRequest; //目标参数

const { status, data } = await apiInstance.goalUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **GoalUpdateGoalRequest**| 目标参数 | |
| **id** | [**number**] | 目标ID | defaults to undefined|


### Return type

**GoalGoalResponse**

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

# **goalUpdateStatus**
> GoalGoalResponse goalUpdateStatus(request)


### Example

```typescript
import {
    GoalsApi,
    Configuration,
    GoalUpdateGoalStatusRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new GoalsApi(configuration);

let id: number; //目标ID (default to undefined)
let request: GoalUpdateGoalStatusRequest; //目标状态参数

const { status, data } = await apiInstance.goalUpdateStatus(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **GoalUpdateGoalStatusRequest**| 目标状态参数 | |
| **id** | [**number**] | 目标ID | defaults to undefined|


### Return type

**GoalGoalResponse**

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

# **goalsList**
> GoalGoalListResponse goalsList()


### Example

```typescript
import {
    GoalsApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new GoalsApi(configuration);

let page: number; //页码，从1开始 (optional) (default to undefined)
let pageSize: number; //每页条数 (optional) (default to undefined)

const { status, data } = await apiInstance.goalsList(
    page,
    pageSize
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | 页码，从1开始 | (optional) defaults to undefined|
| **pageSize** | [**number**] | 每页条数 | (optional) defaults to undefined|


### Return type

**GoalGoalListResponse**

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

