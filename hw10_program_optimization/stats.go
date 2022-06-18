package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

var emailRE *regexp.Regexp

func init() {
	emailRE = regexp.MustCompile(`^[0-9a-z_\-\.]+@([0-9aa-z\.]+\.([a-z]+))$`)
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	s, err := calcStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return s, nil
}

func calcStat(r io.Reader, suffix string) (DomainStat, error) {
	var user User
	stat := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		//line := scanner.Text()
		bytes := scanner.Bytes()

		if err := easyjson.Unmarshal(bytes, &user); err != nil {
			//if err := easyjson.Unmarshal([]byte(line), &user); err != nil {
			return DomainStat{}, err
		}
		if !strings.HasSuffix(user.Email, suffix) {
			continue
		}
		email := strings.ToLower(user.Email)
		m := emailRE.FindStringSubmatch(email)
		if m == nil {
			return DomainStat{}, errors.New("bad email")
		}
		domain, tld := m[1], m[2]

		if tld != suffix {
			continue
		}

		// ЧТО ЗА ДИЧЬ!!!???
		stat[domain] += 1 // Получаем огромный расход RAM (тест по расходу памяти падает) memory used: 39Mb / 30Mb
		//stat[domain] += 2 // Все ОК по RAM (тест по расходу памяти проходит, но конечно падает по статистике)
		// P.S. личится заменой scanner.Text() на scanner.Bytes()

	}

	scanner.Bytes()

	return stat, nil
}
