package integration_tests

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Brialius/antibruteforce/internal/grpc/api"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/Pallinder/go-randomdata"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

type apiStruct struct {
	apiCli                      api.AntiBruteForceServiceClient
	checkAuthRequest            *api.CheckAuthRequest
	checkAuthResponse           *api.CheckAuthResponse
	addToBlackListRequest       *api.AddToBlackListRequest
	addToBlackListResponse      *api.AddToBlackListResponse
	addToWhiteListRequest       *api.AddToWhiteListRequest
	addToWhiteListResponse      *api.AddToWhiteListResponse
	deleteFromBlackListRequest  *api.DeleteFromBlackListRequest
	deleteFromBlackListResponse *api.DeleteFromBlackListResponse
	deleteFromWhiteListRequest  *api.DeleteFromWhiteListRequest
	deleteFromWhiteListResponse *api.DeleteFromWhiteListResponse
	lastCheck                   bool
	lastError                   string
	ip                          string
	login                       string
	password                    string
	resetAt                     int
	passedRequests              int
	timeInterval                time.Duration
}

const (
	serviceUrl   = "antibruteforce-service:8080"
	randomString = "random"
)

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

var ctx = context.Background()

func (a *apiStruct) thereIsAServer(url string) error {
	conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(),
		grpc.WithUserAgent("antibruteforce integration_tests"))
	if err != nil {
		return err
	}
	a.apiCli = api.NewAntiBruteForceServiceClient(conn)
	if a.apiCli == nil {
		return err
	}
	return nil
}

func (a *apiStruct) setAddress(ip string) error {
	a.ip = ip
	return nil
}

func (a *apiStruct) checkAddress() error {
	a.login = randomString
	a.password = randomString
	return a.checkAuth()
}

func (a *apiStruct) checkAuth() error {
	l, p, i := a.login, a.password, a.ip
	if l == randomString {
		l = randomdata.Email()
	}
	if p == randomString {
		p = randomdata.Alphanumeric(10)
	}
	if i == randomString {
		i = randomdata.IpV4Address()
	}
	res, err := a.apiCli.CheckAuth(ctx, &api.CheckAuthRequest{Auth: &api.Auth{
		Login:    l,
		Password: p,
		Ip:       i,
	}})
	if err != nil {
		return err
	}
	a.lastCheck = res.GetOk()
	return nil
}

func (a *apiStruct) requestIsNotBlocked() error {
	if !a.lastCheck {
		return errors.New("request is blocked")
	}
	return nil
}

func (a *apiStruct) requestIsBlocked() error {
	if a.lastCheck {
		return errors.New("request is not blocked")
	}
	return nil
}

func (a *apiStruct) addAddressToWhitelist() error {
	res, err := a.apiCli.AddToWhiteList(ctx, &api.AddToWhiteListRequest{
		Net: a.ip,
	})
	if err != nil {
		return err
	}
	a.lastError = res.GetError()
	return nil
}

func (a *apiStruct) expectNoError() error {
	return nil
}

func (a *apiStruct) addAddressToBlacklist() error {
	res, err := a.apiCli.AddToBlackList(ctx, &api.AddToBlackListRequest{
		Net: a.ip,
	})
	if err != nil {
		return err
	}
	a.lastError = res.GetError()
	return nil
}

func (a *apiStruct) deleteAddressFromWhitelist() error {
	res, err := a.apiCli.DeleteFromWhiteList(ctx, &api.DeleteFromWhiteListRequest{
		Net: a.ip,
	})
	if err != nil {
		return err
	}
	a.lastError = res.GetError()
	return nil
}

func (a *apiStruct) deleteAddressFromBlacklist() error {
	res, err := a.apiCli.DeleteFromBlackList(ctx, &api.DeleteFromBlackListRequest{
		Net: a.ip,
	})
	if err != nil {
		return err
	}
	a.lastError = res.GetError()
	return nil
}

func (a *apiStruct) expectAnError() error {
	if a.lastError == "" {
		return errors.New("erros is expected")
	}
	return nil
}

func (a *apiStruct) setLogin(login string) error {
	a.login = login
	return nil
}

func (a *apiStruct) setIP(ip string) error {
	a.ip = ip
	return nil
}

func (a *apiStruct) setPassword(password string) error {
	a.password = password
	return nil
}

func (a *apiStruct) resetLimit() error {
	res, err := a.apiCli.ResetLimit(ctx, &api.ResetLimitRequest{
		Login: a.login,
		Ip:    a.ip,
	})
	if err != nil {
		return err
	}
	if res.GetError() != "" {
		return errors.New(res.GetError())
	}
	return nil
}

func (a *apiStruct) sendRequests(requests int) error {
	var passed int
	for i := 1; i <= requests; i++ {
		if i == a.resetAt {
			if err := a.resetLimit(); err != nil {
				return err
			}
		}

		err := a.checkAuth()
		if err != nil {
			return err
		}
		if a.lastCheck {
			passed++
		}
		time.Sleep(a.timeInterval)
	}
	a.passedRequests = passed
	return nil
}

func (a *apiStruct) requestArePassed(wantPassed int) error {
	if wantPassed != a.passedRequests {
		return fmt.Errorf("want %d passed requests, but got: %d", wantPassed, a.passedRequests)
	}
	return nil
}

func (a *apiStruct) timeBetweenRequestsIs(interval string) error {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return err
	}
	a.timeInterval = d
	return nil
}

func (a *apiStruct) resetAtRequests(resetAt int) error {
	a.resetAt = resetAt
	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &apiStruct{}
	s.BeforeScenario(func(interface{}) {
		a.resetAt = 0
		_ = a.thereIsAServer(serviceUrl)
	})
	s.Step(`^there is a server "([^"]*)"$`, a.thereIsAServer)
	s.Step(`^address "([^"]*)"$`, a.setAddress)
	s.Step(`^check address$`, a.checkAddress)
	s.Step(`^request is not blocked$`, a.requestIsNotBlocked)
	s.Step(`^request is blocked$`, a.requestIsBlocked)
	s.Step(`^add address to whitelist$`, a.addAddressToWhitelist)
	s.Step(`^expect no error$`, a.expectNoError)
	s.Step(`^add address to blacklist$`, a.addAddressToBlacklist)
	s.Step(`^delete address from whitelist$`, a.deleteAddressFromWhitelist)
	s.Step(`^delete address from blacklist$`, a.deleteAddressFromBlacklist)
	s.Step(`^expect an error$`, a.expectAnError)
	s.Step(`^login "([^"]*)"$`, a.setLogin)
	s.Step(`^IP "([^"]*)"$`, a.setIP)
	s.Step(`^password "([^"]*)"$`, a.setPassword)
	s.Step(`^send (\d+) requests$`, a.sendRequests)
	s.Step(`^time between requests is "([^"]*)"$`, a.timeBetweenRequestsIs)
	s.Step(`^(\d+) requests are passed$`, a.requestArePassed)
	s.Step(`^reset at (\d+) requests$`, a.resetAtRequests)
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
