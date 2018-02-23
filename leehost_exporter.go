package main

import(
	"fmt"
	"net/http"
	"net"
	"os"
	"flag"
	"log"
	"bytes"
	"os/exec"
	"io"
	"encoding/json"
)

type CHECKS struct {
	Http string `json:"http"`
	Interval string `json:"interval"`
}

type Consul struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
	Port int `json:"port"`
	Tag string `json:"tag"`
	Checks CHECKS `json:"checks"`
}

func RegisterConsul(id, servicename_consul, ip, port, tag, interval, consul_register_url string) error {
	metrics := fmt.Sprintf("http://%s:%s/metrics", ip, port)
	port_int, err := strconv.Atoi(port)
	if err != nil {
	        log.Fatal(err)
	}
	cmd := fmt.Sprintf("curl -X PUT -d '{\"id\": \"%s\",\"name\": \"%s\",\"address\": \"%s\",\"port\": %d,\"tags\": [\"%s\"],\"checks\":[{\"http\":\"%s\",\"interval\": \"%s\"}]}'", id, servicename_consul, ip, port_int, tag, metrics, interval)
	cmd = fmt.Sprintf("%s %s", cmd, consul_register_url)
	fmt.Println("cmd==", cmd)
	fmt.Println(exec_shell(cmd))
	return nil
}

func get_hostip() (string){
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var host string = "127.0.0.1"
	for _, a := range addrs {
		if hostip, ok := a.(*net.IPNet); ok && !hostip.IP.IsLoopback() {
			if hostip.IP.To4() != nil {
					fmt.Println(hostip.IP.String())
					host = hostip.IP.String()
					break
				}
		}
	}
	return host
}

func exec_shell(s string) (string){
    cmd := exec.Command("/bin/bash", "-c", s)
    var out bytes.Buffer

    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    return out.String()
}

func ExporterHandler(w http.ResponseWriter, r *http.Request){
	host := get_hostip()
	var metrics_url string = ""
	io.WriteString(w, "#this is monitor for host, dev by lee\n")
    	uptime_cmd := `uptime|awk '{print $11}'|awk -F"," '{print $1}'`
    	metrics_url = fmt.Sprintf("linux_uptime_load{ip=\"%s\",time=\"5min\"} %s\n", host, exec_shell(uptime_cmd))
    	io.WriteString(w, metrics_url)

    	logstash_cmd := `ps -ef|grep -v grep|grep logstash|awk '{print $2}'|wc -l`
    	metrics_url = fmt.Sprintf("logstashIsExist{host=\"%s\"} %s\n", host, exec_shell(logstash_cmd))
   	io.WriteString(w, metrics_url)

    	flume_cmd := `ps -ef|grep -v grep|grep flume|awk '{print $2}'|wc -l`
    	metrics_url = fmt.Sprintf("flumeIsExist{host=\"%s\"} %s\n", host, exec_shell(flume_cmd))
    	io.WriteString(w, metrics_url)

    	zabbix_cmd := `ps -ef|grep -v grep|grep zabbix|awk '{print $2}'|wc -l`
    	metrics_url = fmt.Sprintf("zabbixIsExist{host=\"%s\"} %s\n", host, exec_shell(zabbix_cmd))
    	io.WriteString(w, metrics_url)

    	io.WriteString(w, "#this is monitor for salt process is exist\n")
    	salt_cmd := `ps -ef|grep -v grep|grep salt|awk '{print $2}'|wc -l`
    	metrics_url = fmt.Sprintf("saltIsExist{host=\"%s\"} %s\n", host, exec_shell(salt_cmd))
    	io.WriteString(w, metrics_url)

    	io.WriteString(w, "#this is monitor for openfiles \n")
    	openfile_cmd := `lsof |awk '{print $3}' |grep -v "TID" |sort |uniq -c|awk '{print "user_open_files{host=\"` + host +`\",user=\""$2"\"}\t"$1}'`
    	metrics_url = fmt.Sprintf("%s", exec_shell(openfile_cmd))
    	io.WriteString(w, metrics_url)

	/*
	io.WriteString(w, "#this is monitor for frontend-app connect to redis,only list connect_cnt>10;if your host is redis host,forget it.\n")
	frontend_redis_cmd := `netstat -pan|grep 6379|awk '{print $7}'|awk -F"/" '{print $1}'|sort |uniq -c|awk '{print "pid_redis_cnt{host=\"` +host +`\",pid=\""$2"\"}\t"$1}'`
	//frontend_redis_cmd := `netstat -pan|grep 6379|awk '{print $7}'|awk -F"/" '{print $1}'|sort |uniq -c|awk '{print "pid_redis_cnt{pid=\""$2"\"}\t"$1}'`
	metrics_url = fmt.Sprintf("%s", exec_shell(frontend_redis_cmd))
	io.WriteString(w, metrics_url)

	io.WriteString(w, "#this is monitor for redis connect status on redis host,only list connect_cnt>10;if your host is not redis host,forget it.\n")
	redis_cmd := `netstat -an|grep 6379|awk '{print $5}' |awk -F":" '{print $1}'|grep -v "0.0.0.0"|sort|uniq -c |awk '{if($1>10) print "redis_connect_cnt{remote_host=\""$2"\"}\t"$1}'`
	metrics_url = fmt.Sprintf("%s", exec_shell(redis_cmd))
	io.WriteString(w, metrics_url)*/

	io.WriteString(w, "#this is monitor for disk io on host.only tps>10\n")
	diskio_cmd := `iostat -d -k|awk '{if($2>10) print $0}'|grep -v "Device"|grep -v "Linux"|grep -v "sdb"|awk '{printf "disk_io_tps{disk=\"%s\",host=\"%s\"}\t%f\n",$1,"`+host+`",$2}'`
	//fmt.Println(diskio_cmd)
	metrics_url = fmt.Sprintf("%s", exec_shell(diskio_cmd))
	io.WriteString(w, metrics_url)
}

func main(){
	port := flag.String("port", "30083", "Input your exporter port")
	flag.Parse()
	host := get_hostip()
    	fmt.Println("port==", *port)
	RegisterConsul("yourid","yourname", host, *port, "yourtag", "15s", "http://your_consul_ip:your_consul_port/v1/agent/service/register")
	http.HandleFunc("/metrics", ExporterHandler)
	url := fmt.Sprintf("%s:%s", host, *port)
	fmt.Println("url=", url+"/metrics")
	err := http.ListenAndServe(url, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
