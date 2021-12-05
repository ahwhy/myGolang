# API Server 订阅SCM事件

我们回到 API Server部分的开发，补充应用管理模块, 然后应用关联仓库地址

![](./images/pipeline-ci.png)

## 接口定义

定义应用，主要是scm相关信息

```protobuf
// SCM_TYPE 源码仓库类型
enum SCM_TYPE {
    // gitlab
    GITLAB = 0;
	// github
	GITHUB = 1;
	// coding.net
	CODING = 2;
}

// Application todo
message Application {
    // 唯一ID
    // @gotags: bson:"_id" json:"id"
    string id = 1;
    // 用于加密应用的铭感信息
    // @gotags: bson:"key" json:"key"
    string key = 19;
    // 所属域
    // @gotags: bson:"domain" json:"domain"
    string domain = 2;
    // 所属空间
    // @gotags: bson:"namespace" json:"namespace"
    string namespace = 3;
    // 创建时间
    // @gotags: bson:"create_at" json:"create_at"
    int64 create_at = 4;
    // 创建人
    // @gotags: bson:"create_by" json:"create_by"
    string create_by = 5;
    // 更新时间
    // @gotags: bson:"update_at" json:"update_at"
    int64 update_at = 6;
    // 更新人
    // @gotags: bson:"update_by" json:"update_by"
    string update_by = 7;
    // 名称
    // @gotags: bson:"name" json:"name"
    string name = 8;
    // 应用标签
    // @gotags: bson:"tags" json:"tags"
    map<string, string> tags = 9;
    // 描述
    // @gotags: bson:"description" json:"description"
    string description = 10;
    // 仓库ssh url地址
    // @gotags: json:"repo_ssh_url"
    string repo_ssh_url = 12;
    // 仓库http url地址
    // @gotags: json:"repo_http_url"
    string repo_http_url = 13;
    // 仓库来源类型
    // @gotags: json:"scm_type"
    SCM_TYPE scm_type = 14;
    // 仓库来源类型
    // @gotags: json:"scm_project_id"
    string scm_project_id = 15;
    // scm设置Hook后返回的id, 用于删除应用时，取消hook使用
    // @gotags: json:"scm_hook_id"
    string scm_hook_id = 16;
    // 创建hook过程中的错误信息
    // @gotags: json:"hook_error"
    string hook_error = 17;
    // 仓库的priviate token, 用于设置回调
    // @gotags: json:"scm_private_token"
    string scm_private_token = 18;
    // 用于创建pipeline的请求参数
    // @gotags: json:"pipelines"
    repeated Pipeline pipelines = 11;
}
```

接口 Application 对象的CRUD和一个应用事件处理接口(用于处理回调事件的, 比如SCM commit)
```protobuf
service Service {
    // 应用管理
	rpc CreateApplication(CreateApplicationRequest) returns(Application);
    rpc UpdateApplication(UpdateApplicationRequest) returns(Application);
	rpc QueryApplication(QueryApplicationRequest) returns(ApplicationSet);
    rpc DescribeApplication(DescribeApplicationRequest) returns(Application);
    rpc DeleteApplication(DeleteApplicationRequest) returns(Application);
    // 应用事件处理
    rpc HandleApplicationEvent(ApplicationEvent) returns(Application);
}
```

## 应用同步(gitlab)

我们应用的代码可能在:
+ gitlab
+ github
+ gitee
+ coding

这里以gitlab为例, 我们如何查询我们有哪些仓库?

