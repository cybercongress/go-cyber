# Errors

The resources module may return the following errors:

| Type                     | Code  | Description                         |
| ------------------------ | ------| ----------------------------------- |
| ErrTimeLockCoins         | 2     | error time lock coins               |  
| ErrIssueCoins            | 3     | error issue coins                   |
| ErrMintCoins             | 4     | error mint coins                    |
| ErrBurnCoins             | 5     | error burn coins                    |
| ErrSendMintedCoins       | 6     | error send minted coins             |    
| ErrNotAvailablePeriod    | 7     | not available period                |    
| ErrInvalidAccountType    | 8     | receiver account type not supported |    
| ErrAccountNotFound       | 9     | account not found                   |    
| ErrResourceNotExist      | 10    | resource not exist                  |    
| ErrFullSlots             | 11    | all slots are full                  |
| ErrSmallReturn           | 12    | small resource's return amount      |
| ErrInvalidBaseResource   | 13    | invalid base resource               |
 