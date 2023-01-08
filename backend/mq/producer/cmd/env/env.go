package env

import (
	"fmt"
	"os"
	"strings"
)

var (
	HttpPort  string
	Username  string
	Password  string
	MqHost    string
	Port      string
	Vhost     string
	Queues    []string
	tmpQueues string
)

// Get environment variables
func getUsername() string {
	return os.Getenv("USERNAME")
}

func getPassword() string {
	return os.Getenv("PASSWORD")
}

func getMqHost() string {
	return os.Getenv("MQHOST")
}

func getPort() string {
	return os.Getenv("PORT")
}

func getHttpPort() string {
	return os.Getenv("PORT")
}

func getVhost() string {
	return os.Getenv("VHOST")
}

func getQueues() string {
	return os.Getenv("QUEUE_NAMES")
}

// Set default value to variable if empty

func CheckUsername(username string) string {
	if len(username) != 0 {
		return username
	}
	username = "quest"
	return username
}

func CheckPassword(password string) string {
	if len(password) != 0 {
		return password
	}
	password = "guest"
	return password
}

func CheckPort(port string) string {
	if len(port) != 0 {
		return port
	}
	port = "5672"
	return port
}

func CheckHttpPort(port string) string {
	if len(port) != 0 {
		return port
	}
	port = "8080"
	return port
}

func CheckMqHost(host string) string {
	if len(host) != 0 {
		return host
	}
	host = "127.0.0.1"
	return host
}

func CheckVhost(vhost string) string {
	if len(vhost) != 0 {
		return vhost
	}
	vhost = "/"
	return vhost
}

func CheckQueues(queue string) []string {
	fmt.Println(queue)
	if len(queue) != 0 {
		return strings.Split(queue, ",")
	}
	return []string{"test1", "test2", "test3"}
}

// Set variables value from environment variables
func getValues() {
	Username = getUsername()
	Password = getPassword()
	MqHost = getMqHost()
	Port = getPort()
	Vhost = getVhost()
	tmpQueues = getQueues()
	HttpPort = getHttpPort()

}

// Override variables
func CheckEnvs() {
	getValues()
	Username = CheckUsername(Username)
	Password = CheckPassword(Password)
	MqHost = CheckMqHost(MqHost)
	Port = CheckPort(Port)
	Vhost = CheckVhost(Vhost)
	Queues = CheckQueues(tmpQueues)
	HttpPort = CheckHttpPort(HttpPort)
}
