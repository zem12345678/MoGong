package repositories

import (
	"mogong/internal/app/admin/models"

	"mogong/internal/pkg/common/database"
	"mogong/internal/pkg/common/database/mongo"
	"mogong/internal/pkg/common/database/mysql"

	"github.com/pkg/errors"
	mogongdb "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysUserRepo interface {
	GetUserByID(id int64) (res *models.SysUser, err error)
	GetUserByUserName(userName string) (res *models.SysUser, err error)
}

type DefaultSysUserRepository struct {
	logger  *zap.Logger
	dbWirte *gorm.DB
	dbRead  *gorm.DB
	mongo   *mogongdb.Database
}

func NewSysUserRepository(logger *zap.Logger, dbPool *mysql.DatabasePool, mongo *mongo.Database) SysUserRepo {
	return &DefaultSysUserRepository{
		logger:  logger.With(zap.String("type", "SysUserRepository")),
		dbWirte: dbPool.Write,
		dbRead:  database.GetReadDB(dbPool.Reads, dbPool.ReadCount),
		mongo:   mongo.MongoDb,
	}
}

func (d *DefaultSysUserRepository) GetUserByID(ID int64) (res *models.SysUser, err error) {
	res = new(models.SysUser)
	if err = d.dbRead.Model(res).Where("id = ?", ID).First(res).Error; err != nil {
		return nil, errors.Wrapf(err, "sql user error[id=%s]", ID)
	}
	return
}

func (d *DefaultSysUserRepository) GetUserByUserName(userName string) (res *models.SysUser, err error) {
	res = new(models.SysUser)
	if err = d.dbRead.Model(res).Where("username = ?", userName).First(res).Error; err != nil {
		return nil, errors.Wrapf(err, "get user error[user_name=%s]", userName)
	}
	return
}
