/**
 * @author dengmeiyu
 * @since 20180711
 */
package alipay

import (
	"testing"
	"fmt"
)

func TestAliTrade(t *testing.T) {
	var privateKey= "MIIEpAIBAAKCAQEA8DIJqC8hNAakgp2ihz22uQHQcB1C64VL0akBzMeVFpLlNqSEkJox0vROxzw/G4/2/dv21lIzWzrg9KvIYBIZQEoBsMaL7kSr7864sIsbOFwecjJB3y8Z0lwTNM/omM0sExTumnQ98Da3Z0BuVakb0EC4T9iwbMIPy6N3BoQNL0CdxMS+11H9NdUyfjrdaj+L7tHRnb2PTCT9oNd7dOdh+1HU0R4haD5H7zZd/H5354HX9fRRoteIE+RLlwENuk4qp+4DyQ+J75OH/VJdTUjZrH3naL7stacmcYBVYTq+pHRJzpiiQX6cheABuTcE0/vd4bHTgFLw800hNv1lSETsPQIDAQABAoIBAF2TYOfA1LKsE9M/Dl6Z0GHrLwU+oR+YYp1ftm/NIWwV9m5/UZX0PIzu2PINqphWUec8lsGQbKmSGqm3cuzaoDPHoTj5YPKGlmVqp9E/oG8olEPxCHCmrOyjKmhcx0UgSZd1hC9wMeEvr+vt0g5VP5y2WhfhV3pHcFfpaymCsJBDXvKnfzGWaNYpPoz1LXqy2hdmPOxsSOwXRdP/PPvO46KzvrK3NqwlzVQPzElMwpjIt50oZh0z+l0aO2bgHzD1m7Gt1Xh9SUzqFfp4kqUmGWX89tXFJYsy4ev+j4ZysfuJCNMpmfhaGB02M88IRZ4D9Ab8xTpkzKvBgo72EltkxxECgYEA/RVLYN3QKRGKTm79XP925eEh1Msz+aD1rsUlSvcP/yWjIQs7zsG/BWzrQOenZyNJEc4aKySqLSgY6cSAreAgWenEeFtyZzQmsQWe9KH7FS/49fzRnBP3JhjIBS60ESYP20gcSlKDQlGNAMvRqz0HVsif04uzwoTHJa94o95sbEMCgYEA8va4CPpfl++m41sLWBI5qTMF+QaVZw+mUzJvkAjunQNo+6jJpLb55AoVNPlRB86oBoxHXiUKGxnKahz77AOFUkEDu5KlBMpkG4tOsK7n6VfrIFGeB/7FPx/iWYg2I7SWCjeIH5q3pdoufxcQ1seCKdxhkavEnDKfWMW3rXc//X8CgYEAunwNbhQkBY8CLadFFFi9oMgSaL1O0BtVzXFBeIqyg9yU5o1jhYdoHTRT6SCJTstGVVNcHvxGVT3dlauQ8g5baEWD2vfvRbK86+XrafFNlSAjQAcJ4QspKy5JfOAcGSLFvlvVVMKWK7DxyGtnVNmEZeMxOe0QCT5TjCZPK9iCZgcCgYEAkbunG8uCN1JKWiksHsGf0HuIY5ytVMowS8r+2/hfl9KJ0BmoCaKvNTdPDR1Wm0Y7xuGxlSjGbQcFQKzt9t0NxQ62PHZzgPIdJeBjbNscw/w2ToZmMgmBKqHnVSi8wKH7NVmlzr8w1MyQAy9ErG+zBYTpCUVsgvxiVA8UY3oZ6eUCgYBeY/Inz/Opa3dGlgAv8d6iFaeSGTqLGrPYJ6+wf22uPxzChHyw+o6IONv9E4zr/epZYtA4SEnWVkYcEO2eXuBtmPjlJjNpdq/q/8RTM7gKFR3+r2+M8x3zMvsL5YSroDabfDJpTPcDAt+qZLdhwgy3d2TWqLXfTWh30XOn4LMAeA=="

	app := new(AliTrade)

	app.InitBaseConfig(&BaseAliConfig{
		AppId:    "2018010401585047",
		AppAuthToken:"201801BB635d7c6ab47a459ba579f876d3736X68",
	})

	ret, err := app.Handle(map[string]interface{}{
		"body":           "123",
		"out_trade_no":   "2018061510564487167",
		"trans_currency": "CNY",
		"total_amount":   "0.01",
		"auth_code":      "286050206957832927",
		"scene":          "bar_code",
		"buyer_id":       "1",
		"subject":        "123",

	},privateKey)

	fmt.Printf("%+v", ret.AliPayTradePay.Code)

	fmt.Println(err)
}

