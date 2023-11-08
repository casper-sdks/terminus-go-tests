Feature: Nested CL Map Values

  Scenario: Nested Maps
    Given a map is created {"ONE": 2}
    Then the map's key type is "String" and the maps value type is "U32"
    Then the map's bytes are "01000000030000004f4e4502000000"
    Given that the nested map is deployed in a transfer
    And the transfer containing the nested map is successfully executed
    When the map is read from the deploy
    And the map's bytes are "01000000030000004f4e4502000000"
    Then the map's key type is "String" and the maps value type is "U32"
    And the map's key is "ONE" and value is "2"

    Given a nested map is created {"ONE": {"TWO": 3}, "FOUR": {"FIVE", 6}}
    Then the map's key type is "String" and the maps value type is "Map"
    Then the map's bytes are "02000000030000004f4e45010000000300000054574f0300000004000000464f555201000000040000004649564506000000"
    Given that the nested map is deployed in a transfer
    And the transfer containing the nested map is successfully executed
    When the map is read from the deploy
    Then the map's bytes are "02000000030000004f4e45010000000300000054574f0300000004000000464f555201000000040000004649564506000000"
    And the map's key type is "String" and the maps value type is "Map"
    And the 1st nested map's key is "TWO" and value is "3"

    Given a nested map is created  {1: {11: {111: "ONE_ONE_ONE"}, 12: {121: "ONE_TWO_ONE"}}, 2: {21: {211: "TWO_ONE_ONE"}, 22: {221: "TWO_TWO_ONE"}}}
    Then the map's key type is "U256" and the maps value type is "Map"
    Then the map's bytes are "02000000010102000000010b01000000016f0b0000004f4e455f4f4e455f4f4e45010c01000000790000000b0000004f4e455f54574f5f4f4e4501020200000001150100000001d30b00000054574f5f4f4e455f4f4e45011601000000dd0000000b00000054574f5f54574f5f4f4e45"
    Given that the nested map is deployed in a transfer
    And the transfer containing the nested map is successfully executed
    When the map is read from the deploy
    Then the map's bytes are "02000000010102000000010b01000000016f0b0000004f4e455f4f4e455f4f4e45010c01000000790000000b0000004f4e455f54574f5f4f4e4501020200000001150100000001d30b00000054574f5f4f4e455f4f4e45011601000000dd0000000b00000054574f5f54574f5f4f4e45"
    And the map's key type is "U256" and the maps value type is "Map"
