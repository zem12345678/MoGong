package repositories

import (
	"mogong/internal/app/apigateway/models"
	"mogong/internal/pkg/common/database"
	"mogong/internal/pkg/common/database/mysql"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppRepo interface {
	Create(t *models.App) (*models.App, error)
	Save(t *models.App) error
	List(search string, page, pageSize int) ([]models.App, int64, error)
	Delete(ID int64) error
}

type DefaultAppRepository struct {
	logger  *zap.Logger
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func NewAppRepository(logger *zap.Logger, dbPool *mysql.DatabasePool) AppRepo {
	return &DefaultAppRepository{
		logger:  logger.With(zap.String("type", "AppRepository")),
		dbWrite: dbPool.Write,
		dbRead:  database.GetReadDB(dbPool.Reads, dbPool.ReadCount),
	}
}

func (d *DefaultAppRepository) Create(t *models.App) (*models.App, error) {
	if err := d.dbWrite.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (d *DefaultAppRepository) Save(t *models.App) error {
	if err := d.dbWrite.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (d *DefaultAppRepository) List(search string, page, pageSize int) ([]models.App, int64, error) {
	var list []models.App
	var count int64
	query := d.dbRead.Where("is_delete = ?", 0)
	if search != "" {
		query = query.Where(" (name like ? or app_id like ?)", "%"+search+"%", "%"+search+"%")
	}
	err := query.Limit(pageSize).Offset(page).Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (d *DefaultAppRepository) Delete(ID int64) error {
	if err := d.dbWrite.Where("id = ?", ID).Delete(&models.App{}).Error; err != nil {
		return err
	}
	return nil
}

var AppManagerHandler *AppManager

func init() {
	AppManagerHandler = NewAppManager()
}

type AppManager struct {
	AppMap   map[string]*models.App
	AppSlice []*models.App
	Locker   sync.RWMutex
	init     sync.Once
	err      error
}

func NewAppManager() *AppManager {
	return &AppManager{
		AppMap:   map[string]*models.App{},
		AppSlice: []*models.App{},
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

func (s *AppManager) GetAppList() []*models.App {
	return s.AppSlice
}

func (s *AppManager) LoadOnce(d *DefaultAppRepository) error {
	s.init.Do(func() {
		list, _, err := d.List("", 1, 99999)
		if err != nil {
			s.err = err
			return
		}
		s.Locker.Lock()
		defer s.Locker.Unlock()
		for _, listItem := range list {
			tmpItem := listItem
			s.AppMap[listItem.AppID] = &tmpItem
			s.AppSlice = append(s.AppSlice, &tmpItem)
		}
	})
	return s.err
}
