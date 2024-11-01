package token

import (
	"database/sql"
	"strconv"
	"time"

	"titan-container-platform/chain"
	"titan-container-platform/core/dao"
	"titan-container-platform/errors"
)

const (
	hourlyQuota     = 10000
	maxUserQuota    = 400
	quotaResetHours = 1
)

func getCurrentHour() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
}

// ClaimTokens allows the specified account to claim tokens and returns the number of tokens claimed.
func ClaimTokens(account string) (int, error) {
	currentHour := getCurrentHour()

	distributedAmount, err := dao.GetOrCreateHourlyQuota(currentHour)
	if err != nil {
		return errors.InternalServer, err
	}

	if distributedAmount >= hourlyQuota {
		return errors.QuotaIssued, nil
	}

	lastClaim, err := dao.GetLastClaimByAccount(account)
	if err != nil && err != sql.ErrNoRows {
		return errors.InternalServer, err
	}

	if !lastClaim.IsZero() {
		if isToday(lastClaim) {
			return errors.Received, nil
		}
	}

	if distributedAmount+maxUserQuota > hourlyQuota {
		return errors.QuotaIssued, nil
	}

	if err == sql.ErrNoRows {
		err = dao.SetAccountClaims(account, maxUserQuota)
	} else {
		err = dao.UpdateAccountClaims(account, maxUserQuota)
	}
	if err != nil {
		return errors.InternalServer, err
	}

	err = dao.UpdateHourlyQuota(maxUserQuota, currentHour)
	if err != nil {
		return errors.InternalServer, err
	}

	err = chain.ClaimTokens(account, strconv.Itoa(maxUserQuota))
	if err != nil {
		return errors.InternalServer, err
	}

	return errors.Success, nil
}

func isToday(t time.Time) bool {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	tomorrow := today.AddDate(0, 0, 1)
	return t.After(today) && t.Before(tomorrow)
}

// GetBalance retrieves the balance for a given account.
func GetBalance(account string) (string, error) {
	return chain.GetBalance(account)
}
