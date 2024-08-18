package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = "." + strings.ToLower(domain)
	stat := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}

		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, domain) {
			domainName := strings.SplitN(email, "@", 2)[1]
			stat[domainName]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return stat, nil
}