可以参考下官方文档:
+  [Project List API 接口数据](https://gitlab.com/api/v4/projects?owned=true)
+  [Project API 文档](https://docs.gitlab.com/ce/api/projects.html)

因此整体思路是, 用户通过自己的private token, 拉去自己的仓库, 

### 封装 SDK

下面是在封装Gitlab API

```go
package gitlab

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func NewSCM(addr, token string) *SCM {
	return &SCM{
		Address:      addr,
		PrivateToken: token,
		Version:      "v4",
		client:       &http.Client{Timeout: 5 * time.Second},
	}
}

type SCM struct {
	Address      string
	PrivateToken string
	Version      string

	client *http.Client
}

func (r *SCM) newJSONRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", r.PrivateToken)
	return req, nil
}

func (r SCM) newFormReqeust(method, url string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("PRIVATE-TOKEN", r.PrivateToken)
	return req, nil
}

func (r *SCM) resourceURL(resource string, params map[string]string) string {
	val := make(url.Values)

	for k, v := range params {
		val.Set(k, v)
	}

	return fmt.Sprintf("%s/api/%s/%s?%s", r.Address, r.Version, resource, val.Encode())
}
```

下面通过API文档封装出来的查看项目列表的方法:
```go
// https://gitlab.com/api/v4/projects?owned=true
// https://docs.gitlab.com/ce/api/projects.html
func (r *SCM) ListProjects() (*scm.ProjectSet, error) {
	projectURL := r.resourceURL("projects", map[string]string{"owned": "true", "simple": "true"})
	req, err := r.newJSONRequest("GET", projectURL)
	if err != nil {
		return nil, err
	}

	// 发起请求
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取body
	bytesB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respString := string(bytesB)

	if (resp.StatusCode / 100) != 2 {
		return nil, fmt.Errorf("status code[%d] is not 200, response %s", resp.StatusCode, respString)
	}

	set := NewProjectSet()
	if err := json.Unmarshal(bytesB, &set.Items); err != nil {
		return nil, err
	}

	return set, nil
}
```

### 测试SDK

这里直接使用的gitlab的公有云服务: https://gitlab.com/, 直接去他那儿创建一个Private token, 个人token代表就是个人的身份, 因此可以看到自己的项目，并且有所有权限

![](./images/gitlab-pt.png)

然后自己创建一个项目, 这里我已经创建了一个

![](./images/create-git-repo.png)

编写测试用例测试:
```go
var (
	GitLabAddr    = "https://gitlab.com"
	PraviateToken = ""
)

func TestListProject(t *testing.T) {
	should := assert.New(t)

	repo := gitlab.NewSCM(GitLabAddr, PraviateToken)
	ps, err := repo.ListProjects()
	should.NoError(err)
	fmt.Println(ps)
}
```

使用自己的token测试如下:
```sh
items:{id:29032549  description:"测试使用"  name:"sample-devcloud"  git_ssh_url:"git@gitlab.com:yumaojun03/sample-devcloud.git"  git_http_url:"https://gitlab.com/yumaojun03/sample-devcloud.git"  namespace_path:"yumaojun03/sample-devcloud"}
```

### 封装API

我们能查询出的信息可以作为基础信息, 提供一个API给前端, 这样用户基于同步的信息在补充下就可以创建一个应用了

```go
const (
	GitlabEventHeaderKey = "X-Gitlab-Event"
	GitlabEventTokenKey  = "X-Gitlab-Token"
)

func (h *handler) QuerySCMProject(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	srcType := qs.Get("scm_type")

	var (
		ps  *scm.ProjectSet
		err error
	)
	switch srcType {
	case "gitlab", "":
		repo := gitlab.NewSCM(qs.Get("scm_addr"), qs.Get("token"))
		ps, err = repo.ListProjects()
	case "github":
	case "coding":
	default:
		response.Failed(w, exception.NewBadRequest("unknown scm_type %s", srcType))
		return
	}

	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ps.Items)
}
```

暴露HTTP 接口
```go
	r.BasePath("repo/projects")
	r.Handle("GET", "/", h.QuerySCMProject)
```


## 回调设置与取消

光同步过来还没用, 需要gitlab有仓库有变更的时候通知我们, 因此需要在gitlab那边设置一个webhook, 地址就是我们workflow的地址

![](./images/gitlab-webhook.png)

人为设置容易出错, 既然用户连Private token都给了(这里其实最安全的是使用Oauth2.0来做, 但是毕竟麻烦), 我们直接使用API 操作, 代替用户把Hook配置好

因此我们补充一对方法:
+ 添加Hook(当我们创建应用时，需要添加)
+ 取消Hook(当我们删除应用时，需要取消)

这里我们定义的回调处理API为:
```go
r.BasePath("triggers/scm/gitlab")
r.Handle("POST", "/", h.GitLabHookHanler).DisableAuth()
```

我们只需要将回调http://devcloud.nbtuan.vip/workflow/api/v1/triggers/scm/gitlab, 即可, 注意这里的 http://devcloud.nbtuan.vip 是需要通过配置文件配置的, 而且必须是公网地址, 所以我们本地是调试不通的

```go
// POST /projects/:id/hooks
// https://docs.gitlab.com/ce/api/projects.html#add-project-hook
func (r *SCM) AddProjectHook(in *AddProjectHookRequest) (*AddProjectHookResponse, error) {
	addHookURL := r.resourceURL(fmt.Sprintf("projects/%d/hooks", in.ProjectID), nil)
	req, err := r.newFormReqeust("POST", addHookURL, strings.NewReader(in.Hook.FormValue().Encode()))
	if err != nil {
		return nil, err
	}

	// 发起请求
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取body
	bytesB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respString := string(bytesB)

	if (resp.StatusCode / 100) != 2 {
		return nil, fmt.Errorf("status code[%d] is not 200, response %s", resp.StatusCode, respString)
	}

	ins := NewAddProjectHookResponse()
	if err := json.Unmarshal(bytesB, &ins); err != nil {
		return nil, err
	}

	return ins, nil
}

func NewDeleteProjectReqeust(projectID, hookID int64) *DeleteProjectReqeust {
	return &DeleteProjectReqeust{
		ProjectID: projectID,
		HookID:    hookID,
	}
}

type DeleteProjectReqeust struct {
	ProjectID int64
	HookID    int64
}

// DELETE /projects/:id/hooks/:hook_id
// https://docs.gitlab.com/ce/api/projects.html#delete-project-hook
func (r *SCM) DeleteProjectHook(in *DeleteProjectReqeust) error {
	addHookURL := r.resourceURL(fmt.Sprintf("projects/%d/hooks/%d", in.ProjectID, in.HookID), nil)
	req, err := r.newFormReqeust("DELETE", addHookURL, nil)
	if err != nil {
		return err
	}

	// 发起请求
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取body
	bytesB, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respString := string(bytesB)

	if (resp.StatusCode / 100) != 2 {
		return fmt.Errorf("status code[%d] is not 200, response %s", resp.StatusCode, respString)
	}

	return nil
}
```

然后我们在创建时设置钩子
```go
func (s *service) CreateApplication(ctx context.Context, req *application.CreateApplicationRequest) (
	*application.Application, error) {
	ins, err := application.NewApplication(req)
	if err != nil {
		return nil, err
	}

	hookId, err := s.setWebHook(req, ins.GenWebHook(s.platform))
	if err != nil {
		ins.HookError = fmt.Sprintf("add web hook error, %s", err)
	}
	ins.ScmHookId = hookId

	if _, err := s.col.InsertOne(context.TODO(), ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a application document error, %s", err)
	}

	return ins, nil
}

func (s *service) setWebHook(req *application.CreateApplicationRequest, hook *gitlab.WebHook) (string, error) {
	if !req.NeedSetHook() {
		return "", nil
	}

	addr, err := req.GetScmAddr()
	if err != nil {
		return "", fmt.Errorf("get scm addr from http_url error, %s", err)
	}

	repo := gitlab.NewSCM(addr, req.ScmPrivateToken)
	addHookReq := gitlab.NewAddProjectHookRequest(req.Int64ScmProjectID(), hook)
	resp, err := repo.AddProjectHook(addHookReq)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", resp.ID), nil
}
```

删除应用时，删除钩子
```go
func (s *service) DeleteApplication(ctx context.Context, req *application.DeleteApplicationRequest) (
	*application.Application, error) {
	ins, err := s.DescribeApplication(ctx, application.NewDescribeApplicationRequestWithName(req.Namespace, req.Name))
	if err != nil {
		return nil, err
	}

	// 删除Hook
	if err := s.delWebHook(ins); err != nil {
		s.log.Errorf("delete scm hook error, %s", err)
	}

	if _, err := s.col.DeleteOne(context.TODO(), bson.M{"_id": ins.Id}); err != nil {
		return nil, err
	}
	return ins, nil
}

func (s *service) delWebHook(req *application.Application) error {
	if req.ScmHookId == "" {
		return nil
	}

	if req.ScmPrivateToken == "" {
		s.log.Errorf("delete scm hook error, scm_private_token is empty")
		return nil
	}

	addr, err := req.GetScmAddr()
	if err != nil {
		return fmt.Errorf("get scm addr from http_url error, %s", err)
	}

	repo := gitlab.NewSCM(addr, req.ScmPrivateToken)
	delHookReq := gitlab.NewDeleteProjectReqeust(req.Int64ScmProjectID(), req.Int64ScmHookID())

	return repo.DeleteProjectHook(delHookReq)
}
```

## 创建应用

然后我们创建一个应用, 关联事件上流水线, 当有事件时, 直接由应用创建Pipeline, 这样ci的流程就走通了



## 处理Hook事件

我们测试下, 然后拿到一份样例数据:
```json
{
    "object_kind": "push",
    "event_name": "push",
    "before": "f8a831144634f5810e17014582b5ba21267bb257",
    "after": "f8a831144634f5810e17014582b5ba21267bb257",
    "ref": "refs/heads/main",
    "checkout_sha": "f8a831144634f5810e17014582b5ba21267bb257",
    "message": null,
    "user_id": 9556442,
    "user_name": "紫川秀",
    "user_username": "yumaojun03",
    "user_email": "",
    "user_avatar": "https://secure.gravatar.com/avatar/1c8f622795d244227b2982871bc925d6?s=80&d=identicon",
    "project_id": 29032549,
    "project": {
        "id": 29032549,
        "name": "sample-devcloud",
        "description": "测试使用",
        "web_url": "https://gitlab.com/yumaojun03/sample-devcloud",
        "avatar_url": null,
        "git_ssh_url": "git@gitlab.com:yumaojun03/sample-devcloud.git",
        "git_http_url": "https://gitlab.com/yumaojun03/sample-devcloud.git",
        "namespace": "紫川秀",
        "visibility_level": 0,
        "path_with_namespace": "yumaojun03/sample-devcloud",
        "default_branch": "main",
        "ci_config_path": "",
        "homepage": "https://gitlab.com/yumaojun03/sample-devcloud",
        "url": "git@gitlab.com:yumaojun03/sample-devcloud.git",
        "ssh_url": "git@gitlab.com:yumaojun03/sample-devcloud.git",
        "http_url": "https://gitlab.com/yumaojun03/sample-devcloud.git"
    },
    "commits": [
        {
            "id": "f8a831144634f5810e17014582b5ba21267bb257",
            "message": "Initial commit",
            "title": "Initial commit",
            "timestamp": "2021-08-22T03:44:35+00:00",
            "url": "https://gitlab.com/yumaojun03/sample-devcloud/-/commit/f8a831144634f5810e17014582b5ba21267bb257",
            "author": {
                "name": "紫川秀",
                "email": "9556442-yumaojun03@users.noreply.gitlab.com"
            },
            "added": [
                "README.md"
            ],
            "modified": [],
            "removed": []
        }
    ],
    "total_commits_count": 1,
    "push_options": {},
    "repository": {
        "name": "sample-devcloud",
        "url": "git@gitlab.com:yumaojun03/sample-devcloud.git",
        "description": "测试使用",
        "homepage": "https://gitlab.com/yumaojun03/sample-devcloud",
        "git_http_url": "https://gitlab.com/yumaojun03/sample-devcloud.git",
        "git_ssh_url": "git@gitlab.com:yumaojun03/sample-devcloud.git",
        "visibility_level": 0
    }
}
```

我们Hook处理就基于该样例数据进行进行

```go
func (h *handler) GitLabHookHanler(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get(GitlabEventHeaderKey)
	appID := r.Header.Get(GitlabEventTokenKey)
	switch eventType {
	case "Push Hook":
		event := scm.NewDefaultWebHookEvent()
		if err := request.GetDataFromRequest(r, event); err != nil {
			response.Failed(w, err)
			return
		}

		req := application.NewApplicationEvent(appID, event)
		h.log.Debugf("application %s accept event: %s", appID, event)

		_, err := h.service.HandleApplicationEvent(
			r.Context(),
			req,
		)
		if err != nil {
			response.Failed(w, err)
			return
		}
		response.Success(w, fmt.Sprintf("event %s has accept", event.ShortDesc()))
		return
	default:
		response.Failed(w, fmt.Errorf("known gitlab event type %s", eventType))
		return
	}
}
```

然后我们处理Push 事件
```go
func (s *service) HandleApplicationEvent(ctx context.Context, in *application.ApplicationEvent) (
	*application.Application, error) {
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate ApplicationEvent error, %s", err)
	}

	// 查询应用
	app, err := s.DescribeApplication(ctx, application.NewDescribeApplicationRequestWithID(in.AppId))
	if err != nil {
		return nil, err
	}

	// 找出匹配的pipline
	matched := app.MatchPipeline(in.WebhookEvent)
	if len(matched) == 0 {
		s.log.Infof("application %s no pipeline matched the event %s", app.Id, in.WebhookEvent.ShortDesc())
		return app, nil
	}

	// 运行这些匹配到的pipeline
	for i := range matched {
		req := matched[i]
		req.HookEvent = in.WebhookEvent
		req.Domain = app.Domain
		req.Namespace = app.Namespace
		req.CreateBy = fmt.Sprintf("@app:%s", app.Name)
		status := application.NewPipelineCreateStatus()

		s.log.Debugf("start create pipeline: %s", req.Name)
		p, err := s.pipeline.CreatePipeline(ctx, req)
		if err != nil {
			status.CreateError = err.Error()
		} else {
			status.PipelineId = p.Id
		}
	}

	// 更新应用状态
	if err := s.update(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}
```