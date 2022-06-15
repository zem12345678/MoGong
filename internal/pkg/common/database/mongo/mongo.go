package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type Options struct {
	Mongo struct {
		Name  string `yaml:"name"`
		URL   string `yaml:"url"`
		Debug bool
	}
}

type Database struct {
	MongoDb *mongo.Database
	Client  *mongo.Client
	Context context.Context
}

var DBClient Database

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("mongo", &o.Mongo); err != nil {
		return nil, errors.Wrap(err, "unmarshal database option error")
	}
	logger.Info("load mongo options success", zap.Any("mongo options", o))
	return o, err
}

func New(o *Options) (*Database, error) {
	var d = new(Database)
	if o.Mongo.URL == "" || o.Mongo.Name == "" {
		return nil, errors.New("缺少mongobd配置")
	} else {
		mongodb, err := mongoDB(o)
		if err != nil {
			return nil, err
		}
		d.MongoDb = mongodb
	}
	DBClient.MongoDb = d.MongoDb
	return d, nil
}

func mongoDB(o *Options) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	op := options.Client().ApplyURI(o.Mongo.URL)
	op.SetMaxPoolSize(10)
	client, err := mongo.Connect(ctx, op)
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI(o.mongodb.URL))
	if err != nil {
		return nil, errors.Wrap(err, "mongo driver open mongodb connection error")
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Wrap(err, "mongo ping fail")
	}
	db := client.Database(o.Mongo.Name)
	//res, err := db.Collection("test").InsertOne(ctx,map[string]string{"key":"value"})
	//fmt.Println(res,err)
	return db, nil
}

func (db *Database) Close() {
	db.Client.Disconnect(db.Context)
}
