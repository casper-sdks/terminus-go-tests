Feature: Nested Options

  Scenario: Nested Option
    Given that a nested Option has an inner type of Option with a type of String and a value of "77541f74-c0b6-428e-9f53-ddd298f1bc8c"
    Then the inner type is Option with a type of String and a value of "77541f74-c0b6-428e-9f53-ddd298f1bc8c"
    And the bytes are "01012400000037373534316637342d633062362d343238652d396635332d646464323938663162633863"
    Given that the nested Option is deployed in a transfer
    And the transfer containing the nested Option is successfully executed
    When the Option is read from the deploy
    Then the inner type is Option with a type of String and a value of "77541f74-c0b6-428e-9f53-ddd298f1bc8c"
    And the bytes are "01012400000037373534316637342d633062362d343238652d396635332d646464323938663162633863"

  Scenario: Nested List
    Given that a nested Option has an inner type of List with a type of U256 and a value of (1, 2, 3)
    And the nested list's length is 3
    And the nested list's "1st" item is a CLValue with U256 value of 1
    And the nested list's "2nd" item is a CLValue with U256 value of 2
    And the nested list's "3rd" item is a CLValue with U256 value of 3
    And the bytes are "0103000000010101020103"
    Given that the nested Option is deployed in a transfer
    And the transfer containing the nested Option is successfully executed
    When the Option is read from the deploy
    And the nested list's length is 3
    And the nested list's "1st" item is a CLValue with U256 value of 1
    And the nested list's "2nd" item is a CLValue with U256 value of 2
    And the nested list's "3rd" item is a CLValue with U256 value of 3
    And the bytes are "0103000000010101020103"

  Scenario: Nested Tuple
    Given that a nested Option has an inner type of Tuple2 with a type of "String,U32" values of ("7fc5ae1d-348e-4b01-97cf-a1c2942b28a9", 10000)
    Then the inner type is Tuple2 with a type of "String,U32" and a value of ("7fc5ae1d-348e-4b01-97cf-a1c2942b28a9", 10000)
    And the bytes are "012400000037666335616531642d333438652d346230312d393763662d61316332393432623238613910270000"
    Given that the nested Option is deployed in a transfer
    And the transfer containing the nested Option is successfully executed
    When the Option is read from the deploy
    Then the inner type is Tuple2 with a type of "String,U32" and a value of ("7fc5ae1d-348e-4b01-97cf-a1c2942b28a9", 10000)
    And the bytes are "012400000037666335616531642d333438652d346230312d393763662d61316332393432623238613910270000"

  Scenario: Nested Map
    Given that a nested Option has an inner type of Map with a type of "String,U32" value of {"ONE": 2}
    Then the inner type is Map with a type of "String,U32" and a value of {"ONE": 2}
    And the bytes are "0101000000030000004f4e4502000000"
    Given that the nested Option is deployed in a transfer
    And the transfer containing the nested Option is successfully executed
    When the Option is read from the deploy
    Then the inner type is Map with a type of "String,U32" and a value of {"ONE": 2}
    And the bytes are "0101000000030000004f4e4502000000"

  Scenario: Nested Any
    Given that a nested Option has an inner type of Any with a value of "fdaddd54a11e45f2bbbae200118567f3"
    Then the inner type is Any with a value of "fdaddd54a11e45f2bbbae200118567f3"
    And the bytes are "0101000000030000004f4e4502000000"
    Given that the nested Option is deployed in a transfer
    And the transfer containing the nested Option is successfully executed
    When the Option is read from the deploy
    Then the inner type is Any with a value of "fdaddd54a11e45f2bbbae200118567f3"
    And the bytes are "0101000000030000004f4e4502000000"

