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

//go:generate easyjson
//easyjson:json
type User struct {
	Email string
}

// DomainStat - статистика по доменам
type DomainStat map[string]int

// emailRE - регулярка для парсинга email
var emailRE *regexp.Regexp

func init() {
	emailRE = regexp.MustCompile(`^[0-9a-z_\-\.]+@([0-9aa-z\.]+\.([a-z]+))$`)
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	s, err := calcStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("input data error: %w", err)
	}
	return s, nil
}

// calcStat поточно читает данные и поточно накапливает статистику
func calcStat(r io.Reader, suffix string) (DomainStat, error) {
	var user User
	stat := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		bytes := scanner.Bytes()

		// вытаскиваем email
		if err := easyjson.Unmarshal(bytes, &user); err != nil {
			return DomainStat{}, err
		}

		// если не оканчивается на нужный суффикс, дальше даже не смотрим (оптимизация)
		if !strings.HasSuffix(user.Email, suffix) {
			continue
		}

		// парсим email
		email := strings.ToLower(user.Email)
		m := emailRE.FindStringSubmatch(email)
		if m == nil {
			return DomainStat{}, errors.New("bad email")
		}
		domain, tld := m[1], m[2]

		// нам интересны только заданные tld
		if tld != suffix {
			continue
		}

		// накапливаем статистику по доменам
		stat[domain]++
	}

	return stat, nil
}
