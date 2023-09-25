Feature: query_balance RCP method against nctl

  Scenario: query_balance by main_purse_under_public_key
    Given that a query balance is obtained by main purse public key
    Then a valid query_balance_result is returned
    And the query_balance_result has an API version of "1.0.0"
    And the query_balance_result has a valid balance
  
  Scenario: query_balance by main_purse_under_account_hash 
    Given that a query balance is obtained by main purse account hash
    Then a valid query_balance_result is returned
    And the query_balance_result has an API version of "1.0.0"
    And the query_balance_result has a valid balance
  
  Scenario: query_balance by purse_uref
    Given that a query balance is obtained by main purse uref
    Then a valid query_balance_result is returned
    And the query_balance_result has an API version of "1.0.0"
    And the query_balance_result has a valid balance

