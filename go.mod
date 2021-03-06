module github.com/kassybas/tame

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/antonmedv/expr v1.4.2
	github.com/kassybas/shell-exec v0.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.20.0
	golang.org/x/sys v0.0.0-20190626221950-04f50cda93cb // indirect
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/kassybas/shell-exec v0.1.0 => github.com/kassybas/shell-exec v0.1.2
