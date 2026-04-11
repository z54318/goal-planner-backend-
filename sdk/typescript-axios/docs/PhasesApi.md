# PhasesApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**phaseGet**](#phaseget) | **GET** /api/phases/{id} | 获取阶段详情|
|[**phaseUpdate**](#phaseupdate) | **PUT** /api/phases/{id} | 编辑阶段|

# **phaseGet**
> PhasePhaseResponse phaseGet()


### Example

```typescript
import {
    PhasesApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new PhasesApi(configuration);

let id: number; //阶段ID (default to undefined)

const { status, data } = await apiInstance.phaseGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 阶段ID | defaults to undefined|


### Return type

**PhasePhaseResponse**

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

# **phaseUpdate**
> PhasePhaseResponse phaseUpdate(request)


### Example

```typescript
import {
    PhasesApi,
    Configuration,
    PhaseUpdatePhaseRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new PhasesApi(configuration);

let id: number; //阶段ID (default to undefined)
let request: PhaseUpdatePhaseRequest; //阶段信息

const { status, data } = await apiInstance.phaseUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **PhaseUpdatePhaseRequest**| 阶段信息 | |
| **id** | [**number**] | 阶段ID | defaults to undefined|


### Return type

**PhasePhaseResponse**

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

