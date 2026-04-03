# AuthApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**authLogin**](#authlogin) | **POST** /api/auth/login | 用户登录|
|[**authMenus**](#authmenus) | **GET** /api/auth/menus | 获取当前登录用户可见菜单|
|[**authProfile**](#authprofile) | **GET** /api/auth/profile | 获取当前登录用户信息|
|[**authRegister**](#authregister) | **POST** /api/auth/register | 用户注册|

# **authLogin**
> AuthLoginResponse authLogin(request)


### Example

```typescript
import {
    AuthApi,
    Configuration,
    AuthLoginRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: AuthLoginRequest; //登录参数

const { status, data } = await apiInstance.authLogin(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **AuthLoginRequest**| 登录参数 | |


### Return type

**AuthLoginResponse**

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
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authMenus**
> AuthMenusResponse authMenus()


### Example

```typescript
import {
    AuthApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.authMenus();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**AuthMenusResponse**

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

# **authProfile**
> AuthProfileResponse authProfile()


### Example

```typescript
import {
    AuthApi,
    Configuration
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.authProfile();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**AuthProfileResponse**

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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authRegister**
> AuthRegisterResponse authRegister(request)


### Example

```typescript
import {
    AuthApi,
    Configuration,
    AuthRegisterRequest
} from '@goal-planner/backend-sdk';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: AuthRegisterRequest; //注册参数

const { status, data } = await apiInstance.authRegister(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **AuthRegisterRequest**| 注册参数 | |


### Return type

**AuthRegisterResponse**

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

