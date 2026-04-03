## @goal-planner/backend-sdk@1.0.0

This generator creates TypeScript/JavaScript client that utilizes [axios](https://github.com/axios/axios). The generated Node module can be used in the following environments:

Environment
* Node.js
* Webpack
* Browserify

Language level
* ES5 - you must have a Promises/A+ library installed
* ES6

Module system
* CommonJS
* ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition will be automatically resolved via `package.json`. ([Reference](https://www.typescriptlang.org/docs/handbook/declaration-files/consumption.html))

### Building

To build and compile the typescript sources to javascript use:
```
npm install
npm run build
```

### Publishing

First build the package then run `npm publish`

### Consuming

navigate to the folder of your consuming project and run one of the following commands.

_published:_

```
npm install @goal-planner/backend-sdk@1.0.0 --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
```

### Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AuthApi* | [**authLogin**](docs/AuthApi.md#authlogin) | **POST** /api/auth/login | 用户登录
*AuthApi* | [**authMenus**](docs/AuthApi.md#authmenus) | **GET** /api/auth/menus | 获取当前登录用户可见菜单
*AuthApi* | [**authProfile**](docs/AuthApi.md#authprofile) | **GET** /api/auth/profile | 获取当前登录用户信息
*AuthApi* | [**authRegister**](docs/AuthApi.md#authregister) | **POST** /api/auth/register | 用户注册
*GoalsApi* | [**goalCreate**](docs/GoalsApi.md#goalcreate) | **POST** /api/goals | 创建目标
*GoalsApi* | [**goalGet**](docs/GoalsApi.md#goalget) | **GET** /api/goals/{id} | 获取目标详情
*GoalsApi* | [**goalsList**](docs/GoalsApi.md#goalslist) | **GET** /api/goals | 获取目标列表


### Documentation For Models

 - [AuthLoginData](docs/AuthLoginData.md)
 - [AuthLoginRequest](docs/AuthLoginRequest.md)
 - [AuthLoginResponse](docs/AuthLoginResponse.md)
 - [AuthMenu](docs/AuthMenu.md)
 - [AuthMenusResponse](docs/AuthMenusResponse.md)
 - [AuthProfileData](docs/AuthProfileData.md)
 - [AuthProfileResponse](docs/AuthProfileResponse.md)
 - [AuthRegisterData](docs/AuthRegisterData.md)
 - [AuthRegisterRequest](docs/AuthRegisterRequest.md)
 - [AuthRegisterResponse](docs/AuthRegisterResponse.md)
 - [GoalCreateGoalRequest](docs/GoalCreateGoalRequest.md)
 - [GoalGoal](docs/GoalGoal.md)
 - [GoalGoalListData](docs/GoalGoalListData.md)
 - [GoalGoalListResponse](docs/GoalGoalListResponse.md)
 - [GoalGoalResponse](docs/GoalGoalResponse.md)
 - [ResponseErrorBody](docs/ResponseErrorBody.md)


<a id="documentation-for-authorization"></a>
## Documentation For Authorization


Authentication schemes defined for the API:
<a id="BearerAuth"></a>
### BearerAuth

- **Type**: API key
- **API key parameter name**: Authorization
- **Location**: HTTP header

