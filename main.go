package main

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	//定义命令行参数方式1
	var network string
	var netmask string
	var netcard string
	var gateway string
	var debug_mode_temp string
	var debug_mode bool

	flag.StringVar(&network, "n", "x.x.x.x", "网络地址 | ip address")
	flag.StringVar(&netmask, "m", "255.255.255.0", "子网掩码 | network mask")
	flag.StringVar(&netcard, "c", "eth1/eth1/mgmt1", "网卡名 | network card ")
	flag.StringVar(&gateway, "g", "x.x.x.1", "网关 | gateway")
	flag.StringVar(&debug_mode_temp, "d", "n", "开启debug模式 y/s |start debug mode y/s")

	//解析命令行参数
	flag.Parse()
	fmt.Println(network,netmask,netcard,gateway,debug_mode_temp)
	if string(network) == "x.x.x.x" {
		panic("请指定 -n 必须输" + "入网络地址")
	}
	if string(gateway) == "x.x.x.1" {
		gateway = network[:strings.LastIndex(network, ".")] + ".1"
	}
	if string(netcard) == "eth1/eth1/mgmt1" {

		netcard_temp, err := find_netcard_name()
		if err == nil {
			netcard = netcard_temp
		} else {
			panic("请指定 -c 并输入网卡名称")
		}
	}
	if string(debug_mode_temp)!="y" {
		debug_mode = false
	}else {
		debug_mode = true
	}
	fmt.Println("修改预期目标:\n")
	fmt.Println("    网络地址:" + network)
	fmt.Println("    子网掩码:" + netmask)
	fmt.Println("    网卡名:" + netcard)
	fmt.Println("    网关:" + gateway + "\n")

	fmt.Println("1.开始修改网关")
	run_cmd("rm -rf /etc/sysconfig/network", debug_mode)
	run_cmd("touch /etc/sysconfig/network", debug_mode)
	run_cmd("echo GATEWAY="+gateway+">> /etc/sysconfig/network", debug_mode)

	fmt.Println("2.开始修改网络")
	card_path := "/etc/sysconfig/network-scripts/ifcfg-" + netcard
	run_cmd("mv "+card_path+" "+card_path+".bak", debug_mode)
	run_cmd("touch "+card_path, debug_mode)
	run_cmd("echo DEVICE="+netcard+" >> "+card_path, debug_mode)
	run_cmd("echo TYPE=Ethernet >> "+card_path, debug_mode)
	run_cmd("echo ONBOOT=yes >> "+card_path, debug_mode)
	run_cmd("echo IPV4_FAILURE_FATAL=yes >> "+card_path, debug_mode)
	run_cmd("echo IPV6INIT=yes >> "+card_path, debug_mode)
	run_cmd("echo IPADDR="+network+" >> "+card_path, debug_mode)
	run_cmd("echo NETMASK="+netmask+" >> "+card_path, debug_mode)
	fmt.Println("3.开始重启网卡")
	run_cmd("ifdown "+netcard, debug_mode)
	run_cmd("ifup "+netcard, debug_mode)
}
func find_netcard_name() (string, error) {
	cmd := exec.Command("ip", "a")
	out, _ := cmd.CombinedOutput()
	name_list := []string{"zeth1", "eth0", "eth1", "mgmt1", "enp1s0"}
	for _, v := range name_list {
		if strings.Contains(string(out), v) {
			return v, nil
		}
	}
	//开始
	return "", errors.New("未自动捕获到网卡名，请手动输入网卡。")
}
func run_cmd(cmd string, mode bool) {
	if mode == true {
		fmt.Println("    执行命令:" + cmd)
	}
	err := exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		fmt.Println(err)
	}
}
