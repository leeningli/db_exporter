package main

import (
	"leeconfig"
	"fmt"
	"strings"
	"os"
	"net/http"
	"io"
	"database/sql"
	_ "github.com/wendal/go-oci8"
)
const (
	MAX_CNT_TOPIC int = 100
)
var DB_TOPICS = [MAX_CNT_TOPIC]string{""}

func readConfig() {
	fmt.Println("start read config field:main...")
	TOPIC := leeconfig.GetConfig("main")
	topics := TOPIC["topics"]
	topic_list := strings.Split(topics, ",")
	for k, db_topic := range topic_list {
		DB_TOPICS[k] = db_topic
	}
}

func init() {
	readConfig()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}


func oracle_exporter(appname string) (string){
	var metrics string = ""
	TOPIC_ORACLE := leeconfig.GetConfig(appname)
	ip := TOPIC_ORACLE["ip"]
	port := TOPIC_ORACLE["port"]
	username := TOPIC_ORACLE["username"]
	pwd := TOPIC_ORACLE["pwd"]
	sid := TOPIC_ORACLE["sid"]
	cmd := TOPIC_ORACLE["cmd"]
	fmt.Println(ip,port,username,pwd)
	sql_str := fmt.Sprintf("%s/%s@%s", username, pwd, sid)
	db, err := sql.Open("oci8", sql_str)
	checkError(err)
	rows, err := db.Query(cmd)
	for rows.Next() {
		var tag string
		var value int
		err = rows.Scan(&tag,&value)
		if err == nil {
			res := fmt.Sprintf("%s {host=%s,port=%s} %d\n", tag, ip, port, value)
			metrics = metrics + res
		}
	}
	defer db.Close()
	return metrics
}

func ExporterHandler(w http.ResponseWriter, r *http.Request) {
	for _, value := range DB_TOPICS {
		if value != "" {
			db_index := strings.Split(string(value), ":")[0]
			fmt.Println("db_index==", db_index)
			if strings.ToLower(db_index) == "oracle" {
				metrics := oracle_exporter(value)
				io.WriteString(w, metrics)
			} else {
				fmt.Println("config:", value, " is error.")
				os.Exit(1)
			}
		}
	}
}

func main() {
	port := flag.String("port", "30083", "Input your exporter port")
	flag.Parse()
	addrs, err := net.InterfaceAddrs()
	checkError(err)
	var host string
	for _, a := range addrs {
		if hostip, ok := a.(*net.IPNet); ok && !hostip.IP.IsLoopback() {
			if hostip.IP.To4() != nil {
				fmt.Println(hostip.IP.String())
				host = hostip.IP.String()
				break;
			}
		}
	}
	http.HandleFunc("/metrics", ExporterHandler)
	url := fmt.Sprintf("%s:%s", host, *port)
	fmt.Println("url==", url, "/metrics")
	err = http.ListenAndServe(url, nil)
	checkError(err)
}
