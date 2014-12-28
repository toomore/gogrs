package gogrs

const (
	TWSEURL  string = "http://mis.tse.com.tw/"
	TWSEHOST string = "http://www.twse.com.tw/"
	TWSECSV  string = "/ch/stock/aftertrading/daily_trading_info/st43_download.php?d=%d/%02d&stkno=%s&r=%d"
	TWSEREAL string = "/stock/api/getStockInfo.jsp?ex_ch=%s_%s.tw_%s&json=1&delay=0&_=%d"
)
