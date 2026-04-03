# GoalsApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**goalCreate**](#goalcreate) | **POST** /api/goals | 创建目标|
|[**goalGet**](#goalget) | **GET** /api/goals/{id} | 获取目标详情|
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

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
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

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
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

const { status, data } = await apiInstance.goalsList();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**GoalGoalListResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

