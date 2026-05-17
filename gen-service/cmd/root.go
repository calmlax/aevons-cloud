package cmd

import (
	"fmt"
	"gen-service/internal/repository"
	"gen-service/internal/service"
	"os"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var cfg config.GenConfig
var database *gorm.DB

var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "Aevons 代码生成器",
	Run:   runGenerate, // 执行生成逻辑
	Long: `
==================================================
Aevons 代码生成工具
用于快速生成 Model,DTO,Repository,Service,Handler,Router,Sql,Vue,Api 等标准模块代码
支持自定义模板路径、输出目录，一键生成完整业务模块

使用:
  gen <表名> <归属模块名称>
  gen <表名1,表名2...> <归属模块名称>
示例:
  gen sys_test sys       # 生成 sys_test 表，归属 sys 模块
==================================================
`,
	SilenceUsage:  false,
	SilenceErrors: false,
}

// 核心：代码生成执行函数
func runGenerate(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("==================================================")
		fmt.Println("❌ 请输入要生成的表名！")
		fmt.Println()
		fmt.Println("使用:")
		fmt.Println("  gen <表名> <归属模块名称>")
		fmt.Println("  gen <表名1,表名2...> <归属模块名称>")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  gen sys_test sys")
		fmt.Println("==================================================")
		return
	}
	// 1. 解析参数
	tableNameStr := args[0]
	moduleName := args[1]

	// 2. 输出信息
	fmt.Println("========================================")
	fmt.Println("🚀 开始生成代码")
	fmt.Println("📦 表名：", tableNameStr)
	fmt.Println("📂 模块：", moduleName)
	fmt.Println("========================================")

	// ==============================
	// 3. 在这里写你的真正生成逻辑！

	srv := service.NewGenTableService(repository.NewGenTableRepository(database), repository.NewGenTableColumnRepository(database), cfg)
	tableNames := utils.Split(tableNameStr, ",")
	_, err := srv.CommandGenerateCode(tableNames, moduleName)
	if err != nil {
		fmt.Println("❌ 生成失败：", err.Error())
	} else {
		fmt.Println("✅ 代码生成完成！")
	}
}

func Execute(genConf config.GenConfig, db *gorm.DB) {
	cfg = genConf
	database = db
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultHelpFlag()
	if len(os.Args) >= 2 {
		rootCmd.SetArgs(os.Args[2:])
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
