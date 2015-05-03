package gomysql

import (
	"crypto/sha1"
	"strconv"
	"strings"
)

func parseVersion(versionString string) (version []byte, err error) {
	parts := strings.Split(versionString, "-")
	for _, s := range strings.Split(parts[0], ".") {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		version = append(version, byte(v))
	}
	return version, nil
}

func passwordToken(password string, challange []byte) (token []byte) {
	d := sha1.New()

	d.Write([]byte(password))
	h1 := d.Sum(nil)

	d.Reset()
	d.Write(h1)
	h2 := d.Sum(nil)

	d.Reset()
	d.Write(challange)
	d.Write(h2)
	token = d.Sum(nil)

	for i := range token {
		token[i] ^= h1[i]
	}

	return token
}
