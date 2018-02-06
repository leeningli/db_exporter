package leeconfig

import(
	"flag"
	"log"
	"config"
	"runtime"
)

var(
	configFile = flag.String("configfile", "config.ini", "General configuration file")
)

var TOPIC = make(map[string]string)

func GetConfig(appname string) (map[string]string){
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find",*configFile, err)
	}
	if cfg.HasSection(appname){
		section, err := cfg.SectionOptions(appname)
		if err == nil {
			for _, v := range section {
				options, err := cfg.String(appname, v)
				if err == nil{
					TOPIC[v] = options
				}
			}
		}
	}
	return TOPIC
}
