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
	metrics_url = fmt.Sprintf("uptime{host=\"%s\"} %s\n", host, exec_shell(uptime_cmd))
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
	current_user := exec_shell("whoami")
	openfile_cmd := `lsof |awk -v user=${current_user} '{if($3==user) print $0}'|wc -l`
	metrics_url = fmt.Sprintf("openfiles{host=\"%s\",user=\"%s\"} %s\n", host, current_user, exec_shell(openfile_cmd))
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
