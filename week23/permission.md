# 权限判定

常见的2种权限方案:
+ ABAC: Attribute-based Access Control, 基于属性的访问控制
+ RBAC: Role-based Access Control，基于角色的访问控制

keyauth支持2种权限判断方式:
+ 基于用户属性的权限控制, 用于 keyauth本身的权限控制
+ 基于Policy RBAC的权限控制, 基于 对其他子系统的权限控制, 避免子系统权限设计不好, 被提权, 破坏整个权限系统

## keyauth本身权限判定

### 用户分类

首先我们看下用户分类:
```go
// 为了防止越权, 用户可以调整的权限范围只有10已下的权限
type UserType int32

const (
	// 子账号, 无用户中心后台管理权限
	UserType_SUB UserType = 0
	// 组织管理严, 管理该域的组织结构
	UserType_ORG_ADMIN UserType = 5
	// 审计管理员, 可以查看用户中心相关配置, 相当于用户中心只读权限
	UserType_AUDIT_ADMIN UserType = 6
	// 权限管理员, 管理用户的授权策略, 比如空间管理，资源访问策略的配置
	UserType_PERM_ADMIN UserType = 7
	// 域管理员, 有该域的所有管理权限, 协作主账号进行管理
	UserType_DOMAIN_ADMIN UserType = 8
	// 主账号, 具有本域的所有权限
	UserType_PRIMARY UserType = 10
	// 超级管理员, 系统管理员, 万能的人, 不受权限系统约束
	UserType_SUPPER UserType = 15
)
```

### 权限设置

keyauth会和框架集成, 在路由上打上对应的标识, 最终由权限拦截器负责判定

下面是user模块的权限设定
```go
// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	prmary := router.ResourceRouter("primary_account")
	prmary.Allow(types.UserType_SUPPER)
	prmary.BasePath("users")
	prmary.Handle("POST", "/", h.CreatePrimayAccount)
	prmary.Handle("DELETE", "/", h.DestroyPrimaryAccount)

	ram := router.ResourceRouter("ram_account")
	ram.Allow(types.UserType_ORG_ADMIN)
	ram.BasePath("sub_users")
	ram.Handle("POST", "/", h.CreateSubAccount)
	ram.Handle("GET", "/", h.QuerySubAccount)
	ram.Handle("GET", "/:account", h.DescribeSubAccount)
	ram.Handle("PATCH", "/:account", h.PatchSubAccount)
	ram.Handle("DELETE", "/:account", h.DestroySubAccount)
	ram.BasePath("manage")
	ram.Handle("POST", "/block", h.BlockSubAccount)

	portal := router.ResourceRouter("profile")
	portal.BasePath("profile")
	portal.Handle("GET", "/", h.QueryProfile)
	portal.Handle("GET", "/domain", h.QueryDomain)
	portal.Handle("PUT", "/", h.PutProfile)
	portal.Handle("PATCH", "/", h.PatchProfile)

	dom := router.ResourceRouter("domain")
	dom.Allow(types.UserType_DOMAIN_ADMIN)
	dom.BasePath("settings/domain")
	dom.Handle("PUT", "/info", h.UpdateDomainInfo)
	dom.Handle("PUT", "/security", h.UpdateDomainSecurity)

	pass := router.ResourceRouter("password")
	pass.BasePath("password")
	pass.Handle("POST", "/", h.GeneratePassword)
	pass.Handle("PUT", "/", h.UpdatePassword)
}
```

下面是部门管理的权限设定:
```go
	r.BasePath("departments")
	r.Handle("POST", "/", h.Create).SetAllow(types.UserType_ORG_ADMIN)
	r.Handle("GET", "/", h.List)
	r.Handle("GET", "/:id", h.Get)
	r.Handle("PUT", "/:id", h.Put)
	r.Handle("PATCH", "/:id", h.Patch)
	r.Handle("GET", "/:id/subs", h.GetSub)
	r.Handle("DELETE", "/:id", h.Delete).SetAllow(types.UserType_ORG_ADMIN)
```

### 权限判定

用户认证过后, Token信息里面包含当前用户的类型, 只需要拦截, 和权限条目进行匹配就可以判定

```go
func (a *HTTPAuther) ValidatePermission(ctx context.Context, tk *token.Token, e httpb.Entry) error {
	if tk == nil {
		return exception.NewUnauthorized("validate permission need token")
	}

	// 如果是超级管理员不做权限校验, 直接放行
	if tk.UserType.IsIn(types.UserType_SUPPER) {
		a.l.Debugf("[%s] supper admin skip permission check!", tk.Account)
		return nil
	}

	// 检查是否是允许的类型
	if len(e.Allow) > 0 {
		a.l.Debugf("[%s] start check permission to keyauth ...", tk.Account)
		if !e.IsAllow(tk.UserType) {
			return exception.NewPermissionDeny("no permission, allow: %s, but current: %s", e.Allow, tk.UserType)
		}
		a.l.Debugf("[%s] permission check passed", tk.Account)
	}

	return nil
}
```

