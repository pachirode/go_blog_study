package main

import (
	"github.com/pachirode/go_blog_study/internal/pkg/log"
	"github.com/pachirode/pkg/db"
	"gorm.io/gorm"
	"path/filepath"
)

func initializeDatabase() (*gorm.DB, error) {
	dbOptions := &db.MySQLOptions{
		Addr:     *addr,
		Username: *username,
		Password: *password,
		Database: *database,
	}

	return db.NewMySQL(dbOptions)
}

func resolveModelPackagePath(defaultPath string) string {
	if *modelPath != "" {
		return *modelPath
	}

	absPath, err := filepath.Abs(defaultPath)
	if err != nil {
		log.Errorw("Error resolving path", "err", err)
		return defaultPath
	}

	return absPath
}

func processComponent(component string, dbInstance *gorm.DB) {
	config, ok := generateConfigs[component]
	if !ok {
		log.Errorw("Component not found, skip", "component", component)
		return
	}

	// 解析模型包路径
	modelPackagePath := resolveModelPackagePath(config.ModelPackagePath)

	generator := createGenerator(modelPackagePath)
	generator.UseDB(dbInstance)
	applyGeneratorOptions(generator)

	config.GenerateFunc(generator)
	generator.Execute()
}
