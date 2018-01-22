/*
 *本文件时配置文件相关工具类
 *第一次使用本工具，必须先调用ConfSetPath函数
 */
package tool

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

//配置信息
var config map[string]interface{}

//配置文件路径
var confPath string

//设置配置文件路径，若path为""，则设置默认配置文件路径
func ConfSetPath(path string) error {
	confPath = path

	if err := ConfParse(); err != nil {
		return err
	}

	return nil
}

//获取当前配置文件路径
func ConfGetPath() string {
	return confPath
}

//解析配置文件，若调用了ConfSetPath会自动调用ConfParse
func ConfParse() error {
	//判断文件是否存在
	if err := ExistFile(confPath); err != nil {
		return err
	}

	f, err := os.Open(confPath)
	if err != nil {
		return err
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if !json.Valid(bs) {
		return errors.New("配置文件内容格式有误")
	}

	if err := json.Unmarshal(bs, &config); err != nil {
		return err
	}

	return nil
}

//获取key对应的配置文件值，请注意，本函数并未处理key为空的情况
func Conf(key string) string {
	v, ok := config[key].(string)
	if ok {
		return v
	}
	return ""
}

func ConfStrSlic(key string) []string {
	var str []string

	val, ok := config[key].([]interface{})
	if !ok {
		return nil
	}

	for _, v := range val {
		rv, ok := v.(string)
		if !ok {
			return nil
		}
		str = append(str, rv)
	}

	return str
}
