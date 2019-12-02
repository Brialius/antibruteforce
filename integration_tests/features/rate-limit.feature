Feature: Test antibruteforce service API
	Check API methods

	# Test rate limits

	Scenario: Check rate for login
		Given login "test-login"
		And IP "random"
		And password "random"
		And time between requests is "0s"
		When send 50 requests
		Then 10 requests are passed

	Scenario: Check rate for password
		Given login "random"
		And IP "random"
		And password "test-login-password"
		And time between requests is "0s"
		When send 100 requests
		Then 50 requests are passed

	Scenario: Check rate for IP
		Given login "random"
		And IP "120.110.100.10"
		And password "random"
		And time between requests is "0s"
		When send 120 requests
		Then 100 requests are passed

	Scenario: Check rate for login and IP with reset
		Given login "test-login-reset"
		And IP "123.214.41.52"
		And password "random"
		And reset at 6 requests
		And time between requests is "0s"
		When send 20 requests
		Then 15 requests are passed

	Scenario: Check rate for password with slow rate
		Given login "random"
		And IP "random"
		And password "test-login-password-slow"
		And time between requests is "0.7s"
		When send 60 requests
		Then 60 requests are passed
