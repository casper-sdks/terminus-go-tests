Feature: speculative_execution

  Scenario: speculative_execution
    Given that the "faucet" account transfers 2500000000 to user-1 account with a gas payment amount of 100000000 using the speculative_exec RPC API
    Then the speculative_exec has an api_version of "1.0.0"
    And a valid speculative_exec_result will be returned with 55 transforms
    And the speculative_exec has a valid block_hash
    And the execution_results contains a cost of 100000000
    And the speculative_exec has a valid execution_result
    And the speculative_exec execution_result transform wth the transfer key contains the deploy_hash
    And the speculative_exec execution_result transform with the transfer key has the amount of 2500000000
    And the speculative_exec execution_result transform with the transfer key has the "from" field set to the "faucet" account hash
    And the speculative_exec execution_result transform with the transfer key has the "to" field set to the "user-1" account hash
    And the speculative_exec execution_result transform with the transfer key has the "source" field set to the purse uref of the "faucet" account
    And the speculative_exec execution_result transform with the transfer key has the "target" field set to the purse uref of the "user-1" account
    And the speculative_exec execution_result transform with the deploy key has the deploy_hash of the transfer's hash
    And the speculative_exec execution_result transform with a deploy key has a gas field of 100000000
    And the speculative_exec execution_result transform with a deploy key has 1 transfer with a valid transfer hash
    And the speculative_exec execution_result transform with a deploy key has as from field of the "faucet" account hash
    And the speculative_exec execution_result transform with a deploy key has as source field of the "faucet" account purse uref
    And the speculative_exec execution_result contains at least 6 valid balance transforms
    And the speculative_exec execution_result 1st balance transform is an Identity transform
    And the speculative_exec execution_result last balance transform is an Identity transform is as WriteCLValue of type "U512"
    And the speculative_exec execution_result contains a valid AddUInt512 transform with a value of 100000000
