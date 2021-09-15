package dbtest

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/patrick/test-grpc-docker-gitactions/lib/dbop"
	"github.com/patrick/test-grpc-docker-gitactions/lib/env"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DatabaseID *string
	db         *gorm.DB
}

func (s *Suite) SetDatabaseID(id string) {
	if s.DatabaseID != nil {
		s.FailNowf("testing database ID already set", "database_id: %s", s.DatabaseID)
	}

	s.DatabaseID = &id
}

func (s *Suite) SetupDB() {
	s.SetDatabaseID(setup())
	db := s.DB()

	db.LogMode(true)
}

func (s *Suite) TearDownDB() {
	s.DB().Close()
	s.Close()
	if !s.T().Failed() {
		teardown(*s.DatabaseID)
	} else {
		fmt.Printf("\ndatabase \"%s\" is not removed as test failed\n",
			s.GetDatabaseName())
	}
}

func (s *Suite) GetDatabaseName() string {
	return "testdbtest" + *s.DatabaseID
}

// DB returns a singleton gorm.DB instance for this test suite. Unlike
// dbop.DB() which is a process-wise (hence package level in go test),
// this instance is local to the suite, and is automatically closed
// at the end of suite.
func (s *Suite) DB() *gorm.DB {
	if s.db != nil {
		return s.db
	}
	connectionString := strings.Replace(env.DBConnectionString, "testdb", "testdbtest"+*s.DatabaseID, 1)
	db, err := dbop.NewMSDB(connectionString)
	s.Require().Nil(err)
	s.db = db
	return s.db
}

func getTestDBBakFile() string {
	return env.WorkingDIR + `\testdbtest.bak`
}

//TestDB is a singleton wrapper to the MSSQL gorm database object testdbtest.
func TestDB() *gorm.DB {

	if err := createTestDB(""); err != nil {
		panic(err)
	}

	connectionString := strings.Replace(env.DBConnectionString, "testdb", "testdbtest", 1)
	db, err := dbop.NewMSDB(connectionString)
	if err != nil {
		panic(err)
	}

	return db
}

// Close closes the DB connection if it was opened. This is called
// at the end of test suite automatically.
func (s *Suite) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// ReloadProcedures reinstalls all stored procedures without rebuilding
// the database from scratch, so that procedure sources are reflected
// in the test runs.
// func (s *Suite) ReloadProcedures() {
// 	db, err := dbop.NewDB(dbop.DBOptions{LogMode: &[]bool{false}[0]})
// 	defer db.Close()
// 	s.Require().Nil(err)
// 	// s.Require().Nil(migration.InstallFunctions(db))
// }

func teardown(id string) {
	if err := dropTestDB(id); err != nil {
		panic(err)
	}
}

func setup() (id string) {
	id = generateID()
	connectionString := strings.Replace(env.DBConnectionString, "testdb", fmt.Sprintf("testdbtest%s", id), 1)
	if err := createTestDB(id); err != nil {
		panic(err)
	}

	os.Setenv("TESTDB", connectionString)
	return
}

func createTestDB(id string) error {
	connectionString := strings.Replace(env.DBConnectionString, "testdb", "master", 1)
	msdb, err := dbop.NewMSDB(connectionString)
	if err != nil {
		return err
	}
	defer msdb.Close()

	dbName := fmt.Sprintf("testdbtest%s", id)

	if err := msdb.Exec(fmt.Sprintf(`
	IF DB_ID('%s') IS NOT NULL
	EXEC('DROP DATABASE %s;')`, dbName, dbName)).Error; err != nil {
		return err
	}

	if id == "" {
		return msdb.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)).Error
	}

	logFile := env.WorkingDIR + "/" + dbName
	return msdb.Exec(fmt.Sprintf(`
	RESTORE DATABASE %s 
	FROM DISK = '%s/testdbtest.bak'
	WITH MOVE 'testdb' TO '%s.mdf', 
	MOVE 'testdb_log' TO '%s.ldf'
	`, dbName, env.WorkingDIR, logFile, logFile)).Error
}

func dropTestDB(id string) error {
	connectionString := strings.Replace(env.DBConnectionString, "testdb", "master", 1)
	pgdb, err := dbop.NewMSDB(connectionString)
	if err != nil {
		return err
	}

	defer pgdb.Close()
	dbop.CloseDB()

	return pgdb.Exec(fmt.Sprintf(`DROP DATABASE "testdbtest%s"`, id)).Error
}

func generateID() (id string) {
	u := uuid.Must(uuid.NewV4()).String()
	filePath := getCaller(0)

	shortPath := ""
	lastIdx := strings.LastIndex(filePath, "/")
	if lastIdx > -1 {
		shortPath = filePath[lastIdx+1:]
		// lastIdx = strings.LastIndex(filePath[:lastIdx], "/")
		// if lastIdx > -1 {
		// 	shortPath = filePath[lastIdx+1:]
		// }
	}
	if shortPath == "" {
		shortPath = filePath
	}

	if strings.HasSuffix(shortPath, ".go") {
		shortPath = shortPath[:len(shortPath)-3]
	}
	if strings.HasSuffix(shortPath, "_test") {
		shortPath = shortPath[:len(shortPath)-5]
	}

	return shortPath + u[:4]
}

