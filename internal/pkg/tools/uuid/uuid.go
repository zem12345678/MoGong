package uuid

import (
	"fmt"
	"os/exec"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func NewUUID() (uuidStr string) {
	newUuid := uuid.NewV4()
	uuidStr = newUuid.String()
	return
}

func NewUpperUUID() (uuidStr string) {
	return strings.ToUpper(NewUUID())
}

func GetHostUuid() (uuid string, err error) {
	dmidecode := exec.Command("dmidecode", "-s", "system-uuid")
	bytes, err := dmidecode.Output()
	if err != nil {
		fmt.Println("get uuid command error: ", err)
		return "", nil
	}
	uuid = strings.Trim(string(bytes), "\n")
	return uuid, nil
}
