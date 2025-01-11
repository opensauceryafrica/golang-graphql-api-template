package enum

type WalletType string

const (
	FIAT   WalletType = "FIAT"
	CRYPTO WalletType = "CRYPTO"
)

type CryptoNetwork string

const (
	ETHEREUM CryptoNetwork = "ETHEREUM"
	BITCOIN  CryptoNetwork = "BITCOIN"
	SOLANA   CryptoNetwork = "SOLANA"
	TRON     CryptoNetwork = "TRON"
)
