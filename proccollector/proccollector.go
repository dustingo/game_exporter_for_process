package proccollector

import (
	"bilibili/learn909/03process.yml/config"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// var wg sync.WaitGroup
// var rwlock sync.RWMutex

// GameMetrics 指标结构体
type GameMetrics struct {
	processmetrics *prometheus.Desc
	mutex          sync.Mutex
}

//  newMetric
// 创建指标描述符
func newGlobalMetric(metricName string, docString string, procName []string) *prometheus.Desc {
	return prometheus.NewDesc(metricName, docString, procName, nil)
}

// NewMetrics 初始化GameMetrics 结构体
func NewMetrics() *GameMetrics {
	return &GameMetrics{
		processmetrics: newGlobalMetric("game_procs_num", "number of process", []string{"procname"}),
	}
}

// Describe 传递结构体中的指标描述符到channel
func (c *GameMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.processmetrics
}

// Collect 传递数据给channel
func (c *GameMetrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for procName, procNum := range c.GrabProcessNum() {
		ch <- prometheus.MustNewConstMetric(c.processmetrics, prometheus.GaugeValue, float64(procNum), procName)
	}
}

// GrabProcessNum 实际执行抓取的方法
func (c *GameMetrics) GrabProcessNum() (processNumData map[string]int) {
	var cmd string
	processNumData = make(map[string]int)
	configPath := "./gameprocess.yml"
	configStruct := config.GetConfig(&configPath)
	for _, v := range configStruct.Processnames {

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
	// 只限于当cmdline有两个元素的时候，才去替换
	if len(s) == 2 {
		for i := 0; i < len(s); i++ {
			s[i] = strings.Replace(s[i], "/", "\\/", -1)
		}
	}
	return s
}
