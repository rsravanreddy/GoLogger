package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os/exec"
	"strings"
)

type LockScreenResouce struct{}

func Exe_cmd(cmd string) {
	fmt.Println(cmd)
	parts := strings.Fields(cmd)
	out, err := exec.Command(parts[0], parts[1]).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}

func (tr *LockScreenResouce) LockScreen(c *gin.Context) {

	Exe_cmd(`/System/Library/CoreServices/"Menu Extras"/User.menu/Contents/Resources/CGSession -suspend`)
	c.JSON(201, "success")
}
