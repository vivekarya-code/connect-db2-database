package db2

import (
	"database/sql"
	"encoding/hex"
	"log"
	"os"
	"protectsecrets"
	"strings"

	_ "github.com/alexbrainman/odbc" //External package for ODBC connection to DB2
	"github.com/spf13/viper"
	//"time"
)

// The Constants
const ODBC string = "odbc"
const dsnHeader string = "DSN="
const uidHeader string = "Uid="
const pwdHeader string = "Pwd="

var dsn string
var uid string
var pwd string

// ConnectDB2 lets you connect you to the AS400 Database via ODBC
func ConnectDB2() (*sql.DB, error) {

	//Get database environment
	serverEnv := os.Getenv("GOENVIRONMENT")
	if serverEnv == " " {
		log.Fatal("DB2 Environment not found/Empty in Operating System Variables")
	}
	serverEnv = strings.ToUpper(strings.TrimSpace(serverEnv))

	//Get the connection string based on environment
	var connString string
	switch {
	case serverEnv == "PROD":
		getConfigs("app-prod")
		connString = dsnHeader + strings.TrimSpace(dsn) + ";" +
			uidHeader + strings.TrimSpace(uid) + ";" +
			pwdHeader + pwd
	case serverEnv == "QA":
		getConfigs("app-qa")
		connString = dsnHeader + strings.TrimSpace(dsn) + ";" +
			uidHeader + strings.TrimSpace(uid) + ";" +
			pwdHeader + pwd
	case serverEnv == "DEV":
		getConfigs("app-dev")
		connString = dsnHeader + strings.TrimSpace(dsn) + ";" +
			uidHeader + strings.TrimSpace(uid) + ";" +
			pwdHeader + pwd
	}

	//Open database connection
	db, err := sql.Open(ODBC, connString)
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 5. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(0)
	// Set the maximum number of concurrently idle connections to 5. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	//db.SetMaxIdleConns(0)
	// //Set Connection Life time
	// db.SetConnMaxLifetime(1 * time.Hour)

	return db, err

}

func getConfigs(configFile string) {

	vi := viper.New()
	vi.SetConfigName(configFile)
	vi.AddConfigPath("C:/Go/config")
	err := vi.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//Data Source Name
	dsn = vi.GetString("dsn")

	//User ID
	uid = vi.GetString("uname")

	//Password
	key := "passPhrase"
	pwdhexvalue, _ := hex.DecodeString(vi.GetString("password"))
	ciphertext := []byte(pwdhexvalue)
	pwd = string(protectsecrets.Decrypt(ciphertext, key))

}
