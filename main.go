package main

import (
	"GoBlog/global"
	"GoBlog/internal/model"
	"GoBlog/internal/routers"
	"GoBlog/pkg/logger"
	"GoBlog/pkg/setting"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
)


// main 方法之前自动执行init ,
// 执行顺序是：全局变量初始化 =》init 方法 =》main 方法
func init() {

	// 初始化配置，达到配置文件内容映射到应用配置结构体的作用
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	// 将DBEngine初始化
	err = setupDBEngine()
	if err != nil{
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}


func main() {
	gin.SetMode(global.ServerSetting.RunMode)

	router := routers.NewRouter()

	s:=&http.Server{
		Addr: ":" + global.ServerSetting.HttpPort,
		Handler: router,
		ReadTimeout: global.ServerSetting.ReadTimeout,
		WriteTimeout: global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	//global.Logger.Infof("%s: go_blog/%s", "第一篇文章", "blog-service")


	s.ListenAndServe()

}


func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

// 将全局变量DBEngine初始化
func setupDBEngine() error {
	var err error

	// 这里不能用:= 因为这个符号相当于新创建一个局部变量，全局变量不没有被赋值
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}


func setupLogger() error {
	//使用了 lumberjack 作为日志库的 io.Writer，
	//并且设置日志文件所允许的最大占用空间为 600MB、
	//日志文件最大生存周期为 10 天，并且设置日志文件名的时间格式为本地时间。
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}