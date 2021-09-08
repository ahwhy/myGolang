package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/ahwhy/myGolang/week14/demo/api/conf"
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
	s.db = db

	s.log = zap.L().Named("Host")

	return nil
}
