Feature: Test antibruteforce service API
	Check API methods

	# Test white/black lists

	Scenario: Check not listed IP
		Given address "11.11.11.11"
		When check address
		Then request is not blocked

	Scenario: Add address to whitelist
		Given address "10.0.0.0/8"
		When add address to whitelist

	Scenario: Add another address to whitelist
		Given address "192.168.11.11/32"
		When add address to whitelist

	Scenario: Add address to blacklist
		Given address "192.168.0.0/16"
		When add address to blacklist

	Scenario: Add another address to blacklist
		Given address "10.11.11.11/32"
		When add address to blacklist

	Scenario: Check whitelisted IP
		Given address "10.10.10.10"
		When check address
		Then request is not blocked

	Scenario: Check blacklisted IP
		Given address "192.168.1.1"
		When check address
		Then request is blocked

	Scenario: Check whitelisted and blacklisted IP
		Given address "10.11.11.11"
		When check address
		Then request is not blocked

	Scenario: Check blacklisted and whitelisted IP
		Given address "192.168.11.11"
		When check address
		Then request is not blocked

	Scenario: Delete address from whitelist
		Given address "10.0.0.0/8"
		When delete address from whitelist

	Scenario: Delete another address from whitelist
		Given address "192.168.11.11/32"
		When delete address from whitelist

	Scenario: Delete address from whitelist with error
		Given address "10.0.0.0/8"
		When delete address from whitelist

	Scenario: Delete another address from whitelist with error
		Given address "192.168.11.11/32"
		When delete address from whitelist

	Scenario: Delete address from blacklist
		Given address "192.168.0.0/16"
		When delete address from blacklist

	Scenario: Delete another address from blacklist
		Given address "10.11.11.11/32"
		When delete address from blacklist

	Scenario: Delete address from blacklist with error
		Given address "192.168.0.0/16"
		When delete address from blacklist

	Scenario: Delete another address from blacklist with error
		Given address "10.11.11.11/32"
		When delete address from blacklist

	Scenario: Check former blacklisted IP
		Given address "10.11.11.11"
		When check address
		Then request is not blocked
