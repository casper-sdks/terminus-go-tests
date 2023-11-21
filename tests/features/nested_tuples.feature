Feature: Nested Tuples

  Scenario: Nested Tuples
    Given that a nested Tuple1 is defined as ((1)) using U32 numeric values
    Then the "first" element of the Tuple1 is "(1)"
    And the Tuple1 bytes are "01000000"

    Given that a nested Tuple2 is defined as (1, (2, (3, 4))) using U32 numeric values
    Then the "first" element of the Tuple2 is 1
    And the "second" element of the Tuple2 is "(2, (3, 4))"
    And the Tuple2 bytes are "01000000020000000300000004000000"

    Given that a nested Tuple3 is defined as (1, 2, (3, 4, (5, 6, 7))) using U32 numeric values
    Then the "first" element of the Tuple3 is 1
    Then the "second" element of the Tuple3 is 2
    And the "third" element of the Tuple3 is "(3, 4, (5, 6, 7))"
    #And the Tuple3 bytes are "01000000020000000300000004000000050000000600000007000000"

    Given that the nested tuples are deployed in a transfer
    And the transfer is successful
    When the tuples deploy is obtained from the node
    Then the "first" element of the Tuple1 is "(1)"
    Then the "first" element of the Tuple2 is 1
    And the "second" element of the Tuple2 is "(2, (3, 4))"
    Then the "second" element of the Tuple3 is 2
    And the Tuple1 bytes are "01000000"
    And the Tuple2 bytes are "01000000020000000300000004000000"
    And the Tuple3 bytes are "01000000020000000300000004000000050000000600000007000000"