### 测试

+ 我们先使用超级管理员换取一个token
+ 使用超级管理员创建一个普通用户和一个组织管理员用户
    + 子用户01 ^K%0SFpivmUq
    + 组织管理员01 ^K%0SFpivmUq

然后测试子用户和组织管理员 能否创建用户

## 其他服务权限判定

非keyauth本身的服务, 比如cmdb, 要判断用户是否有权限可以调用怎么办喃?

+ 基于用户类型来判断, 如果这时候新增了一种类型怎么办?

 这种适用于场景固定的服务, 因为我们只使用了一个用户属性做判断


### RBAC权限模型

我们之前的基于用户属性的控制，更偏向于ACL, 而RBAC, 通过角色解绑了 用户直接和权限的关联:

![](./pic/rbac.png)


+ 权限由角色定义, 一个角色拥有多个功能的操作权限
+ 一个人可以有多个角色, 从而获取几个角色的一个功能集合

![](./pic/rbac_rs.jpeg)



### 功能注册

我们的Keyauth就是这样一套中心化的权限系统, 现在用户和角色 都可以在Keyauth定义, 但是怎么知道 服务的功能喃?

这就需要服务把 功能注册到 权限中心

cmdb 添加endpoint注册逻辑:
```go
// registryEndpoints 注册条目
func (s *HTTPService) registryEndpoints() {
	// 注册服务权限条目
	s.l.Info("start registry endpoints ...")

	req := endpoint.NewRegistryRequest(version.Short(), s.r.GetEndpoints().UniquePathEntry())
	_, err := s.kc.Endpoint().RegistryEndpoint(context.Background(), req)
	if err != nil {
		s.l.Warnf("registry endpoints error, %s", err)
	} else {
		s.l.Debug("service endpoints registry success")
	}
}
```

![](./pic/endpoint.png)

注册的时候我们补充了Resource 和 Lables, 我们看看我们何时定义的:
```go
func (h *handler) Registry(r router.SubRouter) {
	hr := r.ResourceRouter("host")
	hr.Handle("GET", "/hosts", h.QueryHost).AddLabel(label.List)
	hr.Handle("POST", "/hosts", h.CreateHost).AddLabel(label.Create)
	hr.Handle("GET", "/hosts/:id", h.DescribeHost).AddLabel(label.Get)
	hr.Handle("DELETE", "/hosts/:id", h.DeleteHost).AddLabel(label.Delete)
	hr.Handle("PUT", "/hosts/:id", h.PutHost).AddLabel(label.Update)
	hr.Handle("PATCH", "/hosts/:id", h.PatchHost).AddLabel(label.Update)
}
```

我们特意添加了Label, 这样我们定义角色的时候，就不用和 具体的功能耦合, 不然功能一变，我们的角色就要重新定义

基于我们注册的Endpoint，我们定义角色:
```
role:          A    
resource:     [host]   
match label:  action: [list, get, create],
```

### 定义角色

创建一个cmdb-reader的角色, 允许访问cmdb所有资源的list和get方法

调用 POST: /keyauth/api/v1/roles 
```json
{
    "name": "cmdb-reader",
    "description": "cmdb 接口只读权限",
    "permissions": [
        {
            "effect": "allow",
            "service_id": "c687sgma0brlmo1a1cag",
            "resource_name": "*",
            "label_key": "action",
            "label_values": [
                "list", "get"
            ]
        }
    ]
}
```

创建一个cmdb-admin的角色, 允许访问cmdb所有资源的所有标签

调用 POST: /keyauth/api/v1/roles 
```json
{
    "name": "cmdb-admin",
    "description": "cmdb 管理员",
    "permissions": [
        {
            "effect": "allow",
            "service_id": "c687sgma0brlmo1a1cag",
            "resource_name": "*",
            "label_key": "*",
            "label_values": [
                "*"
            ]
        }
    ]
}
```

### 策略引擎

单纯依赖rbac也有缺陷，比如:
+ 来了一个外包用户, 需要临时使用你的系统一段时间，直接给他一个角色，又是永久权限
+ 我只允许 用户在 上班时间 才能访问怎么办
+ 我又几条业务线, 我需要允许我业务线的能可以访问这个资源, 当时其他业务线的用户则不可以

基于以上种种问题，rbac不具体这样的灵活性， 因此我们还是使用 PBAC(Policy-based Access Control)

