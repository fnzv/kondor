package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/telegram-bot-api.v4"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Id     string `json:"Id"`
	Camera string `json:"Camera"`
	Label  string `json:"Label"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func stringToTime(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

var mysql_conn string
var bot_token string
var chat_id int64
var frigate_url string

func main() {
	chat_id, err := strconv.ParseInt(os.Getenv("TGBOT_CHATID"), 10, 64)
	if err != nil {
		panic(err)
	}

	if os.Getenv("FRIGATE_URL") != "" {
		frigate_url = os.Getenv("FRIGATE_URL")
	} else {
		fmt.Println("FRIGATE_URL ENV var not set")
	}

	if os.Getenv("MYSQL_CONN") != "" {
		mysql_conn = os.Getenv("MYSQL_CONN")
	} else {
		fmt.Println("MYSQL_CONN ENV var not set")
	}

	res, err := http.Get(frigate_url + "/api/events")
	if err != nil {
		fmt.Println(err)
		return
	}

	db, connerr := sql.Open("mysql", mysql_conn)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `events` (`id` varchar(60) NOT NULL,`camera` varchar(45) DEFAULT NULL,`label` varchar(45) DEFAULT NULL,`time` varchar(70) DEFAULT NULL,`sent` varchar(45) DEFAULT NULL,PRIMARY KEY (`id`), UNIQUE KEY `id_UNIQUE` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")

	fmt.Println("Notify started, creating DB schema if not existing")
	if err != nil {
		panic(err)
		fmt.Println("Create schema failed")
	}
	checkErr(connerr)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	dataJson := string(body)
	var arr []Event
	_ = json.Unmarshal([]byte(dataJson), &arr)
	for i := range arr {
		label := arr[i].Label
		id := arr[i].Id
		camera := arr[i].Camera
		res1 := strings.Split(id, "-")
		timestamp_float := res1[0]
		timestamp := strings.Split(timestamp_float, ".")[0]
		date, _ := stringToTime(timestamp)

		stmt, dberr := db.Prepare("INSERT IGNORE events SET id=?, camera=?, label=?, time=?")
		checkErr(dberr)
		if _, err := stmt.Exec(id, camera, label, date); err != nil {
			fmt.Println("Falied to insert event into DB: ", err, id)
		}
		bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_TOKEN"))
		if err != nil {
			log.Panic(err)
		}

		rows, err := db.Query("SELECT id FROM events WHERE sent IS NULL AND  id='" + id + "'")
		checkErr(err)
		for rows.Next() {
			err = rows.Scan(&id)
			checkErr(err)
			if os.Getenv("TGBOT_TOKEN") != "" && os.Getenv("TGBOT_CHATID") != "" {
				notification := "\360\237\224\224 Frigate notification \360\237\224\224 \n\n<b>Label</b>: " + label + " \n<b>Time</b>:" + date.String() + " \n<b>Camera</b>: " + camera + " \n<b><a href='" + frigate_url + "/events'>Open UI</a></b>"
				msg := tgbotapi.NewMessage(chat_id, notification)
				msg.ParseMode = "HTML"
				bot.Send(msg)

				fileName := "/tmp/" + id + ".jpg"
				URL := frigate_url + "/api/events/" + id + "/thumbnail.jpg"
				err := downloadFile(URL, fileName)
				if err != nil {
					log.Fatal(err)
				}

				photoBytes, err := ioutil.ReadFile(fileName)
				if err != nil {
					panic(err)
				}
				photoFileBytes := tgbotapi.FileBytes{
					Name:  "picture",
					Bytes: photoBytes,
				}
				message, err := bot.Send(tgbotapi.NewPhotoUpload(chat_id, photoFileBytes))
				fmt.Println("NOTIFY SENT: ", message)
				fmt.Println("NOTIFY SENT: " + notification)
				_, err = db.Query("UPDATE events SET sent='true' WHERE id='" + id + "'")
				checkErr(err)
			} else {

				fmt.Println("Telegram env vars are not set... skipping notification")
			}
		}

	}

}