func getCaller(level uint) (name string) {
	// Ask runtime.Callers for up to 10 pcs, including runtime.Callers itself.
	pc := make([]uintptr, level+10)
	n := runtime.Callers(0, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	// Discard "runtime.Callers"
	_, _ = frames.Next()
	// Loop to get frames.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()
		if match, _ := regexp.Match("dbtest/.*.go$", []byte(frame.File)); match {
			continue
		}
		name = frame.File
		// fmt.Printf("- more:%v | %s\n", more, frame.Function)
		if !more || level == 0 {
			break
		}
		level--
	}
	return
}

// ReloadProcedures reinstalls all stored procedures without rebuilding
// the database from scratch, so that procedure sources are reflected
// in the test runs.
// func (s *Suite) ReloadProcedures() {
// 	// db, err := dbop.NewDB(dbop.DBOptions{LogMode: &[]bool{false}[0]})
// 	// defer db.Close()
// 	// s.Require().Nil(err)
// 	s.Require().Nil(migration.InstallFunctions(s.DB()))
// }

// CheckQuery executes the SQL query using s.DB() and compares the tabular
// result with expected. See also FormatQueryResult for the expected string format.
func (s *Suite) CheckQuery(query, expected string) {
	s.T().Helper()

	ex := trimTrailingSpaces(strings.Trim(expected, "\r\n"))
	res := trimTrailingSpaces(FormatQueryResult(s.DB(), query))
	s.Equal(ex, res)
}

// LoadData is a short-hand for dbtest.LoadData()
func (s *Suite) LoadData(t, tableName string) error {
	return LoadData(s.DB(), t, tableName)
}

// LoadData inserts rows into the table, accepting the string
// format of psql explanded display (\x). Use "\pset null '\\N'"
// to distinguish NULL value from empty string.
//
// For example,
//
//   db=> \x
//   Expanded display is on.
//   db=> \pset null '\\N'
//   Null display is "\N".
//   db=> SELECT * FROM table WHERE id IN (1, 2);
//   -[ RECORD 1 ]------+-------------------------------------
//   col1               | foo
//   col2               | \N
//
// then copy and paste this output as string literal in the test code.
//
// There is a special prefix "expr::" for the data, which injects the
// SQL expression directly. This helps use calculated values such as
// md('text')::uuid technique.
func LoadData(tx *gorm.DB, t, tableName string) error {
	items := parseTextFormat(t)
	for _, item := range items {
		stmt := makeStmt(tableName, item)
		params := makeParams(item)
		if err := tx.Exec(stmt, params...).Error; err != nil {
			return err
		}
	}
	return nil
}

// ExpectResult executes the SQL query using s.DB() and compares the tabular
// result with expected. See also FormatQueryResult for the expected string format.
func (s *Suite) ExpectResult(query, expected string) {
	s.T().Helper()

	ex := trimTrailingSpaces(strings.Trim(expected, "\r\n"))
	res := trimTrailingSpaces(FormatQueryResult(s.DB(), query))
	s.Equal(ex, res)
}

type sqlExpr string

const (
	ALIGN_RIGHT = iota
	ALIGN_LEFT
	ALIGN_CENTER
)

type TableOptions struct {
	Aligns []int
}

func makeStmt(tableName string, item map[string]interface{}) string {
	s := "INSERT INTO " + tableName + "("
	keys := []string{}
	placeHolders := []string{}
	for key := range item {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := item[key]
		if expr, ok := val.(sqlExpr); ok {
			placeHolders = append(placeHolders, string(expr))
		} else {
			placeHolders = append(placeHolders, "?")
		}
	}

	s += strings.Join(keys, ", ")
	s += ") VALUES ("
	s += strings.Join(placeHolders, ", ")
	s += ")"
	return s
}

func parseTextFormat(t string) []map[string]interface{} {
	cursor := []byte(t)
	items := []map[string]interface{}{}
	currentItem := map[string]interface{}{}
	for {
		advance, token, _ := bufio.ScanLines(cursor, true)
		if advance == 0 {
			break
		}
		cursor = cursor[advance:]
		if len(strings.TrimSpace(string(token))) == 0 {
			continue
		}
		if match, _ := regexp.Match("\\[ RECORD \\d+ \\]", token); match {
			if len(currentItem) > 0 {
				items = append(items, currentItem)
				currentItem = map[string]interface{}{}
			}
			continue
		}
		splits := strings.SplitN(string(token), "|", 2)
		key := strings.TrimSpace(splits[0])
		val := strings.TrimSpace(splits[1])
		if val == "\\N" {
			currentItem[key] = (*string)(nil)
		} else if strings.HasPrefix(val, "expr::") {
			currentItem[key] = sqlExpr(strings.TrimSpace(val[6:]))
		} else {
			currentItem[key] = val
		}
	}
	if len(currentItem) > 0 {
		items = append(items, currentItem)
	}

	return items
}

