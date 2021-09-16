package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
	"tomato/hosts"
)

// FindValInSlice 查询val是否在 string slice 中
func FindValInSlice(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// 输出模块列表
func ListModules() []string {
	list := []string{}
	files, _ := ioutil.ReadDir("./modules")
	for _, f := range files {
		list = append(list, strings.ReplaceAll(f.Name(), ".go", ""))
	}
	return list
}

// 输出主机列表
func ListHosts(group string) map[string]interface{} {
	if group == "all" {
		list := map[string]interface{}{}
		allhost := hosts.Hosts.GetStringMap("all.hosts")
		for k, v := range allhost {
			list[k] = v
		}
		childrenhost := hosts.Hosts.GetStringMap("all.children")
		for _, v := range childrenhost {
			val := v.(map[string]interface{})
			value := val["hosts"].(map[string]interface{})
			for x, y := range value {
				list[x] = y
			}
		}
		return list
	} else {
		key := fmt.Sprintf("all.children.%s.hosts", group)
		return hosts.Hosts.GetStringMap(key)
	}
}
