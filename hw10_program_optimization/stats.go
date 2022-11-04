package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsonIter "github.com/json-iterator/go"
)

var json = jsonIter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	buf := bufio.NewScanner(r)
	buf.Split(bufio.ScanLines)

	i := 0
	for buf.Scan() {
		var user User
		if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
			return
		}
		result[i] = user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		index := strings.IndexAny(user.Email, "@") + 1

		if strings.HasSuffix(user.Email, "."+domain) {
			str := strings.ToLower(user.Email[index:])

			result[str]++
		}
	}

	return result, nil
}
