package utils

import (
	"fmt"
	"os"
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
	// list := []string{}
	// files, _ := ioutil.ReadDir("./modules")
	// for _, f := range files {
	// 	list = append(list, strings.ReplaceAll(f.Name(), ".go", ""))
	// }
	// return list
	return []string{"ping", "shell", "copy"}
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

// 填充主机列表的认证参数
func FillParams(groups map[string]interface{}) map[string]map[string]interface{} {
	result := map[string]map[string]interface{}{}
	defaultuser := hosts.Hosts.GetString("default.user")
	defaultkeyFile := hosts.Hosts.GetString("default.key")
	defaultport := hosts.Hosts.GetInt("default.port")
	for h, v := range groups {
		result[h] = map[string]interface{}{"user": defaultuser, "key": defaultkeyFile, "port": defaultport, "password": ""}
		if v != nil {
			values := v.(map[string]interface{})
			if u, ok := values["user"]; ok {
				result[h]["user"] = u.(string)
			}
			if p, ok := values["port"]; ok {
				result[h]["port"] = p.(int)
			}
			// 如果自定参数有key，那么就确认key认证
			if k, ok := values["key"]; ok {
				result[h]["key"] = k.(string)
				continue
			}
			// 如果自定参数没有key，但是有password，那么确认password认证
			if pass, ok := values["password"]; ok {
				result[h]["key"] = ""
				result[h]["password"] = pass.(string)
			}
		}
	}
	return result
}

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDir 判断所给路径是否是文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
