package main

import (
	"go_client/ClientAndServer"
	"go_client/DeviceAndServer"
)

func main() {
	//go as client
	ClientAndServer.ClientConnetToServer()

	//go as server
	DeviceAndServer.ListenFromDevice()
}

/*
package main

import (
	"fmt"
)

func main() {
	//创建map
	countryCapitalMap := make(map[string]string)

	//插入每个国家对应的首都
	countryCapitalMap["France"] = "Paris"
	countryCapitalMap["Italy"] = "Rome"
	countryCapitalMap["Japan"] = "Tokyo"
	countryCapitalMap["India"] = "New Delhi"

	//使用key输出map的值
	fmt.Println("第一种输出的方式")
	for country := range countryCapitalMap {
		fmt.Println("Capital of", country, "is", countryCapitalMap[country])
	}

	//直接输出map的key和值
	fmt.Println("第二种输出的方式")
	for k, v := range countryCapitalMap {
		fmt.Println("Capital of", k, "is", v)
	}

	//查看元素是否在map中,变量ok会返回true或者false，当查询的key在map中则返回true并且captial会获取到map中的值
	fmt.Println("***************************************")
	captial, ok := countryCapitalMap["United States"]
	if ok {
		fmt.Println("Capital of", captial, "is", countryCapitalMap[captial])
	} else {
		fmt.Println("Capital of United States is not present")
	}
}
*/
