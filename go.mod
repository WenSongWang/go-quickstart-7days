// 本项目的「模块路径」：Go 用这个名字来标识当前项目。
// 下面两件事容易误导新手，务必弄清：
//
// 1) 为什么 import 本项目的包要写 github.com/... 而不是 ./day3/...？
//    Go 规定：包路径 = 模块路径 + 相对目录。所以本项目里的 day3/internal/config
//    的 import 路径就是 github.com/go-quickstart-7days/day3/internal/config。
//    开发时 Go 从【本地目录】解析，不会从 GitHub 拉代码。
//
// 2) require 里的 github.com/joho/godotenv 等才是第三方依赖，会由 go mod tidy
//    下载到本机缓存；编译、运行用的都是本地缓存，也不是「运行时从 GitHub 拉」。
//
module github.com/go-quickstart-7days

go 1.21

require (
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/spf13/viper v1.18.2
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
