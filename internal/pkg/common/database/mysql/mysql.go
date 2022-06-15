package mysql

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	driverMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Options struct {
	Mysql struct {
		Dsn   string
		Reads []string
		Debug bool
	}
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("mysql", &o.Mysql); err != nil {
		return nil, errors.Wrap(err, "unmarshal mysql option error")
	}

	logger.Info("load mysql options success", zap.Any("mysql options", o))
	return o, err
}

// Database 定义数据库struct
type Database struct {
	Mysql *gorm.DB
}

// DBClient  mysql连接类型
var DBClient Database

type DatabasePool struct {
	Write     *gorm.DB
	Reads     []*gorm.DB
	ReadCount int
}

// DBClient  mysql连接类型
var DBClientPool DatabasePool

func New(o *Options) (*Database, error) {
	var d = new(Database)
	if o.Mysql.Dsn == "" {
		return nil, errors.New("缺少mysql写库配置")
	}
	if o.Mysql.Dsn != "" {
		mysql, err := mysql(o)
		if err != nil {
			return nil, err
		}
		d.Mysql = mysql
	}
	DBClient.Mysql = d.Mysql

	return d, nil
}

func NewPools(o *Options) (*DatabasePool, error) {
	var dpool = new(DatabasePool)
	if o.Mysql.Dsn == "" {
		return nil, errors.New("缺少mysql写库配置")
	}
	if o.Mysql.Dsn != "" {
		mysql, err := mysql(o)
		if err != nil {
			return nil, err
		}
		dpool.Write = mysql
	}
	DBClientPool.Write = dpool.Write
	if len(o.Mysql.Reads) > 0 {
		for _, url := range o.Mysql.Reads {
			dbRead, err := gorm.Open(driverMysql.Open(url), &gorm.Config{})
			if err != nil && dbRead == nil {
				return dpool, errors.Wrap(err, "gorm open mysql only read connection error")
			}
			sqlDb, err := dbRead.DB()
			if err != nil {

				return dpool, errors.Wrap(err, "mysql db fail")
			}
			err = sqlDb.Ping()
			if err != nil {
				return dpool, errors.Wrap(err, "mysql ping fail")
			}
			if o.Mysql.Debug {
				dbRead = dbRead.Debug()
			}
			sqlDb.SetConnMaxLifetime(time.Minute * 10)
			sqlDb.SetMaxIdleConns(10)
			sqlDb.SetMaxOpenConns(100)
			dpool.Reads = append(dpool.Reads, dbRead)
		}
		dpool.ReadCount = len(dpool.Reads)
		DBClientPool.Reads = dpool.Reads
		DBClientPool.ReadCount = dpool.ReadCount
	}

	return dpool, nil
}

// New new database
func NewPool(o *Options) (*Database, *DatabasePool, error) {
	var d = new(Database)
	var dpool = new(DatabasePool)
	if o.Mysql.Dsn == "" {
		return nil, nil, errors.New("缺少mysql写库配置")
	}
	if o.Mysql.Dsn != "" {
		mysql, err := mysql(o)
		if err != nil {
			return nil, nil, err
		}
		d.Mysql = mysql
	}
	DBClient.Mysql = d.Mysql
	DBClientPool.Write = d.Mysql
	if len(o.Mysql.Reads) > 0 {
		for _, url := range o.Mysql.Reads {
			dbRead, err := gorm.Open(driverMysql.Open(url), &gorm.Config{})
			if err != nil && dbRead == nil {
				return d, nil, errors.Wrap(err, "gorm open mysql only read connection error")
			}
			sqlDb, err := dbRead.DB()
			if err != nil {
				return d, nil, errors.Wrap(err, "mysql db fail")
			}
			err = sqlDb.Ping()
			if err != nil {
				return d, nil, errors.Wrap(err, "mysql ping fail")
			}
			if o.Mysql.Debug {
				dbRead = dbRead.Debug()
			}
			sqlDb.SetConnMaxLifetime(time.Minute * 10)
			sqlDb.SetMaxIdleConns(10)
			sqlDb.SetMaxOpenConns(100)
			dpool.Reads = append(dpool.Reads, dbRead)
		}
		dpool.ReadCount = len(dpool.Reads)
		DBClientPool.Reads = dpool.Reads
		DBClientPool.ReadCount = dpool.ReadCount
	}

	return d, dpool, nil
}

// mysql 定义mysql连接信息
func mysql(o *Options) (*gorm.DB, error) {
	db, err := gorm.Open(driverMysql.Open(o.Mysql.Dsn), &gorm.Config{})
	if err != nil && db == nil {
		return nil, errors.Wrap(err, "gorm open mysql connection error")
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "mysql db fail")
	}
	err = sqlDb.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "mysql ping fail")
	}
	if o.Mysql.Debug {
		db = db.Debug()
	}
	sqlDb.SetConnMaxLifetime(time.Minute * 10)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	// 自动迁移模式
	err = db.AutoMigrate()
	if err != nil {
		return db, err
	}

	return db, nil
}

// ProviderSet dependency injection
