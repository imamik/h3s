package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"h3s/internal/utils/logger"
	"log"
	"time"
)

func dialWithRetries(ip string, sshConfig *ssh.ClientConfig, retryInterval time.Duration, maxRetries int) (*ssh.Client, error) {
	for i := 0; i < maxRetries; i++ {
		c, err := ssh.Dial("tcp", ip+":22", sshConfig)
		if err == nil && c != nil {
			return c, err
		}
		logger.LogResourceEvent(logger.Server, "SSH", ip, logger.Failure, err)
		retryBackoff := time.Duration(i+1) * retryInterval
		log.Printf("Failed to dial: %s, retrying in %s", err, retryBackoff)
		time.Sleep(retryBackoff)
	}
	return nil, fmt.Errorf("failed to dial %s after %d retries", ip, maxRetries)
}
