package gogrs

// TWSE base url.
const (
	TWSEURL    string = "http://mis.tse.com.tw/"
	TWSEHOST   string = "http://www.twse.com.tw/"
	TWSEOTCCSV string = "/ch/stock/aftertrading/daily_trading_info/st43_download.php?d=%d/%02d&stkno=%s&r=%d"
	TWSECSV    string = "/ch/trading/exchange/STOCK_DAY/STOCK_DAY_print.php?genpage=genpage/Report%d%02d/%d%02d_F3_1_8_%s.php&type=csv&r=%d"
	TWSEREAL   string = "/stock/api/getStockInfo.jsp?ex_ch=%s_%s.tw_%s&json=1&delay=0&_=%d"
)
