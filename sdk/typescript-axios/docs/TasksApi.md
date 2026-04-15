# TasksApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**phaseTasksSort**](#phasetaskssort) | **PUT** /api/phases/{id}/tasks/sort | 调整阶段任务顺序|
|[**taskCreate**](#taskcreate) | **POST** /api/tasks | 新增任务|
|[**taskDelete**](#taskdelete) | **DELETE** /api/tasks/{id} | 删除任务|
|[**taskGet**](#taskget) | **GET** /api/tasks/{id} | 获取任务详情|
|[**taskNextStepGet**](#tasknextstepget) | **GET** /api/tasks/{id}/next-step | 查询任务执行建议|
|[**taskNextStepSuggest**](#tasknextstepsuggest) | **POST** /api/tasks/{id}/next-step | 生成任务执行建议|
|[**taskUpdate**](#taskupdate) | **PUT** /api/tasks/{id} | 编辑任务|
|[**taskUpdateStatus**](#taskupdatestatus) | **PATCH** /api/tasks/{id}/status | 更新任务状态|
|[**tasksList**](#taskslist) | **GET** /api/tasks | 获取任务列表|

# **phaseTasksSort**
> ResponseBody phaseTasksSort(request)


### Example

```typescript
import {
    TasksApi,
    Configuration,
    TaskSortTasksRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let id: number; //阶段ID (default to undefined)
let request: TaskSortTasksRequest; //任务排序信息

const { status, data } = await apiInstance.phaseTasksSort(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **TaskSortTasksRequest**| 任务排序信息 | |
| **id** | [**number**] | 阶段ID | defaults to undefined|


### Return type

**ResponseBody**

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

# **taskCreate**
> TaskTaskResponse taskCreate(request)


### Example

```typescript
import {
    TasksApi,
    Configuration,
    TaskCreateTaskRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let request: TaskCreateTaskRequest; //任务信息

const { status, data } = await apiInstance.taskCreate(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **TaskCreateTaskRequest**| 任务信息 | |


### Return type

**TaskTaskResponse**

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

# **taskDelete**
> TaskDeleteTaskResponse taskDelete()


### Example

```typescript
import {
    TasksApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let id: number; //任务ID (default to undefined)

const { status, data } = await apiInstance.taskDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 任务ID | defaults to undefined|


### Return type

**TaskDeleteTaskResponse**

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

# **taskGet**
> TaskTaskResponse taskGet()


### Example

```typescript
import {
    TasksApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let id: number; //任务ID (default to undefined)

const { status, data } = await apiInstance.taskGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 任务ID | defaults to undefined|


### Return type

**TaskTaskResponse**

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

# **taskNextStepGet**
> TaskNextStepSuggestionResponse taskNextStepGet()


### Example

```typescript
import {
    TasksApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let id: number; //任务ID (default to undefined)

const { status, data } = await apiInstance.taskNextStepGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 任务ID | defaults to undefined|


### Return type

**TaskNextStepSuggestionResponse**

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

# **taskNextStepSuggest**
> TaskNextStepSuggestionResponse taskNextStepSuggest()


### Example

```typescript
import {
    TasksApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let id: number; //任务ID (default to undefined)

const { status, data } = await apiInstance.taskNextStepSuggest(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | 任务ID | defaults to undefined|


### Return type

**TaskNextStepSuggestionResponse**

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
|**502** | Bad Gateway |  -  |
|**503** | Service Unavailable |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **taskUpdate**
> TaskTaskResponse taskUpdate(request)


### Example

```typescript
import {
    TasksApi,
    Configuration,
    TaskUpdateTaskRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let id: number; //任务ID (default to undefined)
let request: TaskUpdateTaskRequest; //任务信息

const { status, data } = await apiInstance.taskUpdate(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **TaskUpdateTaskRequest**| 任务信息 | |
| **id** | [**number**] | 任务ID | defaults to undefined|


### Return type

**TaskTaskResponse**

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

# **taskUpdateStatus**
> TaskTaskResponse taskUpdateStatus(request)


### Example

```typescript
import {
    TasksApi,
    Configuration,
    TaskUpdateTaskStatusRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let id: number; //任务ID (default to undefined)
let request: TaskUpdateTaskStatusRequest; //任务状态

const { status, data } = await apiInstance.taskUpdateStatus(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **TaskUpdateTaskStatusRequest**| 任务状态 | |
| **id** | [**number**] | 任务ID | defaults to undefined|


### Return type

**TaskTaskResponse**

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

# **tasksList**
> TaskTaskListResponse tasksList()


### Example

```typescript
import {
    TasksApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new TasksApi(configuration);

let status: string; //任务状态 (optional) (default to undefined)
let goalId: number; //目标ID (optional) (default to undefined)
let phaseId: number; //阶段ID (optional) (default to undefined)
let page: number; //页码，从1开始 (optional) (default to undefined)
let pageSize: number; //每页条数 (optional) (default to undefined)

const { status, data } = await apiInstance.tasksList(
    status,
    goalId,
    phaseId,
    page,
    pageSize
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **status** | [**string**] | 任务状态 | (optional) defaults to undefined|
| **goalId** | [**number**] | 目标ID | (optional) defaults to undefined|
| **phaseId** | [**number**] | 阶段ID | (optional) defaults to undefined|
| **page** | [**number**] | 页码，从1开始 | (optional) defaults to undefined|
| **pageSize** | [**number**] | 每页条数 | (optional) defaults to undefined|


### Return type

**TaskTaskListResponse**

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

