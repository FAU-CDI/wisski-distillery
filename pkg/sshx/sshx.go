package sshx

import "github.com/gliderlabs/ssh"

// ParseAllKeys parses all keys from the list of bytes
func ParseAllKeys(bytes []byte) (keys []ssh.PublicKey) {
	var key ssh.PublicKey
	var err error
	for {
		key, _, _, bytes, err = ssh.ParseAuthorizedKey(bytes)
		if err != nil {
			break
		}
		keys = append(keys, key)
	}
	return
}
