package config

import (
	"cendit.io/garage/database"
	"cendit.io/garage/database/migration"
	"cendit.io/garage/logger"

	"cendit.io/garage/primer/typing"
	"cendit.io/garage/redis"
	"cendit.io/signal/email"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Dependencies struct {
	// DAL - Database Access Layer
	DB *bun.DB

	// RAL - Redis Access Layer
	Redis *redis.RAL

	Email typing.EmailClient

	SMS typing.SMSClient
}

var Deps *Dependencies

// Deps initiliazes the project dependencies based on the config
func ResolveDeps(cfg *Variable) (*Dependencies, error) {

	if Deps != nil {
		return Deps, nil
	}

	db, err := database.New(cfg.CenditDatabaseURL, cfg.DebugDatabase, cfg.DatabaseConnectionLimit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up data access layer")
	}
	logger.GetLogger().Info("[DB]: OK")

	createErr := migration.CreateTables(db, migration.Tables)
	if createErr != nil {
		logger.GetLogger().Sugar().Fatalf("[DB]: unable to create new tables: %v", createErr.Error())
		return nil, createErr
	}
	logger.GetLogger().Info("[DB]: tables updated")

	migrateErr := migration.Migrate(db)
	if migrateErr != nil {
		logger.GetLogger().Sugar().Fatalf("[DB]: unable to migrate schema: %v", migrateErr.Error())
		return nil, migrateErr
	}
	logger.GetLogger().Info("[DB]: migration competed")

	redisS, err := redis.New(cfg.RedisURL)
	if err != nil {
		logger.GetLogger().Sugar().Fatalf("[REDIS]: unable to connect to redis: %v", err.Error())
		return nil, err
	}
	logger.GetLogger().Info("[REDIS]: OK")

	emailClient := email.NewClient(cfg.Email.ClientID, cfg.Email.Secret, cfg.Email.SendPulseBaseURL)
	if err := emailClient.Authenticate(); err != nil {
		logger.GetLogger().Sugar().Fatalf("[EMAIL]: unable to authenticate email client: %v", err.Error())
		return nil, err
	}
	logger.GetLogger().Info("[EMAIL]: OK")

	// smsClient := sms.NewClient(cfg.SMS.SendChampKey, cfg.SMS.SendChampBaseURL, cfg.SMS.SendChampSender, cfg.SMS.SendChampKey)
	// if err = smsClient.Authenticate(); err != nil {
	// 	logger.GetLogger().Sugar().Fatalf("[SMS]: unable to authenticate sms client: %v", err.Error())
	// 	return nil, err
	// }
	// logger.GetLogger().Info("[SMS]: OK")

	Deps = &Dependencies{
		DB:    db,
		Redis: redisS,
		Email: emailClient,
		// SMS:   smsClient,
	}

	return Deps, nil
}
