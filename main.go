package main

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"

	"pii-encrypt-example/cmd/user/v1"
	"pii-encrypt-example/configs"
	"pii-encrypt-example/server"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // for development
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	ddlogrus "gopkg.in/DataDog/dd-trace-go.v1/contrib/sirupsen/logrus"

	"pii-encrypt-example/pkg/crypto"
	"pii-encrypt-example/pkg/hook"
	"pii-encrypt-example/pkg/middleware"
	"pii-encrypt-example/pkg/response"
)

var (
	// tracer       *apm.Tracer
	cfg          *configs.Config
	indexMessage string = "Application is running properly"
)

func init() {
	// tracer = apm.DefaultTracer
	cfg = configs.Load()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)
	logger.AddHook(&ddlogrus.DDContextLogHook{})
	logger.AddHook(hook.NewStdoutLoggerHook(logrus.New(), cfg.Logger.Formatter))

	// set crypto
	crypto := crypto.NewAESGCM(cfg.Crypto.Secret, cfg.Crypto.Pepper)

	// set mariadb read only object
	dbReadOnly, err := sql.Open(cfg.MariadbReadOnly.Driver, cfg.MariadbReadOnly.DSN)
	if err != nil {
		logger.Fatal(err)
	}
	if err := dbReadOnly.Ping(); err != nil {
		logger.Fatal(err)
	}
	dbReadOnly.SetConnMaxLifetime(time.Minute * 3)
	dbReadOnly.SetMaxOpenConns(cfg.MariadbReadOnly.MaxOpenConnections)
	dbReadOnly.SetMaxIdleConns(cfg.MariadbReadOnly.MaxIdleConnections)

	// set mariadb read write object
	dbReadWrite, err := sql.Open(cfg.MariadbReadWrite.Driver, cfg.MariadbReadWrite.DSN)
	if err != nil {
		logger.Fatal(err)
	}
	if err := dbReadWrite.Ping(); err != nil {
		logger.Fatal(err)
	}
	dbReadWrite.SetConnMaxLifetime(time.Minute * 3)
	dbReadWrite.SetMaxOpenConns(cfg.MariadbReadWrite.MaxOpenConnections)
	dbReadWrite.SetMaxIdleConns(cfg.MariadbReadWrite.MaxIdleConnections)

	router := mux.NewRouter()
	router.HandleFunc("/todo", index)

	basicAuthMiddleware := middleware.NewBasicAuth(cfg.BasicAuth.Username, cfg.BasicAuth.Password)

	// set validator
	validator := validator.New()
	// validator.RegisterTagNameFunc(customvalidator.SetTagName)
	// validator.RegisterValidation("default-name", customvalidator.SetDefaultName)
	// validator.RegisterValidation("idn-mobile-number", customvalidator.SetIDNMobileNumber)
	// validator.RegisterValidation("ISO8601date", customvalidator.SetISO8601dateFormat)

	userRepository := user.NewUserRepository(logger, dbReadOnly, dbReadWrite, "user_encrypt")
	userUsecase := user.NewUserUsecase(logger, cfg.Application.Timezone, crypto, userRepository)
	user.NewUserHTTPHandler(logger, router, basicAuthMiddleware, validator, userUsecase)

	handler := middleware.ClientDeviceMiddleware(router)
	// set cors
	handler = cors.New(cors.Options{
		AllowedOrigins:   cfg.Application.AllowedOrigins,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
	}).Handler(handler)
	handler = middleware.NewRecovery(logger, true).Handler(handler)

	// initiate server
	srv := server.NewServer(logger, handler, cfg.Application.Port)
	srv.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	srv.Close()
	dbReadOnly.Close()
	dbReadWrite.Close()

}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
