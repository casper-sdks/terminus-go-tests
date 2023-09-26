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

  Scenario: query_balance by main_purse_under_public_key using a block hash state identifier
    When a transfer of 2500000000 is made to user-1's purse
    And that a query balance is obtained by user-1's main purse public and latest block identifier
    Then the balance includes the transferred amount
    When that a query balance is obtained by user-1's main purse public key and previous block identifier
    Then the balance is the pre transfer amount

  Scenario: query_balance by main_purse_under_public_key using a state root hash state identifier
    When a transfer of 2500000000 is made to user-1's purse
    And that a query balance is obtained by user-1's main purse public and latest state root hash identifier
    Then the balance includes the transferred amount
    When that a query balance is obtained by user-1's main purse public key and previous state root hash identifier
    Then the balance is the pre transfer amount




