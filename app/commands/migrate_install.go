package commands

import (
	"fmt"
	"github.com/largezhou/gin_starter/app/model"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"io/fs"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

func NewMigrateInstallCommand() *cli.Command {
	return &cli.Command{
		Name:      "migrate:install",
		Usage:     "使用迁移文件",
		UsageText: "migrate:install",
		Category:  "migrate",
		Action: func(c *cli.Context) error {
			if err := createMigrationTableIfNeeded(); err != nil {
				return err
			}

			filenames, err := getMigrationFilenames()
			if err != nil {
				return err
			}

			migrationMap, maxBatch := getExistsMigrationInfo()
			var newFilenames []string
			for _, filename := range filenames {
				if _, ok := migrationMap[filename]; !ok {
					newFilenames = append(newFilenames, filename)
				}
			}

			if len(newFilenames) == 0 {
				fmt.Println("没有新的迁移")
				return nil
			}

			upRegexp, _ := regexp.Compile("(-- __UP__[\\s\\S]*?)-- __DOWN__")
			err = db.Transaction(func(tx *gorm.DB) error {
				var newMigrations []*Migration
				for _, filename := range newFilenames {
					newMigrations = append(newMigrations, &Migration{
						File:  filename,
						Batch: maxBatch + 1,
					})

					file, err := os.Open("./migration/" + filename)
					if err != nil {
						return nil
					}

					fileContent, err := ioutil.ReadAll(file)
					if err != nil {
						return err
					}

					matches := upRegexp.FindSubmatch(fileContent)
					if len(matches) != 2 {
						return fmt.Errorf("%s 未找到升级 SQL", filename)
					}

					fmt.Printf("正在迁移 %s\n", filename)

					sqlList := strings.Split(string(matches[1]), ";")
					for _, sql := range sqlList {
						if strings.TrimSpace(sql) == "" {
							continue
						}
						if result := tx.Exec(sql); result.Error != nil {
							return result.Error
						}
					}
				}

				if result := tx.Create(newMigrations); result.Error != nil {
					return result.Error
				}

				return nil
			})

			return err
		},
	}
}

type Migration struct {
	model.Model
	File  string `gorm:"type:varchar(200);unique;not null"`
	Batch uint   `gorm:"type:integer;not null"`
}

func createMigrationTableIfNeeded() error {
	if db.Migrator().HasTable(&Migration{}) {
		return nil
	}

	if err := db.Migrator().CreateTable(&Migration{}); err != nil {
		return err
	}

	return nil
}

func getMigrationFilenames() ([]string, error) {
	var files []fs.FileInfo
	var filenames []string
	files, err := ioutil.ReadDir("./migration/")
	if err != nil {
		return filenames, err
	}

	for _, fileInfo := range files {
		name := strings.ToLower(fileInfo.Name())
		if strings.HasSuffix(name, ".sql") {
			filenames = append(filenames, name)
		}
	}

	sort.Strings(filenames)

	return filenames, nil
}

func getExistsMigrationInfo() (map[string]Migration, uint) {
	msMap := make(map[string]Migration)
	var maxBatch uint = 0

	var ms []Migration
	db.Model(&Migration{}).Select("file, batch").Find(&ms)
	for _, m := range ms {
		if m.Batch > maxBatch {
			maxBatch = m.Batch
		}
		msMap[m.File] = m
	}
	return msMap, maxBatch
}
