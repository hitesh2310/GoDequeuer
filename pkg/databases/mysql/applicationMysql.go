package database

import (
	"database/sql"
	"fmt"
	"main/logs"
	"main/pkg/constants"

	_ "github.com/go-sql-driver/mysql"
)

var DbConn *sql.DB

func EstablishDbConnection() {
	logs.InfoLog("Establishing DB connection")
	host := constants.ApplicationConfig.Database.Host
	port := constants.ApplicationConfig.Database.Port
	username := constants.ApplicationConfig.Database.Username
	password := constants.ApplicationConfig.Database.Password
	dbname := constants.ApplicationConfig.Database.Dbname
	// fmt.Println("HEREE")
	// fmt.Println(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	logs.InfoLog(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		logs.InfoLog(err.Error())
	}
	// defer db.Close()
	logs.InfoLog("Conn opened")
	// Check the database connection
	logs.InfoLog("Ping check ")
	err = db.Ping()
	logs.InfoLog("Ping checked")

	if err != nil {
		// logs.InfoLog(err.Error())
		logs.ErrorLog("Not connected %v", err)
		return
	}

	DbConn = db
	logs.InfoLog("Connected to the database")

}

func GetClientDetail(phonenumber string) map[string]string {
	if DbConn == nil {
		EstablishDbConnection()
	}

	query1 := fmt.Sprintf("SELECT * from clients where  concat(country_code, '', phone_number ) = '%s'", phonenumber)
	rows, err := DbConn.Query(query1)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	data := make(map[string]string)

	// Check if the phone number exists in the database
	if rows.Next() {
		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			logs.InfoLog("Error to get Columns")
			return nil
		}

		// Create a slice for values
		values := make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(sql.RawBytes)
		}

		// Scan the row into the map
		err = rows.Scan(values...)
		if err != nil {

			return nil
		}

		// Populate the map
		for i, column := range columns {
			val := *(values[i].(*sql.RawBytes))
			data[column] = string(val)
		}
	}

	// Return the map
	return data
}
