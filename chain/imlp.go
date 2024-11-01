package chain

// ClaimTokens sends tokens to the specified account.
func ClaimTokens(account string) error {
	// TODO 检查次数
	return faucetSend(account)
}

// GetBalance retrieves the balance for a given account.
func GetBalance(account string) (string, error) {
	return balance(account)
}
