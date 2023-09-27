Feature: speculative_execution

  Scenario: speculative_execution
    Given that a deploy is executed against a node using the speculative_exec RPC API
    Then a valid speculative_exec_result will be returned
    And the speculative_exec has an api_version of "1.0.0"
    And the speculative_exec has a valid block_hash
    And the speculative_exec has a valid execution_results
    And the execution_results contains a cost of 123456
