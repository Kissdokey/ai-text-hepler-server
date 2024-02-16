package mysql

import (
	"ai-text-helper-server/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 传入表参数
type SQLData struct {
	Username         string `json:"username"`
	Avatar           string `json:"avatar"`
	ApiRequestNumber int    `json:"apiRequestNumber"`
	FileId           string `json:"fileId"`
	Belongs          string `json:"belongs"`
	Content          string `json:"content"`
	Comments         string `json:"comments"`
	LastModifiedTime string `json:"lastModifiedTime"`
}

// 字段约束描述map
var fieldDescription = map[string]string{
	"username":         " varchar(64) NOT NULL",
	"avatar":           " longtext ",
	"apiRequestNumber": " int(4) NOT NULL DEFAULT 10",
	"fileId":           " varchar(64) NOT NULL",
	"belongs":          " varchar(64) NOT NULL  DEFAULT ''",
	"content":          " longtext ",
	"comments":         " longtext ",
	"lastModifiedTime": " date",
}

// 每个表字段map
var tables = map[string][]string{
	"user_info_table": {"username", "avatar", "apiRequestNumber"},
	"file_table":      {"fileId", "belongs", "content", "comments", "lastModifiedTime"},
}

// 创建表指令集map
var createTableCommands = map[string]string{}

// 数据库连接的配置
const (
	username = "root"
	hostname = "127.0.0.1:3306"
	dbname   = "ai_text_helper"
)

// 数据库对象
var SqlDb *sql.DB

// 初始化指令集函数
func initCommands() {
	//初始化创建表指令集
	var command string
	for table, fields := range tables {
		command = "CREATE TABLE `" + table + "` ("
		for _, field := range fields {
			command += fmt.Sprintf("`%s` %s,", field, fieldDescription[field])
		}
		command = strings.TrimRight(command, ",")
		command += ");"
		createTableCommands[table] = command
	}
}

// 生成连接指令函数
func formConnectCommand(passWord string, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", username, passWord, hostname, dbName)
}

// 生成插入指令
func formInserCommand(data SQLData, insertType string) (string, error) {
	var (
		fieldArray []string
		ok         bool
		command    string
	)
	if fieldArray, ok = tables[insertType]; !ok {
		return "", errors.New("生成插入指令失败，未找到插入表")
	}
	command = "insert into " + insertType + " values("
	dataValues := reflect.ValueOf(data)
	for _, value := range fieldArray {
		field := strings.ToUpper(value[:1]) + value[1:]
		switch dataValues.FieldByName(field).Interface().(type) {
		case string:
			command += ("'" + dataValues.FieldByName(field).Interface().(string) + "'")

		case int:
			command += strconv.Itoa(dataValues.FieldByName(field).Interface().(int))
		}
		command += ","
	}
	command = strings.TrimRight(command, ",")
	command += ");"
	return command, nil
}

// 创建表
func createTable(createTableStr string) {
	_, err := SqlDb.Exec(createTableStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("表创建成功")
}

// 判断指定表是否存在，这里默认使用ai_text_helper database
func ifTableExist(tableName string) bool {
	var result string
	err := SqlDb.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dbname, tableName).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("table does not exist")
		} else {
			fmt.Println(err.Error())
		}
		return false
	} else {
		return true
	}
}

//外部调用函数

// 初始化数据库，包括连接，判断表存在与否新建表。只需要确保数据库已经允许在主机上，且配置变量和环境变量中的MYSQL_PASSWORD正确
func InitMySQL() {
	//获取数据库密码
	passWord, passwordErr := utils.GetEnvSQL()
	if passwordErr != nil {
		fmt.Println("数据库密码获取出现问题", passwordErr)
		return
	}
	var err error
	//尝试连接到连接池基
	SqlDb, err = sql.Open("mysql", formConnectCommand(passWord, ""))
	if err != nil {
		fmt.Println("数据库打开出现了问题", err)
		return
	}
	// 2.测试数据库是否连接成功
	err = SqlDb.Ping()
	if err != nil {
		fmt.Println("数据库连接出现了问题", err)
		return
	}
	//下面是创建指定数据库并重新连接
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := SqlDb.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n%s\n", err, res)
		return
	}
	SqlDb.Close()
	SqlDb, err = sql.Open("mysql", formConnectCommand(passWord, dbname))
	if err != nil {
		fmt.Println("数据库打开出现了问题", err)
		return
	}
	err = SqlDb.Ping()
	if err != nil {
		fmt.Println("数据库连接出现了问题", err)
		return
	}
	fmt.Println("数据库连接成功", dbname)
	//初始化指令map
	initCommands()
	//依次判断需要的表是否存在，如不存在就创建
	for table := range tables {
		if exist := ifTableExist(table); !exist {
			createTable(createTableCommands[table])
		}
	}
}

// 往table插入数据
func InsertRecord(data SQLData, table string) error {
	insertCommand, error := formInserCommand(data, table)
	fmt.Println(insertCommand)
	if error != nil {
		fmt.Println(error)
		return error
	}
	_, err := SqlDb.Exec(insertCommand)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("插入成功")
	return nil
}

func UpdateFile(data SQLData, fileId string) error {
	deleteCommand := "delete  from file_table where fileId = " + fileId + ";"
	_, err := SqlDb.Exec(deleteCommand)
	if err != nil {
		fmt.Println(err)
	}
	insertErr := InsertRecord(data, "file_table")
	if insertErr != nil {
		fmt.Println(insertErr)
		return insertErr
	}
	return nil
}
func UpdateAvatar(avatar string, username string) error {
	command := fmt.Sprintf("update user_info_table set avatar = '%s' where username = '%s';", avatar, username)
	_, err := SqlDb.Exec(command)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func UpdateRequestNumber(username string, deltaNum int) (int,error) {
	selectCommand := "select apiRequestNumber from user_info_table where username = ?;"
	numRows := SqlDb.QueryRow(selectCommand,username)
	var num int
	selectErr := numRows.Scan(&num)
	if selectErr != nil {
		fmt.Println(selectErr)
		return 0,selectErr
	}
	num += deltaNum
	if num < 0 {
		num = 0
	}
	command := fmt.Sprintf("update %s set apiRequestNumber = %d where username = '%s';", "user_info_table", num, username)
	_, err := SqlDb.Exec(command)
	if err != nil {
		fmt.Println(err)
		return 0,err
	}
	return num,nil
}

func GetUserInfo(username string) (SQLData, error) {
	command := "select avatar,apiRequestNumber from user_info_table where username = ?"
	infoRows := SqlDb.QueryRow(command,username)
	var (
		avatar           string
		apiRequestNumber int
	)
	err :=infoRows.Scan(&avatar, &apiRequestNumber)
	if err != nil {
		fmt.Println(err)
		return SQLData{}, err
	}
	return SQLData{Avatar: avatar, ApiRequestNumber: apiRequestNumber},nil
}
