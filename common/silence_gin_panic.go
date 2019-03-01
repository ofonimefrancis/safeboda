package common

import (
	"net"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ofonimefrancis/safeboda/common/log"
)

func SilenceSomePanics() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if isBrokenPipeError(r) {
				log.Warningf("caught a broken pipe, it shouldn't be an exception, so silencing it")
				return
			}
			panic(r)
		}()
		c.Next()
	}
}

func isBrokenPipeError(r interface{}) bool {
	if r == nil {
		return false
	}
	if nerr, ok := r.(*net.OpError); ok && nerr.Op == "write" {
		if sysErr, ok := nerr.Err.(*os.SyscallError); ok && sysErr.Syscall == "write" {
			if e, ok := sysErr.Err.(syscall.Errno); ok && e.Error() == "broken pipe" {
				return true
			}
		}
	}

	return false
}
