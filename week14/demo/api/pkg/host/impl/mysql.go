package impl

import (
	"database/sql"

	"gitee.com/infraboard/go-course/day14/demo/api/conf"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	// Service 服务实例
	Service = &service{}
)

type service struct {
	db  *sql.DB
	log logger.Logger
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named("Host")
	s.db = db
	return nil
}
