package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func get_session_key() string {
	b := ""
	for i := 0; i < 10; i++ {
		b += strconv.Itoa(rand.Intn(9) + 1)
	}
	return b
}

func get_hash_str() string {
	li := ""
	for i := 0; i < 5; i++ {
		li += strconv.Itoa(rand.Intn(10))
	}
	return li
}

type Session_protector struct {
	__hash string
}

func (self Session_protector) __calc_hash(session_key string, val int) string {
	switch val {
	case 1:
		result := ""
		ret := ""
		for idx := 0; idx < 5; idx++ {
			result += string(session_key[idx])
		}
		i, _ := strconv.Atoi(result)
		result = "00" + strconv.Itoa(i%97)
		for idx := len(result) - 2; idx < len(result); idx++ {
			ret += string(result[idx])
		}
		return ret
	case 2:
		result := ""
		for idx := 0; idx < len(session_key); idx++ {
			result += string(session_key[len(session_key)-idx-1])
		}
		return result
	case 3:
		return session_key[len(session_key)-5:] + session_key[0:5]
	case 4:
		num := 0
		for i := 1; i < 9; i++ {
			num += int(session_key[i]) + 41
		}
		return string(num)
	case 5:
		ch := ""
		result := 0
		for idx := 0; idx < len(session_key); idx++ {
			ch = string(int(int(session_key[idx]) ^ 43))
			if _, err := strconv.Atoi(ch); err != nil {
				ch = string(int(ch[0]))
			}
			num, _ := strconv.Atoi(ch)
			result += num
		}
		return strconv.Itoa(result)
	default:
		result, _ := strconv.Atoi(session_key)
		return strconv.Itoa(result + val)
	}
}

func (self Session_protector) next_session_key(session_key string) string {
	result := 0
	if self.__hash == "" {
		fmt.Println("hash is empty")
		return get_session_key()
	}
	for idx := 0; idx < len(self.__hash); idx++ {
		i := string(self.__hash[idx])
		if _, err := strconv.Atoi(i); err != nil {
			fmt.Println("Here is letter")
			return get_session_key()
		}
	}
	ret := ""
	for idx := 0; idx < len(self.__hash); idx++ {
		num, _ := strconv.Atoi(string(self.__hash[idx]))
		k, _ := strconv.Atoi(self.__calc_hash(session_key, num))
		result += k
	}
	for idx := 0; idx < 10 && idx < len(strconv.Itoa(result)); idx++ {
		ret += string((strconv.Itoa(result))[idx])
	}
	m := ""
	ret = "0000000000" + ret
	for idx := len(ret) - 10; idx < len(ret); idx++ {
		m += string(ret[idx])
	}
	return m
}

func main() {
	fmt.Print("Server-IP:port ")
	IPAdr := ""
	fmt.Fscan(os.Stdin, &IPAdr)
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	conn, err := net.Dial("tcp", IPAdr)
	if err != nil {
		fmt.Println("Server not found...")
	} else {
		cl_hash_string := get_hash_str()
		key1 := get_session_key()
		fmt.Print(cl_hash_string + "\n")
		fmt.Fprintf(conn, cl_hash_string+key1+"\n")
		client_protector := Session_protector{cl_hash_string}
		key2, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Server is not responding.")
		}
		key1 = client_protector.next_session_key(key1)
		for {
			text := ""
			// send to socket
			fmt.Fprintf(conn, strings.Replace(text, "\n", "", -1)+key1+"\n")
			// listen for reply
			fmt.Println("Waiting for answer...")
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Server is not responding.")
			}
			key2 = ""
			for i := len(message) - 11; i < len(message)-1; i++ {
				key2 += string(message[i])
			}
			for i := 0; i < len(message)-11; i++ {
				text += string(message[i])
			}
			key1 = client_protector.next_session_key(key1)
			fmt.Println("key: ", key1, " ", key2)
			key1 = client_protector.next_session_key(key1)
		}
	}
}
