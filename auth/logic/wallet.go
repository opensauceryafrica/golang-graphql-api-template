package logic

import (
	"context"
	"database/sql"

	"cendit.io/gate/http/graphql/exception"
	"cendit.io/gate/http/graphql/model"

	"cendit.io/auth/repository"
	"cendit.io/auth/schema"
	"cendit.io/garage/function"
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/xiao"
)

// function for user to add wallet
func AddWallet(input model.WalletInput, user schema.User) (*schema.Wallet, error) {
	query := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"user_id":  user.ID,
					"currency": input.Currency,
				},
				JoinOperator:       xiao.And,
				ComparisonOperator: xiao.Equal,
			},
		},
		RMap: xiao.SQLMap{
			Map: map[string]interface{}{
				"*": nil,
			},
		},
		WJoinOperator: xiao.And,
	}

	// Check if user already has a wallet for this currency
	w, err := repository.Wallet().FindByMap(context.Background(), query, true)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if w.ID != "" {
		return nil, exception.MakeError("You already have a wallet for this currency!", 400)
	}

	if enum.Currency(input.Currency.String()) == enum.NGN {
		w.Name = "Nigerian Naira"
		w.Currency = enum.NGN
		w.Address = []schema.WalletAddress{}
		w.Type = enum.FIAT
	} else if enum.Currency(input.Currency.String()) == enum.BTC {
		w.Name = "Bitcoin"
		w.Currency = enum.BTC
		w.Address = []schema.WalletAddress{}
		w.Type = enum.CRYPTO

		// TODO: generate wallet address via bysect
		w.Address = append(w.Address, []schema.WalletAddress{
			{
				Address: "0x1234567890",
				Network: enum.BITCOIN,
			},
		}...)
	} else if enum.Currency(input.Currency.String()) == enum.USDT {
		w.Name = "Tether"
		w.Currency = enum.USDT
		w.Address = []schema.WalletAddress{}
		w.Type = enum.CRYPTO

		// TODO: generate wallet address via bysect
		w.Address = append(w.Address, []schema.WalletAddress{
			{
				Address: "0x1234567890",
				Network: enum.TRON,
			},
			{
				Address: "0x1234567890",
				Network: enum.SOLANA,
			},
		}...)
	} else {
		return nil, exception.MakeError("You cannot add a wallet for this currency!", 400)
	}

	// generate unique wallet id
	w.ID = function.GenerateUUID()

	w.UserID = user.ID
	w.Date()

	err = repository.Wallet().Create(context.Background(), xiao.SQLMaps{
		IMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"id":         w.ID,
					"user_id":    w.UserID,
					"type":       w.Type,
					"name":       w.Name,
					"balance":    w.Balance,
					"address":    w.Address,
					"currency":   w.Currency,
					"created_at": w.CreatedAt,
					"updated_at": w.UpdatedAt,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return w, nil
}

// function to get signed in user's wallets
func Wallets(input model.WalletFilter, user schema.User) (schema.Wallets, error) {

	query := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"user_id": user.ID,
				},
				JoinOperator:       xiao.And,
				ComparisonOperator: xiao.Equal,
			},
		},
		RMap: xiao.SQLMap{
			Map: map[string]interface{}{
				"*": nil,
			},
		},
		WJoinOperator: xiao.And,
	}

	if input.Type != nil {
		query.WMaps[0].Map["type"] = *input.Type
	}

	w, err := repository.Wallet().FindAllByMap(context.Background(), query, true)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return w, nil
}
