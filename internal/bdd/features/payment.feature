Feature: PIX Payment

  Scenario: Create a new PIX Billing
    Given the payment has valid data with order_id and amount
    When i send a message to create a new pix billing
    Then the pix billing should be created successfully