我们需要一个策略来制定用户的访问权限，比如: 用户A 只能在 上班时间 访问业务线A 的资源(比如发布上线):
```
# policy 策略
用户： A
空间:  业务线A
时间:  上班时间
角色:  publisher
```

先创建一个空间: POST /keyauth/api/v1/namesapces
```json
{
    "department_id": ".1",
	"name": "默认空间",
	"description": "默认空间"
}
```

我们调用接口 创建一条策略: POST /keyauth/api/v1/polices/
```json
{
	"namespace_id": "c6auqqea0brvia231k5g",
	"account": "子用户01",
	"role_id": "c6agl7ua0brlc2ie3j0g"
}
```

我们我们使用子用户01来分别访问 list,get,create 来进行测试:

### 权限测试

keyauth提供了一个HTTP接口可以查看用户的权限列表 GET: keyauth/api/v1/namespaces/c6auqqea0brvia231k5g/permissions
```json
{
    "code": 0,
    "data": {
        "total": 0,
        "items": [
            {
                "id": "a1492f4b",
                "role_id": "c6agl7ua0brlc2ie3j0g",
                "create_at": 1637157535069,
                "creater": "admin",
                "effect": "ALLOW",
                "service_id": "c687sgma0brlmo1a1cag",
                "resource_name": "*",
                "label_key": "action",
                "match_all": false,
                "label_values": [
                    "list",
                    "get"
                ],
                "scope": "",
                "desc": ""
            }
        ]
    }
}
```

keyauth提供了专门用户验证用户是否有权限的API: POST keyauth/api/v1/namespaces/c6auqqea0brvia231k5g/permissions

检查用户在namespace: c6auqqea0brvia231k5g 下是否有访问服务: c687sgma0brlmo1a1cag 的GET./cmdb/api/v1/hosts path的权限
```json
{
    "service_id": "c687sgma0brlmo1a1cag",
    "path": "GET./cmdb/api/v1/hosts"
}
```

下面是检测通过, 并返回允许通行的策略信息:
```json
{
    "code": 0,
    "data": {
        "id": "a1492f4b",
        "role_id": "c6agl7ua0brlc2ie3j0g",
        "create_at": 1637157535069,
        "creater": "admin",
        "effect": "ALLOW",
        "service_id": "c687sgma0brlmo1a1cag",
        "resource_name": "*",
        "label_key": "action",
        "match_all": false,
        "label_values": [
            "list",
            "get"
        ],
        "scope": "",
        "desc": ""
    }
}
```

如果我们把 GET换成POST,  鉴权接口就返回鉴权失败:
```json
{
    "code": 403,
    "namespace": "global",
    "reason": "访问未授权",
    "message": "in namespace c6auqqea0brvia231k5g, role [cmdb-reader] has no permission access endpoint: POST./cmdb/api/v1/hosts"
}
```

你可能很奇怪, 怎么多个ID，我怎么记得住喃? 真正对接权限的时候是系统对接, 这些都不要人输入的, 用户只需要提供token即可

## 权限对接

那我们CMDB服务如何接入 这套权限系统喃?


### 权限中间件

