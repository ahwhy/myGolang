# 云资产凭证管理

我们对4个厂商的凭证进行抽象:
```
TX_CLOUD_SECRET_ID=xxx
TX_CLOUD_SECRET_KEY=xxx
AL_CLOUD_ACCESS_KEY=xx
AL_CLOUD_ACCESS_SECRET=xxx
HW_CLOUD_ACCESS_KEY=xxx
HW_CLOUD_ACCESS_SECRET=xxx
VS_HOST=xxx
VS_USERNAME=xxxx
VS_PASSWORD=xxx
```

我们定义2中凭证类型
```go
const (
	CrendentialAPIKey CrendentialType = iota
	CrendentialPassword
)

type CrendentialType int
```

针对vsphere, 需要独立出一个字段: VS_HOST 用来保存 vshpere服务的地址


我们定义secret表，来存储我们通过过程中要使用到的这些密钥:
```sql
CREATE TABLE `secret` (
  `id` varchar(64) NOT NULL,
  `create_at` bigint(13) NOT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 NOT NULL,
  `vendor` tinyint(1) NOT NULL,
  `address` varchar(255) NOT NULL,
  `allow_regions` text NOT NULL,
  `crendential_type` tinyint(1) NOT NULL,
  `api_key` varchar(255) NOT NULL,
  `api_secret` text NOT NULL,
  `request_rate` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`api_key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=latin1
```

## Secret 相关接口定义

定义数据结构
```go
type Secret struct {
	Id       string `json:"id"`        // 全局唯一Id
	CreateAt int64  `json:"create_at"` // 创建时间

	*CreateSecretRequest
}

type CreateSecretRequest struct {
	Description     string          `json:"description" validate:"required,lte=100"` // 描述
	Vendor          resource.Vendor `json:"vendor"`                                  // 厂商
	AllowRegions    []string        `json:"allow_regions"`                           // 允许同步的区域
	CrendentialType CrendentialType `json:"crendential_type"`                        // 凭证类型
	Address         string          `json:"address"`                                 // 服务地址, 云商不用填写
	APIKey          string          `json:"api_key" validate:"required,lte=100"`     // key
	APISecret       string          `json:"api_secret" validate:"required,lte=100"`  // secrete
	RequestRate     int             `json:"request_rate"`                            // 请求速率限制, 默认1秒5个
}
```

定义接口
```go
type Service interface {
	SecretService
}

type SecretService interface {
	CreateSecret(context.Context, *CreateSecretRequest) (*Secret, error)
	QuerySecret(context.Context, *QuerySecretRequest) (*SecretSet, error)
	DescribeSecret(context.Context, *DescribeSecretRequest) (*Secret, error)
}
```

全局实例添加该Service
```go
package pkg

import (
	"github.com/infraboard/cmdb/pkg/host"
	"github.com/infraboard/cmdb/pkg/syncer"
)

var (
	Host   host.Service
	Syncer syncer.Service
)
```

## Secret CRUD实现

impl模块里面定义个secret.go,  用于实现SecretService, 具体参考: [secret CRUD](https://github.com/infraboard/cmdb/blob/master/pkg/syncer/impl/secret.go)

然后编写HTTP暴露模块:
```go
func (h *handler) QuerySecret(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := syncer.NewQuerySecretRequest()
	set, err := h.service.QuerySecret(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) CreateSecret(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := syncer.NewCreateSecretRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.CreateSecret(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *handler) DescribeSecret(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := syncer.NewDescribeSecretRequest(ps.ByName("id"))
	set, err := h.service.DescribeSecret(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}
```

最后由handler注册:
```go
var (
	api = &handler{}
)

type handler struct {
	service syncer.Service
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named("Syncer")
	if pkg.Syncer == nil {
		return fmt.Errorf("dependence service syncer not ready")
	}
	h.service = pkg.Syncer
	return nil
}

func RegistAPI(r *httprouter.Router) {
	api.Config()
	r.POST("/secrets", api.CreateSecret)
	r.GET("/secrets", api.QuerySecret)
	r.GET("/secrets/:id", api.DescribeSecret)
}
```

## 注册服务与加载路由

1. secret 服务注册: cmd/start.go
```go
// 初始化服务层 Ioc初始化
if err := impl.Service.Config(); err != nil {
	return err
}
pkg.Host = impl.Service

// Syncer Service
if err := syncer.Service.Config(); err != nil {
	return err
}
pkg.Syncer = syncer.Service
```

2. 注册http服务  protocol/http.go

```go
// Start 启动服务
func (s *HTTPService) Start() error {
	// 装置子服务路由
	hostAPI.RegistAPI(s.r)
	syncerAPI.RegistAPI(s.r)
	...
}
```

3.  Post测试

通过Postman测试接口的使用性

## Secret Key加密与脱敏

现在我们的Secret Key都是明文存储的, 为了安全, 我们加密存储, 

为了确保安全, Secret在返回的时候，需要做脱敏, 及时是秘文也不应该显示

### Secret加密存储

我们采用cbc来实现对称加密: 这是之前封装的模块 [cbc加密](github.com/infraboard/mcube/crypto/cbc)

为了能区分 原始数据和密码数据，我们给秘文数据加一个前缀，以避免重复加密或者解密

conf/config.go: 
```go
const (
	CIPHER_TEXT_PREFIX = "@ciphered@"
)
```

然后为我们secret对象 添加加解密的方法:
```go
func (s *Secret) EncryptAPISecret(key string) error {
	// 判断文本是否已经加密
	if strings.HasPrefix(s.APISecret, conf.CIPHER_TEXT_PREFIX) {
		return fmt.Errorf("text has ciphered")
	}

	cipherText, err := cbc.Encrypt([]byte(s.APISecret), []byte(key))
	if err != nil {
		return err
	}

	base64Str := base64.StdEncoding.EncodeToString(cipherText)
	s.APISecret = fmt.Sprintf("%s%s", conf.CIPHER_TEXT_PREFIX, base64Str)
	return nil
}

func (s *Secret) DecryptAPISecret(key string) error {
	// 判断文本是否已经是明文
	if !strings.HasPrefix(s.APISecret, conf.CIPHER_TEXT_PREFIX) {
		return fmt.Errorf("text is plan text")
	}

	base64CipherText := strings.TrimPrefix(s.APISecret, conf.CIPHER_TEXT_PREFIX)

	cipherText, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return err
	}

	planText, err := cbc.Decrypt([]byte(cipherText), []byte(key))
	if err != nil {
		return err
	}

	s.APISecret = string(planText)
	return nil
}
```


### Secret脱敏显示

为secret对象补充一个脱敏方法:
```go
func (s *Secret) Desense() {
	if s.APISecret != "" {
		s.APISecret = "******"
	}
}
```

对List方法查询出来的数据进行过敏:
```go
	for rows.Next() {
		ins := syncer.NewDefaultSecret()
		...
		ins.Desense()
		set.Add(ins)
	}
```

针对Get方法, 在API暴露时脱敏:
```go
func (h *handler) DescribeSecret(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := syncer.NewDescribeSecretRequest(ps.ByName("id"))
	ins, err := h.service.DescribeSecret(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	ins.Desense()
	response.Success(w, ins)
}
```

这样我们系统内部交互时, 通过Describe拿到的数据 依然是可以解密的


## 参考

+ [CMDB项目源码](https://github.com/infraboard/cmdb)
