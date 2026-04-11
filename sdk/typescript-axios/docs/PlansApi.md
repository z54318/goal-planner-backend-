# PlansApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**goalPlanDelete**](#goalplandelete) | **DELETE** /api/goals/{id}/plan | 删除目标计划|
|[**goalPlanGenerate**](#goalplangenerate) | **POST** /api/goals/{id}/generate-plan | 生成目标计划|
|[**goalPlanGet**](#goalplanget) | **GET** /api/goals/{id}/plan | 获取目标计划|
|[**goalPlanRegenerate**](#goalplanregenerate) | **POST** /api/goals/{id}/regenerate-plan | 重新生成目标计划|
|[**goalPlanUpdate**](#goalplanupdate) | **PUT** /api/goals/{id}/plan | 编辑目标计划|

# **goalPlanDelete**
> ResponseBody goalPlanDelete()


### Example

```typescript
import {
    PlansApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new PlansApi(configuration);

let id: number; //目标ID (default to undefined)

const { status, data } = await apiInstance.goalPlanDelete(
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

# **goalPlanGenerate**
> PlanPlanResponse goalPlanGenerate()


### Example

```typescript
import {
    PlansApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new PlansApi(configuration);

let id: number; //目标ID (default to undefined)

const { status, data } = await apiInstance.goalPlanGenerate(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 目标ID | defaults to undefined|


### Return type

**PlanPlanResponse**

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
|**404** | Not Found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal Server Error |  -  |
|**503** | Service Unavailable |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **goalPlanGet**
> PlanPlanResponse goalPlanGet()


### Example

```typescript
import {
    PlansApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new PlansApi(configuration);

let id: number; //目标ID (default to undefined)

const { status, data } = await apiInstance.goalPlanGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 目标ID | defaults to undefined|


### Return type

**PlanPlanResponse**

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

# **goalPlanRegenerate**
> PlanPlanResponse goalPlanRegenerate()


### Example

```typescript
import {
    PlansApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new PlansApi(configuration);

let id: number; //目标ID (default to undefined)

const { status, data } = await apiInstance.goalPlanRegenerate(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 目标ID | defaults to undefined|


### Return type

**PlanPlanResponse**

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
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |
|**503** | Service Unavailable |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **goalPlanUpdate**
> PlanPlanResponse goalPlanUpdate(request)


### Example

```typescript
import {
    PlansApi,
    Configuration,
    PlanUpdatePlanRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new PlansApi(configuration);

let id: number; //目标ID (default to undefined)
let request: PlanUpdatePlanRequest; //计划参数

const { status, data } = await apiInstance.goalPlanUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **PlanUpdatePlanRequest**| 计划参数 | |
| **id** | [**number**] | 目标ID | defaults to undefined|


### Return type

**PlanPlanResponse**

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

