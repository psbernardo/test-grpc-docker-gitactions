package env

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	GOClearingGRPCPort     = GetValue("MSCLEARING_GRPC_PORT")
	GOClearingGRPCWebPort  = GetValue("MSCLEARING_GRPC_WEB_PORT")
	DBConnectionString     = GetValue("TESTDB")
	DBTestConnectionString = strings.Replace(GetValue("TESTDB"), "testdb", "testdbtest", 1)
	SecretKey              = GetValue("GO_CLEARING_SECRET_KEY")   //SecretKey for encryption i.e. tax id
	JwtAudience            = GetValue("GO_CLEARING_JWT_AUDIENCE") // must  match with UI/Client's client id
	JwtSubject             = GetValue("GO_CLEARING_JWT_SUBJECT")
	JwtIssuer              = GetValue("GO_CLEARING_JWT_ISSUER")
	ApplicationMode        = GetValue("GO_CLEARING_MODE")                 //ApplicationMode (STG, PROD, SANDBOX, PAPER_PROD)
	TokenDuration          = GetInt("GO_CLEARING_TOKEN_DURATION_MINUTES") //TokenDuration for generation access token
	Debug                  = GetBool("GO_CLEARING_DEBUG", true)
	GRPCMaxMsgSize         = GetInt("GO_CLEARING_GRPC_MAX_MSG_SIZE")
	AttachmentsPath        = GetValue("GO_CLEARING_ATTACHMENTS_PATH")
	WorkingDIR             = GetValue("GO_CLEARING_Working_DIR") //WorkingDIR use for saving temporary files i.e. test database
	CasBinConfigPath       = GetValue("CASBIN_CONFIG_PATH")      // CasBinConfigPath get enviroment variable
	SlackBotToken          = GetValue("SLACK_BOT_TOKEN")
	AppPublishPath         = GetValue("GO_CLEARING_PUBLISH_PATH")
	PdfTemplatePath        = GetValue("GO_CLEARING_PDF_TEMPLATE_PATH")

	GoogleApplicatoinCredential  = GetValue("GOOGLE_APPLICATION_CREDENTIALS")
	PubsubCredential             = GetValue("GO_CLEARING_PUBSUB_APPLICATION_CREDENTIALS")
	PubsubNotificationCredential = GetValue("GO_CLEARING_PUBSUB_NOTIFICATION_CREDENTIALS")
	ApiConfiguration             = GetValue("GO_CLEARING_API_CONFIGURATION")
)

var (
	BMODailyStatementFilePrefix = GetValue("BMO_DAILY_STATEMENT_FILE_PREFIX")
	BMOSFTPHost                 = GetValue("BMO_SFTP_HOST")
	BMOSFTPPort                 = GetValue("BMO_SFTP_PORT")
	BMOSFTPUser                 = GetValue("BMO_SFTP_USER")
	BMOSFTPPass                 = GetValue("BMO_SFTP_PASS")

	FRBDailyStatementFilePrefix = GetValue("FRB_DAILY_STATEMENT_FILE_PREFIX")
	FRBSFTPHost                 = GetValue("FRB_SFTP_HOST")
	FRBSFTPPort                 = GetValue("FRB_SFTP_PORT")
	FRBSFTPUser                 = GetValue("FRB_SFTP_USER")
	FRBSFTPPass                 = GetValue("FRB_SFTP_PASS")

	KnownHostFile = GetValue("KNOWN_HOSTS_FILE")

	AWSAccessKeyID     = GetValue("AWS_ACCESS_KEY_ID")
	AWSSecretAccessKey = GetValue("AWS_SECRET_ACCESS_KEY")
	AWSS3BucketName    = GetValue("AWS_S3_BUCKETNAME")
	AWSS3NameSpace     = GetValue("AWS_S3_NAMESPACE")
	StartTime          = GetValue("START_TIME")
	LogDB              = true
	DBMaxOpenConns     = GetInt("DB_MAX_OPEN_CONNS")
	DBMaxIdleConns     = GetInt("DB_MAX_IDLE_CONNS")
)

var dVal sync.Map

//  RegisterDefault by key and defaultValue
func RegisterDefault(key, defaultValue string) {
	dVal.Store(key, defaultValue)
}

// GetValue by key
func GetValue(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if v, _ := dVal.Load(key); v != nil {
			return v.(string)
		} else {
			return ""
		}
	}
	return value
}

// GetBool returns a true if the environment variable is
// "true", "t", "yes" or "y" with case insensitive.
func GetBool(key string, defaultValue bool) bool {
	strval := GetValue(key)
	if strval == "" {
		return defaultValue
	} else if strings.EqualFold(strval, "true") ||
		strings.EqualFold(strval, "t") ||
		strings.EqualFold(strval, "yes") ||
		strings.EqualFold(strval, "1") ||
		strings.EqualFold(strval, "y") {
		return true
	}

	return false
}

// GetInt returns the integer value stored in the enviroment variable.
// When the key is missing, or it cannot be parsed to int, it returns 0
// this function DOES NOT return a default value when a missing key, or non-int value is encountered!
func GetInt(key string) int {
	v, err := strconv.Atoi(GetValue(key))
	if err != nil {
		return 0
	}

	return v
}
