package repositories

import (
	"mogong/internal/app/apigateway/models"

	"mogong/internal/pkg/common/database"
	"mogong/internal/pkg/common/database/mongo"
	"mogong/internal/pkg/common/database/mysql"

	"github.com/pkg/errors"
	mogongdb "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByID(id int64) (res *models.User, err error)
}

type DefaultUserRepository struct {
	logger  *zap.Logger
	dbWirte *gorm.DB
	dbRead  *gorm.DB
	mongo   *mogongdb.Database
}

func NewUserRepository(logger *zap.Logger, dbPool *mysql.DatabasePool, mongo *mongo.Database) UserRepo {
	return &DefaultUserRepository{
		logger:  logger.With(zap.String("type", "UserRepository")),
		dbWirte: dbPool.Write,
		dbRead:  database.GetReadDB(dbPool.Reads, dbPool.ReadCount),
		mongo:   mongo.MongoDb,
	}
}

func (d *DefaultUserRepository) GetUserByID(ID int64) (result *models.User, err error) {
	res := new(models.User)
	if err = d.dbRead.Model(res).Where("id = ?", ID).First(res).Error; err != nil {
		return nil, errors.Wrapf(err, "sql user error[id=%s]", ID)
	}
	return
}
