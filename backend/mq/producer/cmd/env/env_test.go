package env

import (
	"fmt"
	"testing"
)

func TestCheckHttpPort(t *testing.T) {
	testData := "9000"
	defaultData := CheckHttpPort(testData)
	if defaultData != testData {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}

	testData = ""
	defaultData = CheckHttpPort(testData)
	if defaultData != "8080" {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}
}

func TestCheckUsername(t *testing.T) {
	testData := "user"
	defaultData := CheckUsername(testData)
	if defaultData != testData {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}

	testData = ""
	defaultData = CheckUsername(testData)
	if defaultData != "guest" {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}
}

func TestCheckPassword(t *testing.T) {
	testData := "12345"
	defaultData := CheckPassword(testData)
	if defaultData != testData {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}

	testData = ""
	defaultData = CheckPassword(testData)
	if defaultData != "guest" {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}
}

func TestCheckMqHost(t *testing.T) {
	testData := "rabbitmq.local.domain"
	defaultData := CheckMqHost(testData)
	if defaultData != testData {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}

	testData = ""
	defaultData = CheckMqHost(testData)
	if defaultData != "127.0.0.1" {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}
}

func TestCheckPort(t *testing.T) {
	testData := "5555"
	defaultData := CheckPort(testData)
	if defaultData != testData {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}

	testData = ""
	defaultData = CheckPort(testData)
	if defaultData != "5672" {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}
}

func TestCheckVhost(t *testing.T) {
	testData := "/test"
	defaultData := CheckVhost(testData)
	if defaultData != testData {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}

	testData = ""
	defaultData = CheckVhost(testData)
	if defaultData != "/" {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testData, defaultData)
		t.Error(message)
	}
}

func TestCheckQueues(t *testing.T) {
	testData := "test1,test2"
	testArray := []string{"test1", "test2"}
	defaultData := CheckQueues(testData)
	if len(testArray) != len(defaultData) {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testArray, defaultData)
		t.Error(message)
	}

	testData = ""
	defaultData = CheckQueues(testData)
	if len(defaultData) != 3 {
		message := fmt.Sprintf("Testdata: %s and function return value: %s not match", testArray, defaultData)
		t.Error(message)
	}
}

func TestCheckEnvs(t *testing.T) {
	CheckEnvs()
	message := fmt.Sprintf(" Default http port: %s\n Default dusername: %s\n Default password: %s\n Default messagequeue host: %s\n Default port: %s \n Default vhost: %s \n Default queues: %s \n", HttpPort, Username, Password, MqHost, Port, Vhost, Queues)
	fmt.Println(message)
}
