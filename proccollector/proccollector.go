package proccollector

import (
	"fmt"
	"game_exporter/config"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// GameMetrics 指标结构体
type GameMetrics struct {
	processmetrics *prometheus.Desc
	mutex          sync.Mutex
}

//  newMetric
// 创建指标描述符
func newGlobalMetric(metricName string, docString string, procName []string) *prometheus.Desc {
	log.Printf("调用newGlobalMetric")
	return prometheus.NewDesc(metricName, docString, procName, nil)
}

// NewMetrics 初始化GameMetrics 结构体
func NewMetrics() *GameMetrics {
	log.Printf("调用NewMetric")
	return &GameMetrics{
		processmetrics: newGlobalMetric("game_procs_num", "number of process", []string{"procname"}),
	}
}

// Describe 传递结构体中的指标描述符到channel
func (c *GameMetrics) Describe(ch chan<- *prometheus.Desc) {
	log.Printf("调用Describe")
	ch <- c.processmetrics
}

// Collect 传递数据给channel
func (c *GameMetrics) Collect(ch chan<- prometheus.Metric) {
	log.Printf("调用Collect")
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for procName, procNum := range c.GrabProcessNum() {
		ch <- prometheus.MustNewConstMetric(c.processmetrics, prometheus.GaugeValue, float64(procNum), procName)
	}
}

// GrabProcessNum 实际执行抓取的方法
func (c *GameMetrics) GrabProcessNum() (processNumData map[string]int) {
	log.Printf("调用GrabProcessNum")
	var cmd string
	processNumData = make(map[string]int)
	configStruct, err := config.GetConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range configStruct.Processnames {
		log.Println(v.Cmdline)

		if len(v.Cmdline) == 1 {
			cmd = `ps aux | awk '/` + v.Cmdline[0] + `/ && !/awk/ '|wc -l`
		} else {
			newcmdline := modifyString(v.Cmdline)
			cmd = `ps aux | awk '/` + newcmdline[0] + `/ && /` + newcmdline[1] + `/  && !/awk/ '|wc -l`
		}
		result, err := exec.Command("/bin/bash", "-c", cmd).Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		pronum, _ := strconv.Atoi(strings.TrimSuffix(string(result), "\n"))
		processNumData[v.Name] = pronum

	}
	return

}

// 将结构体内cmdline中，涉及到“/”全部添加转义符 “\”
func modifyString(s []string) []string {
	log.Printf("调用了modifyString")
	// 只限于当cmdline有两个元素的时候，才去替换
	if len(s) == 2 {
		for i := 0; i < len(s); i++ {
			s[i] = strings.Replace(s[i], "/", "\\/", -1)
		}
	}
	return s
}
