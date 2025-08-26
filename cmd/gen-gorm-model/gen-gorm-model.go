package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"gorm.io/gen"
	"gorm.io/gen/field"

	"github.com/pachirode/go_blog_study/internal/pkg/log"
)

const helpText = `Usage: main [flags] arg [arg...]

This is a pflag example.

Flags:
`

// IQuery 定义了数据库查询接口.
type IQuery interface {
	// FilterWithNameAndRole 按名称和角色查询记录
	FilterWithNameAndRole(name string) ([]gen.T, error)
}

// GenerateConfig 保存代码生成配置
type GenerateConfig struct {
	ModelPackagePath string
	GenerateFunc     func(query *gen.Generator)
}

// 预定义的生成配置
var generateConfigs = map[string]GenerateConfig{
	"mb": {
		ModelPackagePath: "internal/apiserver/model",
		GenerateFunc:     GenerateMiniBlogModels,
	},
}

// 命令行参数
var (
	addr       = pflag.StringP("addr", "a", "127.0.0.1:3306", "MySQL host address.")
	username   = pflag.StringP("username", "u", "miniblog", "Username to connect to the database.")
	password   = pflag.StringP("password", "p", "miniblog1234", "Password to use when connecting to the database.")
	database   = pflag.StringP("db", "d", "miniblog", "Database name to connect to.")
	modelPath  = pflag.String("model-pkg-path", "", "Generated model code's package name.")
	components = pflag.StringSlice("component", []string{"mb"}, "Generated model code's for specified component.")
	help       = pflag.BoolP("help", "h", false, "Show this help message.")
)

func main() {

	pflag.Usage = func() {
		fmt.Printf("%s", helpText)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}

	dbInstance, err := initializeDatabase()
	if err != nil {
		log.Fatalw("Failed to connect to database: %v", err)
	}

	for _, component := range *components {
		processComponent(component, dbInstance)
	}
}

// createGenerator 初始化一个新的生成器实例
func createGenerator(packagePath string) *gen.Generator {
	return gen.NewGenerator(gen.Config{
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		ModelPkgPath:      packagePath,
		WithUnitTest:      true,
		FieldNullable:     true,  // 数据库中可空的字段，使用指针类型
		FieldSignable:     false, // 禁用无符号属性提高兼容性
		FieldWithIndexTag: false, // 不包含 GORM 的索引标签
		FieldWithTypeTag:  false, // 不包含 GORM 的类型标签
	})
}

// applyGeneratorOptions  设置自定义生成器选项
func applyGeneratorOptions(g *gen.Generator) {
	g.WithOpts(
		gen.FieldGORMTag("createdAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),

		gen.FieldGORMTag("updatedAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
	)
}

// GenerateMiniBlogModels 为 miniblog 组件生成模型.
func GenerateMiniBlogModels(g *gen.Generator) {
	g.GenerateModelAs(
		"user",
		"UserM",
		gen.FieldIgnore("placeholder"),
		gen.FieldGORMTag("username", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_username")
			return tag
		}),
		gen.FieldGORMTag("userID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_userID")
			return tag
		}),
		gen.FieldGORMTag("phone", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_phone")
			return tag
		}),
	)
	g.GenerateModelAs(
		"post",
		"PostM",
		gen.FieldIgnore("placeholder"),
		gen.FieldGORMTag("postID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_post_postID")
			return tag
		}),
	)
	g.GenerateModelAs(
		"casbin_rule",
		"CasbinRuleM",
		gen.FieldRename("ptype", "PType"),
		gen.FieldIgnore("placeholder"),
	)
}
