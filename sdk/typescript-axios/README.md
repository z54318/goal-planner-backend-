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
*GoalsApi* | [**goalDelete**](docs/GoalsApi.md#goaldelete) | **DELETE** /api/goals/{id} | 删除目标
*GoalsApi* | [**goalGet**](docs/GoalsApi.md#goalget) | **GET** /api/goals/{id} | 获取目标详情
*GoalsApi* | [**goalUpdate**](docs/GoalsApi.md#goalupdate) | **PUT** /api/goals/{id} | 更新目标
*GoalsApi* | [**goalUpdateStatus**](docs/GoalsApi.md#goalupdatestatus) | **PATCH** /api/goals/{id}/status | 更新目标状态
*GoalsApi* | [**goalsList**](docs/GoalsApi.md#goalslist) | **GET** /api/goals | 获取目标列表
*MenusApi* | [**adminMenuCreate**](docs/MenusApi.md#adminmenucreate) | **POST** /api/admin/menus | 新增菜单
*MenusApi* | [**adminMenuDelete**](docs/MenusApi.md#adminmenudelete) | **DELETE** /api/admin/menus/{id} | 删除菜单
*MenusApi* | [**adminMenuGet**](docs/MenusApi.md#adminmenuget) | **GET** /api/admin/menus/{id} | 获取菜单详情
*MenusApi* | [**adminMenuUpdate**](docs/MenusApi.md#adminmenuupdate) | **PUT** /api/admin/menus/{id} | 更新菜单
*MenusApi* | [**adminMenusList**](docs/MenusApi.md#adminmenuslist) | **GET** /api/admin/menus | 获取菜单列表
*PlansApi* | [**goalPlanGenerate**](docs/PlansApi.md#goalplangenerate) | **POST** /api/goals/{id}/generate-plan | 生成目标计划
*PlansApi* | [**goalPlanGet**](docs/PlansApi.md#goalplanget) | **GET** /api/goals/{id}/plan | 获取目标计划
*PlansApi* | [**goalPlanRegenerate**](docs/PlansApi.md#goalplanregenerate) | **POST** /api/goals/{id}/regenerate-plan | 重新生成目标计划
*RbacApi* | [**adminPermissionCreate**](docs/RbacApi.md#adminpermissioncreate) | **POST** /api/admin/permissions | 新增权限
*RbacApi* | [**adminPermissionDelete**](docs/RbacApi.md#adminpermissiondelete) | **DELETE** /api/admin/permissions/{id} | 删除权限
*RbacApi* | [**adminPermissionUpdate**](docs/RbacApi.md#adminpermissionupdate) | **PUT** /api/admin/permissions/{id} | 更新权限
*RbacApi* | [**adminPermissionsList**](docs/RbacApi.md#adminpermissionslist) | **GET** /api/admin/permissions | 获取权限列表
*RbacApi* | [**adminRolePermissionsGet**](docs/RbacApi.md#adminrolepermissionsget) | **GET** /api/admin/roles/{id}/permissions | 获取角色权限绑定
*RbacApi* | [**adminRolePermissionsUpdate**](docs/RbacApi.md#adminrolepermissionsupdate) | **PUT** /api/admin/roles/{id}/permissions | 更新角色权限绑定
*RbacApi* | [**adminRolesList**](docs/RbacApi.md#adminroleslist) | **GET** /api/admin/roles | 获取角色列表
*TasksApi* | [**taskCreate**](docs/TasksApi.md#taskcreate) | **POST** /api/tasks | 新增任务
*TasksApi* | [**taskDelete**](docs/TasksApi.md#taskdelete) | **DELETE** /api/tasks/{id} | 删除任务
*TasksApi* | [**taskGet**](docs/TasksApi.md#taskget) | **GET** /api/tasks/{id} | 获取任务详情
*TasksApi* | [**taskUpdate**](docs/TasksApi.md#taskupdate) | **PUT** /api/tasks/{id} | 编辑任务
*TasksApi* | [**taskUpdateStatus**](docs/TasksApi.md#taskupdatestatus) | **PATCH** /api/tasks/{id}/status | 更新任务状态
*TasksApi* | [**tasksList**](docs/TasksApi.md#taskslist) | **GET** /api/tasks | 获取任务列表
*UsersApi* | [**adminUserRolesUpdate**](docs/UsersApi.md#adminuserrolesupdate) | **PUT** /api/admin/users/{id}/roles | 更新用户角色
*UsersApi* | [**adminUsersList**](docs/UsersApi.md#adminuserslist) | **GET** /api/admin/users | 获取用户列表


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
 - [GoalGoalStatus](docs/GoalGoalStatus.md)
 - [GoalUpdateGoalRequest](docs/GoalUpdateGoalRequest.md)
 - [GoalUpdateGoalStatusRequest](docs/GoalUpdateGoalStatusRequest.md)
 - [MenuCreateMenuRequest](docs/MenuCreateMenuRequest.md)
 - [MenuMenu](docs/MenuMenu.md)
 - [MenuMenuListResponse](docs/MenuMenuListResponse.md)
 - [MenuMenuResponse](docs/MenuMenuResponse.md)
 - [MenuUpdateMenuRequest](docs/MenuUpdateMenuRequest.md)
 - [PlanPhase](docs/PlanPhase.md)
 - [PlanPlan](docs/PlanPlan.md)
 - [PlanPlanResponse](docs/PlanPlanResponse.md)
 - [PlanTask](docs/PlanTask.md)
 - [PlanTaskPriority](docs/PlanTaskPriority.md)
 - [PlanTaskStatus](docs/PlanTaskStatus.md)
 - [RbacCreatePermissionRequest](docs/RbacCreatePermissionRequest.md)
 - [RbacPermission](docs/RbacPermission.md)
 - [RbacPermissionListResponse](docs/RbacPermissionListResponse.md)
 - [RbacPermissionResponse](docs/RbacPermissionResponse.md)
 - [RbacRole](docs/RbacRole.md)
 - [RbacRoleListResponse](docs/RbacRoleListResponse.md)
 - [RbacRolePermissionIDsData](docs/RbacRolePermissionIDsData.md)
 - [RbacRolePermissionIDsResponse](docs/RbacRolePermissionIDsResponse.md)
 - [RbacUpdatePermissionRequest](docs/RbacUpdatePermissionRequest.md)
 - [RbacUpdateRolePermissionsRequest](docs/RbacUpdateRolePermissionsRequest.md)
 - [ResponseBody](docs/ResponseBody.md)
 - [ResponseErrorBody](docs/ResponseErrorBody.md)
 - [TaskCreateTaskRequest](docs/TaskCreateTaskRequest.md)
 - [TaskDeleteTaskData](docs/TaskDeleteTaskData.md)
 - [TaskDeleteTaskResponse](docs/TaskDeleteTaskResponse.md)
 - [TaskTask](docs/TaskTask.md)
 - [TaskTaskListData](docs/TaskTaskListData.md)
 - [TaskTaskListResponse](docs/TaskTaskListResponse.md)
 - [TaskTaskPriority](docs/TaskTaskPriority.md)
 - [TaskTaskResponse](docs/TaskTaskResponse.md)
 - [TaskTaskStatus](docs/TaskTaskStatus.md)
 - [TaskUpdateTaskRequest](docs/TaskUpdateTaskRequest.md)
 - [TaskUpdateTaskStatusRequest](docs/TaskUpdateTaskStatusRequest.md)
 - [UserRole](docs/UserRole.md)
 - [UserUpdateUserRolesRequest](docs/UserUpdateUserRolesRequest.md)
 - [UserUser](docs/UserUser.md)
 - [UserUserListResponse](docs/UserUserListResponse.md)
 - [UserUserResponse](docs/UserUserResponse.md)


<a id="documentation-for-authorization"></a>
## Documentation For Authorization


Authentication schemes defined for the API:
<a id="BearerAuth"></a>
### BearerAuth

- **Type**: API key
- **API key parameter name**: Authorization
- **Location**: HTTP header

