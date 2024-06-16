package cmd

import (
    "flag"
    "fmt"
    "io"
    "os"
    "strconv"

    storeModule "store-api/internal/store/handler"
    storeRepo "store-api/internal/store/repository"
    storeService "store-api/internal/store/service"

    "gopkg.in/Graylog2/go-gelf.v2/gelf"

    fcmToken "store-api/internal/base/service/firebase"
    cache "store-api/internal/base/service/redisser"

    gelfFormatter "github.com/seatgeek/logrus-gelf-formatter"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cast"

    "store-api/internal/base/handler"

    "store-api/pkg/db"
    "store-api/pkg/httpclient"
    "store-api/pkg/metric"
)

var (
    params map[string]string

    baseHandler  *handler.BaseHTTPHandler
    storeHandler *storeModule.HTTPHandler
    httpClient   httpclient.Client

    mysqlClientRepo  *db.MySQLClientRepository
    redisClient      cache.RedisClient
    firebaseClient   fcmToken.FirebaseClient
    statsdMonitoring metric.StatsdMonitoring
)

func initMySQL() {
    host := os.Getenv("DB_HOST")
    port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
    dbname := os.Getenv("DB_NAME")
    uname := os.Getenv("DB_USERNAME")
    pass := os.Getenv("DB_PASSWORD")

    mysqlClientRepo, _ = db.NewMySQLRepository(host, uname, pass, dbname, port)
}

func initInfrastructure() {
    initMySQL()
    initLog() // Init log after baseHandler
    httpClientFactory := httpclient.New()
    httpClient = httpClientFactory.CreateClient()

    var err error
    if err != nil {
        logrus.Errorln(err)
    }
}

func isProd() bool {
    return os.Getenv("APP_ENV") == "production"
}

func initHTTP() {
    params = initParams()
    initInfrastructure()

    params["mysql_tz"] = mysqlClientRepo.TZ

    storeRepo := storeRepo.NewStoreRepository(mysqlClientRepo.DB)

    storeService := storeService.NewService(storeRepo)

    baseHandler = handler.NewBaseHTTPHandler(mysqlClientRepo.DB, httpClient, params, statsdMonitoring, storeService)

    storeHandler = storeModule.NewHTTPHandler(baseHandler, storeService)

    fmt.Println("INFO: Init and load module completed. Server started.\n---")
}

func initLog() {
    logrus.SetFormatter(&gelfFormatter.GelfFormatter{})

    checkGraylog := cast.ToBool(os.Getenv("USE_GRAYLOG"))
    if checkGraylog {
        var graylogAddr string

        flag.StringVar(&graylogAddr, "graylog", os.Getenv("GREYLOG_HOST"), "graylog server addr")
        flag.Parse()

        if graylogAddr != "" {
            gelfWriter, err := gelf.NewTCPWriter(graylogAddr)
            if err != nil {
                logrus.Fatal("gelf.NewWriter: %w", err)
            }

            // log to both stderr and graylog2
            logrus.SetOutput(io.MultiWriter(os.Stdout, gelfWriter))
            logrus.Info("Graylog is running")
        }
    }

    lv := os.Getenv("LOG_LEVEL_DEV")
    level := logrus.InfoLevel
    switch lv {
    case "PanicLevel":
        level = logrus.PanicLevel
    case "FatalLevel":
        level = logrus.FatalLevel
    case "ErrorLevel":
        level = logrus.ErrorLevel
    case "WarnLevel":
        level = logrus.WarnLevel
    case "InfoLevel":
        level = logrus.InfoLevel
    case "DebugLevel":
        level = logrus.DebugLevel
    case "TraceLevel":
        level = logrus.TraceLevel
    default:
    }

    // Only log above level
    if isProd() {
        // Only Warn and Error log for prod
        logrus.SetLevel(logrus.WarnLevel) // Set default InfoLevel
    } else {
        // Set log level for staging.
        if lv == "" && os.Getenv("APP_DEBUG") == "True" {
            level = logrus.DebugLevel
        }
        logrus.SetLevel(level)
    }
}