我们之前做了认证的中间件, 在用户认证后, 再添加鉴权的逻辑:
```go
func (a *HTTPAuther) Auth(r *http.Request, entry httpb.Entry) (
	authInfo interface{}, err error) {
	var tk *token.Token

	// 从请求中获取access token
	acessToken := r.Header.Get(header.OAuthTokenHeader)

	if entry.AuthEnable {
		ctx := r.Context()

		// 校验身份
		tk, err = a.ValidateIdentity(ctx, acessToken)
		if err != nil {
			return nil, err
		}

		// namesapce检查
		if entry.RequiredNamespace && tk.Namespace == "" {
			return nil, exception.NewBadRequest("namespace required!")
		}

		// 权限检查
		err = a.CheckPermission(ctx, a.mode, tk, entry)
		if err != nil {
			return nil, err
		}
	}

	// 设置RequestID
	if r.Header.Get(header.RequestIdHeader) == "" {
		r.Header.Set(header.RequestIdHeader, xid.New().String())
	}

	return tk, nil
}

func (a *HTTPAuther) CheckPermission(ctx context.Context, mod PermissionCheckMode, tk *token.Token, e httpb.Entry) error {
	if tk == nil {
		return exception.NewUnauthorized("validate permission need token")
	}

	// 如果是超级管理员不做权限校验, 直接放行
	if tk.UserType.IsIn(types.UserType_SUPPER) {
		a.l.Debugf("[%s] supper admin skip permission check!", tk.Account)
		return nil
	}

	switch a.mode {
	case ACL_MODE:
		return a.ValidatePermissionByACL(ctx, tk, e)
	case PRBAC_MODE:
		return a.ValidatePermissionByPRBAC(ctx, tk, e)
	default:
		return fmt.Errorf("only support acl and prbac")
	}

}

func (a *HTTPAuther) ValidatePermissionByACL(ctx context.Context, tk *token.Token, e httpb.Entry) error {
	// 检查是否是允许的类型
	if len(e.Allow) > 0 {
		a.l.Debugf("[%s] start check permission to keyauth ...", tk.Account)
		if !e.IsAllow(tk.UserType) {
			return exception.NewPermissionDeny("no permission, allow: %s, but current: %s", e.Allow, tk.UserType)
		}
		a.l.Debugf("[%s] permission check passed", tk.Account)
	}

	return nil
}

func (a *HTTPAuther) ValidatePermissionByPRBAC(ctx context.Context, tk *token.Token, e httpb.Entry) error {
	svr, err := a.GetClientService(ctx)
	if err != nil {
		return err
	}

	req := permission.NewCheckPermissionRequest()
	req.Account = tk.Account
	req.NamespaceId = tk.Namespace
	req.ServiceId = svr.Id
	req.Path = e.UniquePath()
	_, err = a.keyauth.Permission().CheckPermission(ctx, req)
	if err != nil {
		return exception.NewPermissionDeny(err.Error())
	}
	a.l.Debugf("[%s] permission check passed", tk.Account)
	return nil
}

func (a *HTTPAuther) GetClientService(ctx context.Context) (*micro.Micro, error) {
	if a.svr != nil {
		return a.svr, nil
	}
	a.lock.Lock()
	defer a.lock.Unlock()

	req := micro.NewDescribeServiceRequestWithClientID(a.keyauth.GetClientID())
	ins, err := a.keyauth.Micro().DescribeService(ctx, req)
	if err != nil {
		return nil, err
	}
	a.svr = ins
	return ins, nil
}
```

则有我们cmdb就已经接入了 keyauth的权限系统

### cmdb接口测试

我们使用刚才的用户进行测试, 但是当前用户并没有 进入到具体的某个空间之中, 需要使用API 切换到当前用户的工作空间: POST /keyauth/api/v1/self/tokens/
```json
{
    "namespace": "c6auqqea0brvia231k5g"
}
```

否则开启鉴权后会报参数不够:
```json
{
    "code": 403,
    "namespace": "global",
    "reason": "访问未授权",
    "message": "rpc error: code = Unknown desc = validate param error, namespace required"
}
```

此时我们访问cmdb secret list接口, 权限是正常的:
```json
{
    "code": 0,
    "data": {
        "id": "c5pcffua0bro7e7a05j0",
        "create_at": 1634912191581,
        "description": "腾讯云资源同步",
        "vendor": "TENCENT",
        "allow_regions": [
            "*"
        ],
        "crendential_type": "API_KEY",
        "address": "",
        "api_key": "AKIDc3E9nAG5ewyjESP5VfYBXdS0GI3XkhOG",
        "api_secret": "******",
        "request_rate": 5
    }
}
```

但是如果我们要创建secret: http://127.0.0.1:8060/cmdb/api/v1/secrets， 则会报错:
```json
{
    "code": 403,
    "namespace": "global",
    "reason": "访问未授权",
    "message": "rpc error: code = Unknown desc = in namespace c6auqqea0brvia231k5g, role [cmdb-reader] has no permission access endpoint: POST./cmdb/api/v1/secrets"
}
```

到处我们整套cmdb的权限都对接完成


## 流程解读

cmdb 子系统鉴权流程: 
+ 用户携带自己的令牌 访问服务CMDB的接口
+ cmdb 询问 keyauth 令牌是否合法, 并返回了令牌关联的用户信息
+ cmdb 询问 keyauth 用户是否有访问空间: c6auqqea0brvia231k5g 里面的接口: POST./cmdb/api/v1/secrets 的权限
+ 如果都通过 返回用户访问的资源: secrets

keyauth服务鉴权流程:
+ token validate, 查表 并判定token是否过期
+ permission check
	+ 查询用户在namespace c6auqqea0brvia231k5g 下的关联策略
	+ 查询策略关联的 角色
	+ 判定 当前访问的endpoint 是否能匹配到 当前角色, 通过标签 list，get ... 进行endpint表匹配

```
check roles [cmdb-reader] has permission access endpoint [function_name:"DescribeSecret"  path:"GET./cmdb/api/v1/secrets/:id"  method:"GET"  resource:"secret"  auth_enable:true  permission_enable:true  labels:{key:"action"  value:"get"}]
```