func makeParams(item map[string]interface{}) []interface{} {
	params := []interface{}{}
	keys := []string{}
	for key := range item {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if _, ok := item[key].(sqlExpr); ok {
			continue
		}
		params = append(params, item[key])
	}
	return params
}

func trimTrailingSpaces(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}

// FormatQueryResult performs the DB query and returns a string
// formatted in psql-like output.
// Example:
//
//   FormatQueryResult(`
//     SELECT
//       account_no
//       ,contra_account_no
//       ,side
//       ,qty
//       ,price
//       ,commission
//     FROM report.activity_vw
//     ORDER BY external_id, account_no`))
//
// will yield
//
//    account_no | contra_account_no | side |      qty      |  price  | commission
//   ------------+-------------------+------+---------------+---------+------------
//    1200       | 1210              | sell |          -1.0 | 23.3999 |          0
//    1210       | 1200              | buy  |           1.0 | 23.3999 |          0
//    1210       | 923800204         | sell | -0.8974320620 | 23.4001 |          0
//    923800204  | 1210              | buy  |  0.8974320620 | 23.4001 |       0.99
//
// Note there is no leading/trailing empty lines are added. Also notice
// there are a white space at the end of line as a padding of last column.
// When you test the result, you need to make sure the expected string
// follows the same way (which is wanted to fix but later.)
func FormatQueryResult(tx *gorm.DB, query string) string {
	rows, err := tx.DB().Query(query)
	if err != nil {
		return err.Error()
	}
	cols, err := rows.Columns()
	colTypes, err := rows.ColumnTypes()
	var outputs [][]string
	for rows.Next() {
		tuple := make([]string, len(cols))
		receiver := make([]*string, len(cols))
		tuplePointers := make([]interface{}, len(cols))
		for i := range cols {
			tuplePointers[i] = &receiver[i]
		}
		if err := rows.Scan(tuplePointers...); err != nil {
			return err.Error()
		}
		for i := range colTypes {
			if receiver[i] == nil {
				tuple[i] = ""
			} else if strings.ToLower(colTypes[i].DatabaseTypeName()) == "date" {
				// 2021-01-27T00:00:00Z -> 2021-01-27
				tuple[i] = (*receiver[i])[:10]
			} else {
				tuple[i] = *receiver[i]
			}
		}
		outputs = append(outputs, tuple)
	}
	var aligns []int
	for _, typ := range colTypes {
		var align int
		switch strings.ToLower(typ.DatabaseTypeName()) {
		case "int", "numeric", "bigint", "float":
			align = ALIGN_RIGHT
		default:
			align = ALIGN_LEFT
		}
		aligns = append(aligns, align)
	}

	return writeTableString(outputs, cols, TableOptions{Aligns: aligns})
}

func writeTableString(rows [][]string, header []string, opts TableOptions) string {
	writer := &bytes.Buffer{}
	colwidths := make([]int, len(header))
	for i, h := range header {
		colwidths[i] = len(h) + 2
	}
	for _, r := range rows {
		for i, val := range r {
			l := len(val) + 2
			if l > colwidths[i] {
				colwidths[i] = l
			}
		}
	}
	for i, h := range header {
		writer.Write([]byte(center(h, colwidths[i])))
		if i < len(header)-1 {
			writer.Write([]byte("|"))
		}
	}
	writer.Write([]byte("\n"))
	for i := range header {
		sep := make([]byte, colwidths[i])
		for j := 0; j < colwidths[i]; j++ {
			sep[j] = '-'
		}
		writer.Write(sep)
		if i < len(header)-1 {
			writer.Write([]byte("+"))
		}
	}
	writer.Write([]byte("\n"))
	for i, r := range rows {
		for j := range r {
			align := ALIGN_RIGHT
			if len(opts.Aligns) > j {
				align = opts.Aligns[j]
			}
			var s string
			switch align {
			case ALIGN_LEFT:
				s = left(r[j], colwidths[j])
			case ALIGN_RIGHT:
				s = right(r[j], colwidths[j])
			case ALIGN_CENTER:
				s = center(r[j], colwidths[j])
			}
			writer.Write([]byte(s))
			if j < len(colwidths)-1 {
				writer.Write([]byte("|"))
			}
		}
		if i < len(rows)-1 {
			writer.Write([]byte("\n"))
		}
	}
	return strings.Trim(writer.String(), "\r\n")
}

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func left(s string, w int) string {
	return fmt.Sprintf(" %*s", -(w - 1), s)
}

func right(s string, w int) string {
	return fmt.Sprintf("% *s ", w-1, s)
}
