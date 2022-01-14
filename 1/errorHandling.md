# Error

[Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments/5a40ba36d388ff1b8b2dd4c1c3fe820b8313152f): `Indent Error Flow`

Try to keep the normal code path at a minimal indentation, and indent the error handling, dealing with it first. This improves the readability of the code by permitting visually scanning the normal path quickly.

For instance:
![](pic/IndentErrorFlow1.png)
![](pic/IndentErrorFlow2.png)

## Example 


### Source code
```go
if cur_game_stage == game.EGST_FEATURE_GAME {
	if cached_coin_size_index, cached_num_lines, err := GetCachedClientSpinRequest(csr.UserID, cache); err == nil {
		if _, err = game_conf.GetCoinSizeByIndex(cached_coin_size_index); err == nil {
			if game_conf.IsValidNumLines(cached_num_lines) == true {
				coin_size, _ = game_conf.GetCoinSizeByIndex(cached_coin_size_index)
				num_lines = cached_num_lines
				bet_money = coin_size * float64(num_lines)
			} else {
				c.JSON(http.StatusBadRequest, ginutils.Resp{
					Status: core.EseInvalidFields.Int(),
				})

				return
			}
		} else {
			c.JSON(http.StatusBadRequest, ginutils.Resp{
				Status: core.EseInvalidFields.Int(),
				Msg:    err.Error(),
			})

			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, ginutils.Resp{
			Status: core.EseAccessDB.Int(),
			Msg:    err.Error(),
		})

		return
	}
}
```

### Pseudocode
```go
if 處於 Feature 階段 {
	if CoinSizeIndex, NumLines, err := 獲取緩存請求數據(); err == nil {
		if _, err := 獲得具體CoinSize(); err == nil{
			if 是合理的NumLines嗎() == true {
				SetCoinSize()
				SetNumLines()
				SetBetMoney()
			} else {
				JSON(BadRequest, 錯誤的Fiedls)
				return
			}
		} else {
			JSON(BadRequest, 錯誤的Fiedls)
			return		
		}
	} else {
		JSON(BadRequest, 錯誤的AccessDB)
		return	
	}
}
```

### Final 
```go
if 處於 Feature 階段 {
	CoinSizeIndex, NumLines, err := 獲取緩存請求數據();
	if err != nil {
		JSON(BadRequest, 錯誤的AccessDB)
		return	
	} 

	_, err := 獲得具體CoinSize()
	if err != nil{
		JSON(BadRequest, 錯誤的Fiedls)
		return		
	}	

	if 是合理的NumLines嗎() != true {
		JSON(BadRequest, 錯誤的Fiedls)
		return
	} 

	SetCoinSize()
	SetNumLines()
	SetBetMoney()
}
```