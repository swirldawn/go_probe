package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	//"github.com/tidwall/gjson"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	re := getLocal()
	fmt.Fprintln(w, re)
}

//获取本服务器的
func getLocal() string {
	disk := getDiskStatus()

	cpu_mem := getCPUstatus()

	version := getLinuxVersion()
	re := disk + "|" + cpu_mem + "|" + version
	return re
}

//获取其他服务器的
func getAll() []string {

}

func main() {
	http.HandleFunc("/get_local", IndexHandler)
	http.ListenAndServe("127.0.0.1:8084", nil)

}

//获取系统发行版本
func getLinuxVersion() string {
	str := ""
	if runtime.GOOS == "darwin" {
		str = `CentOS release 6.5 (Final)
Kernel \r on an \m

		`
	} else {
		str, _ = exec_shell("cat /etc/issue")
	}
	s := strings.Split(str, "\n")
	return s[0]
}

//获取cpu使用率 改命令只能在linux下执行
func getCPUstatus() string {
	str := ""
	if runtime.GOOS == "darwin" {
		str = `top - 05:00:03 up 200 days,  3:08,  1 user,  load average: 0.00, 0.00, 0.00
Tasks: 123 total,   2 running, 121 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.4 us,  0.5 sy,  0.0 ni, 99.1 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
KiB Mem :   516496 total,    62732 free,    57264 used,   396500 buff/cache
KiB Swap:   135164 total,   119964 free,    15200 used.   415156 avail Mem

 PID USER      PR  NI    VIRT    RES    SHR S %CPU %MEM     TIME+ COMMAND
11415 nobody    20   0  803940   5556   2488 R  6.7  1.1   2490:33 ./mtproto-proxy -u nobody -p 8888 -H 9538 -S 685a0b340aaf588de2aeee9447017a4d --aes-pwd proxy-secret proxy-+
25433 root      20   0   40388   3584   3068 R  6.7  0.7   0:00.01 top -bn 1 -i -c
`
	} else {
		str, _ = exec_shell("top -bn 1 -i -c")
	}

	s := strings.Split(str, "\n")
	mem := ""
	cpu := ""
	for value, _ := range s {
		if strings.Index(s[value]+" ", "Mem:") != -1 || strings.Index(s[value]+" ", "Mem :") != -1 {
			mem = s[value]
		}
		if strings.Index(s[value]+" ", "Cpu(s):") != -1 {
			cpu = s[value]
		}
	}
	mem = strings.Replace(mem, " ", "", -1)
	mem = strings.Replace(mem, "KiB", "", -1)
	mem = strings.Replace(mem, "Mem:", "", -1)

	cpu = strings.Replace(cpu, "Cpu(s):", "", -1)
	cpu = strings.Replace(cpu, "%", "", -1)
	cpu = strings.Replace(cpu, "us", "用户使用", -1)
	cpu = strings.Replace(cpu, "sy", "系统使用", -1)
	cpu = strings.Replace(cpu, "id", "空闲", -1)
	cpu = strings.Replace(cpu, " ", "", -1)

	re := "cpu:" + cpu + "|mem:" + mem

	return re
}

//获取磁盘占用量
func getDiskStatus() string {
	str, _ := exec_shell("df -lh")
	s := strings.Split(str, "\n")
	real_str := ""
	for value, _ := range s {
		if strings.Index(s[value]+" ", " / ") != -1 {
			real_str = s[value]
		}
	}
	real_str = strings.Replace(real_str, "  ", " ", -1)
	real_str_arr := strings.Split(real_str, " ")
	var new_arr []string = make([]string, 0)

	for index := range real_str_arr {
		if strings.Replace(real_str_arr[index], " ", "", -1) != "" {
			new_arr = append(new_arr, real_str_arr[index])
		}
	}
	// fmt.Println(real_str, new_real_str_arr, len(new_real_str_arr))
	re := "disk:磁盘总量" + new_arr[1] + ",已使用" + new_arr[2] + ",剩余" + new_arr[3] + ",使用率" + new_arr[4]
	return re
}

//获取系统名称
func getOsName() string {
	name, err := os.Hostname()
	if err == nil {
		return name
	}
	return ""
}

//获取系统核心数量
func getCoreNum() (num int) {
	num = runtime.GOMAXPROCS(0)
	return num
}

//执行系统命令
func exec_shell(s string) (string, error) {

	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	checkErr(err)

	return out.String(), err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
