# CMDB 同步API

联合secret和provider 我们就可以出一个 使用secret同步资产的接口了:
```go
r.POST("/secrets/:id/sync", api.Sync)
```


我们接着扩展secret模块

## 接口定义

需要同步某个key下的 某个区域的资源，因此请求有3个参数:
+ SecretId
+ Region
+ ResourceType: 现在只有Host


接口需要的数据结构: sync.go
```go
func NewSyncRequest(secretId string) *SyncRequest {
	return &SyncRequest{
		SecretId: secretId,
	}
}


type SyncRequest struct {
	Region       string
	SecretId     string `validate:"required,lte=100"`
	ResourceType resource.Type
}

func NewSyncReponse() *SyncReponse {
	return &SyncReponse{
		Details: []*SyncDetail{},
	}
}

// 同步后返回 那些同步成功，那些同步失败, 方便界面暂时
type SyncReponse struct {
	TotolSucceed int64         `json:"total_succeed"`
	TotalFailed  int64         `json:"total_failed"`
	Details      []*SyncDetail `json:"details"`
}
```

定义接口: interface.go:
```go
type Service interface {
	SecretService
	SyncService
}

type SyncService interface {
	Sync(context.Context, *SyncRequest) (*SyncReponse, error)
}
```

## 实现同步

利用secret 和 provider 就可以实现同步的逻辑


### Secret与Provider结合

+ 通过SecretID 找到对应Secret
+ 利用Secret保存的数据, 集合Provider同步资源

Sync逻辑: impl/sync.go
```go
func (s *service) Sync(ctx context.Context, req *syncer.SyncRequest) (
	*syncer.SyncReponse, error) {
	var (
		resp *syncer.SyncReponse
		err  error
	)

	// 校验参数合法性
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate sync request error, %s", err)
	}

	// 通过SecretID找到对应的Secret对象
	secret, err := s.DescribeSecret(ctx, syncer.NewDescribeSecretRequest(req.SecretId))
	if err != nil {
		return nil, err
	}

	// 如果不是vsphere 需要检查region
	if !secret.Vendor.Equal(resource.VendorVsphere) {
		if req.Region == "" {
			return nil, exception.NewBadRequest("region required")
		}
		if !secret.IsAllowRegion(req.Region) {
			return nil, exception.NewBadRequest("this secret not allow sync region %s", req.Region)
		}
	}

	// 根据请求的资源进行具体的资源同步
	switch req.ResourceType {
	case resource.HostResource:
		resp, err = s.syncHost(ctx, secret, req.Region)
	case resource.RdsResource:
		resp, err = s.syncRds(ctx, secret, req.Region)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
```

### 同步主机

+ 解密secret key
+ new ecs client
+ page query
+ 读取所有Page数据
+ 调用host服务保存
+ 所有host同步完成后，返回结果

```go
package impl

import (
	"context"
	"fmt"

	"github.com/infraboard/cmdb/conf"
	"github.com/infraboard/cmdb/pkg/host"
	"github.com/infraboard/cmdb/pkg/resource"
	"github.com/infraboard/cmdb/pkg/syncer"
	"github.com/infraboard/mcube/exception"

	aliConn "github.com/infraboard/cmdb/provider/aliyun/connectivity"
	ecsOp "github.com/infraboard/cmdb/provider/aliyun/ecs"
	hwConn "github.com/infraboard/cmdb/provider/huawei/connectivity"
	hwEcsOp "github.com/infraboard/cmdb/provider/huawei/ecs"
	txConn "github.com/infraboard/cmdb/provider/txyun/connectivity"
	cvmOp "github.com/infraboard/cmdb/provider/txyun/cvm"
	vsConn "github.com/infraboard/cmdb/provider/vsphere/connectivity"
	vmOp "github.com/infraboard/cmdb/provider/vsphere/vm"
)

func (s *service) syncHost(ctx context.Context, secret *syncer.Secret, region string) (
	*syncer.SyncReponse, error) {
	var (
		pager host.Pager
	)

	// 解密secret
	err := secret.DecryptAPISecret(conf.C().App.EncryptKey)
	if err != nil {
		s.log.Warnf("decrypt api secret error, %s", err)
	}

	hs := host.NewHostSet()
	switch secret.Vendor {
	case resource.VendorAliYun:
		s.log.Debugf("sync aliyun host ...")
		client := aliConn.NewAliCloudClient(secret.APIKey, secret.APISecret, region)
		ec, err := client.EcsClient()
		if err != nil {
			return nil, err
		}
		operater := ecsOp.NewEcsOperater(ec)
		req := ecsOp.NewPageQueryRequest()
		req.Rate = secret.RequestRate
		pager = operater.PageQuery(req)
	case resource.VendorTencent:
		s.log.Debugf("sync txyun host ...")
		client := txConn.NewTencentCloudClient(secret.APIKey, secret.APISecret, region)
		operater := cvmOp.NewCVMOperater(client.CvmClient())
		pager = operater.PageQuery()
	case resource.VendorHuaWei:
		s.log.Debugf("sync hwyun host ...")
		client := hwConn.NewHuaweiCloudClient(secret.APIKey, secret.APISecret, region)
		ec, err := client.EcsClient()
		if err != nil {
			return nil, err
		}
		operater := hwEcsOp.NewEcsOperater(ec)
		pager = operater.PageQuery()
	case resource.VendorVsphere:
		s.log.Debugf("sync vshpere host ...")
		client := vsConn.NewVsphereClient(secret.Address, secret.APIKey, secret.APISecret)
		ec, err := client.VimClient()
		if err != nil {
			return nil, err
		}
		operater := vmOp.NewVmOperater(ec)
		hs, err = operater.Query()
		if err != nil {
			return nil, err
		}
	default:
		return nil, exception.NewBadRequest("unsuport vendor %s", secret.Vendor)
	}

	set := syncer.NewSyncReponse()
	// 分页查询数据
	if pager != nil {
		hasNext := true
		for hasNext {
			p := pager.Next()
			hasNext = p.HasNext

			if p.Err != nil {
				return nil, fmt.Errorf("sync error, %s", p.Err)
			}

			// 调用host服务保持数据
			for i := range p.Data.Items {
				hs.Add(p.Data.Items[i])
			}
		}
	}

	// 调用host服务保持数据
	for i := range hs.Items {
		target := hs.Items[i]
		_, err := s.host.SaveHost(ctx, target)
		if err != nil {
			set.AddFailed(target.Name, err.Error())
		} else {
			set.AddSucceed(target.Name, "")
		}
	}

	return set, nil
}
```

## 参考

+ [CMDB项目源码](https://github.com/infraboard/cmdb)