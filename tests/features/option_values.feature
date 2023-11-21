Feature: Option Values

  Scenario: Option Values
    Given that an Option value has an empty value
    Then the Option value is not present
    And the Option value's bytes are ""

    Given an Option value contains a "Bool" value of "true"
    Then the Option value's bytes are "0101"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "0101"
    And the type of the Option is "Bool" with a value of "true"

    Given an Option value contains a "ByteArray" value of "d2029649"
    Then the Option value's bytes are "01d2029649"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "01d2029649"

    Given an Option value contains a "I32" value of "1234567890"
    Then the Option value's bytes are "01d2029649"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "01d2029649"

    Given an Option value contains a "I64" value of "102030405060708090"
    Then the Option value's bytes are "01fac6eee11c7c6a01"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "01fac6eee11c7c6a01"

    Given an Option value contains a "String" value of "The quick brown fox jumps over the lazy dog."
    Then the Option value's bytes are "012c00000054686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f672e"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "012c00000054686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f672e"

    Given an Option value contains a "U8" value of "15"
    Then the Option value's bytes are "010f"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "010f"

    Given an Option value contains a "U32" value of "1234567890"
    Then the Option value's bytes are "01d2029649"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "01d2029649"

    Given an Option value contains a "U64" value of "102030405060708090"
    Then the Option value's bytes are "01fac6eee11c7c6a01"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "01fac6eee11c7c6a01"

    Given an Option value contains a "U128" value of "102030405060708090"
    Then the Option value's bytes are "0108fac6eee11c7c6a01"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "0108fac6eee11c7c6a01"

    Given an Option value contains a "U256" value of "102030405060708090102030405060708090"
    Then the Option value's bytes are "010ffac696b1fd5d40dc69265dc37ca613"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "010ffac696b1fd5d40dc69265dc37ca613"

    Given an Option value contains a "U512" value of "102030405060708090102030405060708090102030405060708090"
    Then the Option value's bytes are "0117fac696b19d1300622227dcf324306f3cdea13d09b41001"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "0117fac696b19d1300622227dcf324306f3cdea13d09b41001"

    Given an Option value contains a "Any" value of "102030405060708090102030405060708090102030405060708090"
    Then the Option value's bytes are "01102030405060708090102030405060708090102030405060708090"
    Given that the Option value is deployed in a transfer as a named argument
    And the transfer containing the Option value is successfully executed
    When the Option is read from the deploy
    Then the Option value's bytes are "01102030405060708090102030405060708090102030405060708090"
