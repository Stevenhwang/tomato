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

// 输出选中的主机列表
func ListHosts(group string) map[string]interface{} {
	// 获取所有主机和主机组名
	list := map[string]interface{}{}
	keys := []string{}
	allhost := hosts.Hosts.GetStringMap("all.hosts")
	for k, v := range allhost {
		list[k] = v
	}
	childrenhost := hosts.Hosts.GetStringMap("all.children")
	for k, v := range childrenhost {
		keys = append(keys, k)
		val := v.(map[string]interface{})
		value := val["hosts"].(map[string]interface{})
		for x, y := range value {
			list[x] = y
		}
	}
	if group == "all" {
		return list
	} else {
		if FindValInSlice(keys, group) {
			key := fmt.Sprintf("all.children.%s.hosts", group)
			return hosts.Hosts.GetStringMap(key)
		} else {
			value, ok := list[group]
			if !ok {
				return nil
			} else {
				return map[string]interface{}{group: value}
			}
		}
	}
}
