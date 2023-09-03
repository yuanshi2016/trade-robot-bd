package biz

import (
	"testing"
	"trade-robot-bd/app/wallet-svc/internal/dao"
)

func TestWalletService_AddIfcBalance(t *testing.T) {
	srv := &WalletRepo{
		dao:          dao.New(),
		cacheService: nil,
		binance:      nil,
		UserSrv:      nil,
	}
	if err := srv.AddIfcBalance("1273211817757249536", "hhhhhhhh", "register", "", 1.0); err != nil {
		t.Errorf("错误 %v", err.Error())
	}
	t.Log("okok")
}
