package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type sensordata struct {
	ID         int
	Devicename string
	Temprature float64
	IsActive   bool
}

func main() {
	conn := "postgres://postgres:3020@localhost:5432/gopgtest?sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal("Error in opening the connection: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Connection not validated since request ain't hitting the Db\n", err)
	}
	createTable(db)
	id := insertData(db)
	fmt.Printf("The id of newly inserted data is: %d\n", id)

	data := getSensorDataByID(db, id)
	fmt.Printf("Pura data: %v\n", data) // %v format mai struct ka data field ka name automatically print krdega

	allsensors, errr := getAllsensordata(db)
	if errr != nil {
		log.Fatal("Data laane me error: ", errr)
	}
	fmt.Println("\\\\\\\\\\\\Printintg the data we received\\\\\\\\")
	for index, value := range allsensors {
		fmt.Printf("index:%d, \n", index+1)
		fmt.Printf("The data: %v, \n", value)
	}
}

func createTable(db *sql.DB) { //here i learned that when a string is breaked into multi-line we use smthng like ``
	query := `CREATE TABLE IF NOT EXISTS sensordata (
    id SERIAL PRIMARY KEY,
    devicename VARCHAR(100) NOT NULL,              
    temperature NUMERIC(5, 2) NOT NULL,
    isactive BOOL
);`

	_, err := db.Exec(query) //here db.Exec won't return me the row
	if err != nil {
		log.Printf("table bnane mai glti hui hai\n")
	}

}
func insertData(db *sql.DB) int {

	var id int
	query := `INSERT INTO sensordata (devicename, temperature, isactive) 
VALUES ($1, $2, $3) RETURNING id;`
	err := db.QueryRow(query, "motion_tracker", 53, true).Scan(&id)
	if err != nil {
		log.Fatal("the error in logging the data:\n", err)
	}
	return id
}

func getSensorDataByID(db *sql.DB, id int) sensordata {
	query := `SELECT devicename, temperature, isactive FROM sensordata WHERE id = $1;`
	var data sensordata
	err := db.QueryRow(query, id).Scan(&data.Devicename, &data.Temprature, &data.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Errors due to no row present with the row id: %d", id)
			return data
		}
		log.Print("there's an error", err)
	}
	return data
}

func getAllsensordata(db *sql.DB) ([]sensordata, error) {

	query := `SELECT id, devicename, temperature, isactive FROM sensordata;`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("the error is: \n", err)

	}
	defer rows.Close()
	var allData []sensordata
	for rows.Next() {
		var u sensordata
		err := rows.Scan(&u.ID, &u.Devicename, &u.Temprature, &u.IsActive)
		if err != nil {
			log.Fatal("scan error: %w", err)
		}
		allData = append(allData, u)

	}
	return allData, err
}
