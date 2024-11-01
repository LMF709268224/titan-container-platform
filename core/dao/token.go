package dao

import (
	"database/sql"
	"fmt"
	"time"
)

func GetOrCreateHourlyQuota(currentHour time.Time) (int, error) {
	queryS := fmt.Sprintf(`SELECT distributed_amount FROM %s WHERE hour = ? `, hourlyQuotasTable)
	queryI := fmt.Sprintf(`INSERT INTO %s (hour, distributed_amount) VALUES (?, 0) `, hourlyQuotasTable)

	var distributedAmount int
	err := mDB.QueryRow(queryS, currentHour).Scan(&distributedAmount)
	if err == sql.ErrNoRows {
		_, err = mDB.Exec(queryI, currentHour)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return distributedAmount, nil
}

func UpdateHourlyQuota(amount int, currentHour time.Time) error {
	query := fmt.Sprintf(`UPDATE %s SET distributed_amount = distributed_amount + ? WHERE hour = ? `, hourlyQuotasTable)
	_, err := mDB.Exec(query, amount, currentHour)
	return err
}

func GetLastClaimByAccount(account string) (time.Time, error) {
	query := fmt.Sprintf(`SELECT %s FROM user_claims WHERE account = ? `, userClaimsTable)
	var lastClaim time.Time
	err := mDB.QueryRow(query, account).Scan(&lastClaim)
	if err != nil && err != sql.ErrNoRows {
		return lastClaim, err
	}

	return lastClaim, err
}

func SetAccountClaims(account string, maxUserQuota int) error {
	query := fmt.Sprintf(`INSERT INTO %s (account, amount, last_claim) VALUES (?, ?, ?) `, userClaimsTable)
	_, err := mDB.Exec(query, account, maxUserQuota, time.Now())
	return err
}

func UpdateAccountClaims(account string, maxUserQuota int) error {
	query := fmt.Sprintf(`UPDATE %s SET amount = amount + ?, last_claim = ? WHERE account = ? `, userClaimsTable)
	_, err := mDB.Exec(query, maxUserQuota, time.Now(), account)
	return err
}
