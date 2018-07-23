package BDYString

import (
	"strconv"
	"strings"
)

//去掉字符串中最后位
func DeleteTail(str string) string {
	lenght := len([]rune(str))
	buf := []byte(str)
	str_out := string(buf[0 : lenght-1])
	return str_out
}

//解析BDT01协议中的状态位
func ParseStatusData(str string) (signal string, sat_num string, bat string, mode string) {
	buf := []byte(str)
	signal = string(buf[0:3])
	sat_num = string(buf[3:6])
	bat = string(buf[6:9])
	mode = string(buf[10:12])
	return signal, sat_num, bat, mode
}

//从字符串的某位开始，获取指定长度的字符串
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//获取指定开始到指定结束的字符串
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

//10进制整形转16进制字符串
func Int2HexString(lenght int) string {
	var num [4]int
	var buf [4]string
	num[0] = lenght / 4096
	num[1] = lenght % 4096 / 256
	num[2] = lenght % 4096 % 256 / 16
	num[3] = lenght % 16
	for i := 0; i < 4; i++ {
		if num[i] == 10 {
			buf[i] = "a"
		} else if num[i] == 11 {
			buf[i] = "b"
		} else if num[i] == 12 {
			buf[i] = "c"
		} else if num[i] == 13 {
			buf[i] = "d"
		} else if num[i] == 14 {
			buf[i] = "e"
		} else if num[i] == 15 {
			buf[i] = "f"
		} else {
			buf[i] = strconv.Itoa(num[i])
		}
	}
	str_out := buf[0] + buf[1] + buf[2] + buf[3]
	return str_out
}

//16进制字符串转10进制整形
func HexString2Int(str string) int {
	lenght := len([]rune(str))
	var num [4]int

	if lenght == 3 {
		str = "0" + str
	} else if lenght == 2 {
		str = "00" + str
	} else if lenght == 1 {
		str = "000" + str
	} else {
	}
	buf := []byte(str)
	for i := 0; i < 4; i++ {
		if string(buf[i]) == "a" {
			num[i] = 10
		} else if string(buf[i]) == "b" {
			num[i] = 11
		} else if string(buf[i]) == "c" {
			num[i] = 12
		} else if string(buf[i]) == "d" {
			num[i] = 13
		} else if string(buf[i]) == "e" {
			num[i] = 14
		} else if string(buf[i]) == "f" {
			num[i] = 15
		} else {
			num[i], _ = strconv.Atoi(string(buf[i]))
		}
	}
	flag := num[0]*4096 + num[1]*256 + num[2]*16 + num[3]
	return flag
}
