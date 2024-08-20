package forwardtest

import (
	"testing"

	"github.com/lerenn/cryptellation/pkg/models/account"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/suite"
)

func TestForwardTestSuite(t *testing.T) {
	suite.Run(t, new(ForwardTestSuite))
}

type ForwardTestSuite struct {
	suite.Suite
}

func (suite *ForwardTestSuite) TestGetAccountsSymbols() {
	cases := []struct {
		Accounts map[string]account.Account
		Expected []string
	}{
		// Only one account with one symbol
		{
			Accounts: map[string]account.Account{
				"exchange": {
					Balances: map[string]float64{"DAI": 1000},
				},
			},
			Expected: []string{"DAI"},
		},
		// One account with 3 different symbols
		{
			Accounts: map[string]account.Account{
				"exchange": {
					Balances: map[string]float64{
						"DAI":  1000,
						"USDT": 1000,
						"BTC":  1000,
					},
				},
			},
			Expected: []string{"BTC", "DAI", "USDT"},
		},
		// Two accounts with common symbols
		{
			Accounts: map[string]account.Account{
				"exchange1": {
					Balances: map[string]float64{
						"DAI":  1000,
						"USDT": 1000,
						"BTC":  1000,
					},
				},
				"exchange2": {
					Balances: map[string]float64{
						"DAI":  1000,
						"USDT": 1000,
						"ETH":  1000,
					},
				},
			},
			Expected: []string{"BTC", "DAI", "ETH", "USDT"},
		},
	}

	less := func(a, b string) bool { return a < b }
	for _, c := range cases {
		ft := ForwardTest{
			Accounts: c.Accounts,
		}
		suite.Require().True(cmp.Diff(c.Expected, ft.GetAccountsSymbols(), cmpopts.SortSlices(less)) == "")
	}
}
