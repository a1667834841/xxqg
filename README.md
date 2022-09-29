## 学习强国自动化学习


## 鸣谢

+ ### [johlanse/study_xxqg](https://github.com/johlanse/study_xxqg)


编译
go env -w CGO_ENABLED=0
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -o study_xxqg ./


启动
nohup ./study_xxqg > ./study_xxqg.log 2>&1 & echo $!>pid.pid