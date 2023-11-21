Feature: Nested CL List Values

  Scenario: Bool Lists
    Given a list is created with "Bool" values of ("true", "false", "true")
    Then the list's bytes are "03000000010001"
    And the list's length is 3
    And the list's "1st" item is a CLValue with "Bool" value of "true"
    And the list's "2nd" item is a CLValue with "Bool" value of "false"
    And the list's "3rd" item is a CLValue with "Bool" value of "true"
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000010001"
    And the list's length is 3
    And the list's "1st" item is a CLValue with "Bool" value of "true"
    And the list's "2nd" item is a CLValue with "Bool" value of "false"
    And the list's "3rd" item is a CLValue with "Bool" value of "true"

  Scenario: I32 Lists
    Given a list is created with I32 values of (1, 2, 3)
    Then the list's bytes are "03000000010000000200000003000000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with I32 value of 1
    And the list's "2nd" item is a CLValue with I32 value of 2
    And the list's "3rd" item is a CLValue with I32 value of 3
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000010000000200000003000000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with I32 value of 1
    And the list's "2nd" item is a CLValue with I32 value of 2
    And the list's "3rd" item is a CLValue with I32 value of 3

  Scenario: I64 Lists
    Given a list is created with I64 values of (1001, 2002, 3003)
    Then the list's bytes are "03000000e903000000000000d207000000000000bb0b000000000000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with I64 value of 1001
    And the list's "2nd" item is a CLValue with I64 value of 2002
    And the list's "3rd" item is a CLValue with I64 value of 3003
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000e903000000000000d207000000000000bb0b000000000000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with I64 value of 1001
    And the list's "2nd" item is a CLValue with I64 value of 2002
    And the list's "3rd" item is a CLValue with I64 value of 3003

  Scenario: String Lists
    Given a list is created with "String" values of ("ONE", "TWO", "THREE")
    Then the list's bytes are "03000000030000004f4e450300000054574f050000005448524545"
    And the list's length is 3
    And the list's "1st" item is a CLValue with "String" value of "ONE"
    And the list's "2nd" item is a CLValue with "String" value of "TWO"
    And the list's "3rd" item is a CLValue with "String" value of "THREE"
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000030000004f4e450300000054574f050000005448524545"
    And the list's length is 3
    And the list's "1st" item is a CLValue with "String" value of "ONE"
    And the list's "2nd" item is a CLValue with "String" value of "TWO"
    And the list's "3rd" item is a CLValue with "String" value of "THREE"

  Scenario: U8 Lists
    Given a list is created with U8 values of (4, 5, 6)
    Then the list's bytes are "03000000040506"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U8 value of 4
    And the list's "2nd" item is a CLValue with U8 value of 5
    And the list's "3rd" item is a CLValue with U8 value of 6
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000040506"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U8 value of 4
    And the list's "2nd" item is a CLValue with U8 value of 5
    And the list's "3rd" item is a CLValue with U8 value of 6
    Given that the list is deployed in a transfer

  Scenario: U32 Lists
    Given a list is created with U32 values of (700, 800, 900)
    Then the list's bytes are "03000000bc0200002003000084030000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U32 value of 700
    And the list's "2nd" item is a CLValue with U32 value of 800
    And the list's "3rd" item is a CLValue with U32 value of 900
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000bc0200002003000084030000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U32 value of 700
    And the list's "2nd" item is a CLValue with U32 value of 800
    And the list's "3rd" item is a CLValue with U32 value of 900

  Scenario: U64 Lists
    Given a list is created with U64 values of (170000000, 180000000, 190000000)
    Then the list's bytes are "0300000080fe210a000000000095ba0a00000000802b530b00000000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U64 value of 170000000
    And the list's "2nd" item is a CLValue with U64 value of 180000000
    And the list's "3rd" item is a CLValue with U64 value of 190000000
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "0300000080fe210a000000000095ba0a00000000802b530b00000000"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U64 value of 170000000
    And the list's "2nd" item is a CLValue with U64 value of 180000000
    And the list's "3rd" item is a CLValue with U64 value of 190000000

  Scenario: U128 Lists
    Given a list is created with U128 values of (770000000, 880000000, 990000000)
    Then the list's bytes are "03000000048044e52d0400bc7334048033023b"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U128 value of 770000000
    And the list's "2nd" item is a CLValue with U128 value of 880000000
    And the list's "3rd" item is a CLValue with U128 value of 990000000
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000048044e52d0400bc7334048033023b"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U128 value of 770000000
    And the list's "2nd" item is a CLValue with U128 value of 880000000
    And the list's "3rd" item is a CLValue with U128 value of 990000000

  Scenario: U256 Lists
    Given a list is created with U256 values of (1, 2, 3)
    Then the list's bytes are "03000000010101020103"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U256 value of 1
    And the list's "2nd" item is a CLValue with U256 value of 2
    And the list's "3rd" item is a CLValue with U256 value of 3
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000010101020103"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U256 value of 1
    And the list's "2nd" item is a CLValue with U256 value of 2
    And the list's "3rd" item is a CLValue with U256 value of 3

  Scenario: U512 Lists
    Given a list is created with U512 values of (1, 2, 3)
    Then the list's bytes are "03000000010101020103"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U512 value of 1
    And the list's "2nd" item is a CLValue with U512 value of 2
    And the list's "3rd" item is a CLValue with U512 value of 3
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "03000000010101020103"
    And the list's length is 3
    And the list's "1st" item is a CLValue with U512 value of 1
    And the list's "2nd" item is a CLValue with U512 value of 2
    And the list's "3rd" item is a CLValue with U512 value of 3

  Scenario: Nested String U256 Lists
    Given a nested list is created with U256 values of ((1, 2, 3),(4, 5, 6))
    Then the list's bytes are "020000000300000001010102010303000000010401050106"
    And the list's length is 2
    And the "1st" nested list's "1st" item is a CLValue with U256 value of 1
    And the "2nd" nested list's "3rd" item is a CLValue with U256 value of 6
    Given that the list is deployed in a transfer
    And the transfer containing the list is successfully executed
    When the list is read from the deploy
    Then the list's bytes are "020000000300000001010102010303000000010401050106"
    And the list's length is 2
    And the "1st" nested list's "1st" item is a CLValue with U256 value of 1
    And the "2nd" nested list's "3rd" item is a CLValue with U256 value of 6
