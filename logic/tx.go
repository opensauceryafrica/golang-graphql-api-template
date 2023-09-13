package logic

import (
	"blacheapi/http/graphql/model"
	"blacheapi/primer/constant"
	"blacheapi/primer/enum"
	"blacheapi/primer/function"
	"blacheapi/primer/gql"
	"blacheapi/primer/primitive"
	"blacheapi/primer/typing"
	"blacheapi/repository"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/uptrace/bun/schema"
)

// Webhook performs post-fuding & post-withdrawal processing on the input and makes a repository call to update balances
func Webhook(input model.WebhookInput, debitSavings bool) (*repository.Transaction, error) {

	// start a database transaction
	btx, err := repository.BeginBlacheTx()
	if err != nil {
		return nil, err
	}
	defer btx.Rollback()

	// the transaction
	var tx repository.Transaction

	// find saving
	filter := typing.SQLMaps{
		WMaps: []typing.SQLMap{
			{
				Map: map[string]interface{}{
					"reference": input.Reference,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},
		WJoinOperator: enum.And,
	}

	// find tx and lock
	err = tx.FUByMap(btx, filter, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found!`, input.Reference), 404)
		}
		return nil, err
	}

	if tx.Paid && !tx.Failed {
		return &tx, nil
	}

	var txfailed, cardtx *bool

	card := repository.Card{}

	if tx.Gateway != "" && tx.Gateway == enum.Flutterwave && input.Provider == model.EPaymentGateway(enum.Flutterwave) {
		// handle Flutterwave tx

		txresponse := input.Tx.(repository.FlutterwaveTx)

		// if Flutterwave can't find the tx
		if txresponse.Status == string(enum.FlutterwaveError) {
			return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found on %s!`, input.Reference, strings.ToTitle(input.Provider.String())), 404)
		}

		// if Flutterwave tx failed
		if txresponse.Data.Status == string(enum.FlutterwaveFailed) {
			// set failed flag
			txfailed = function.PBool(true)
		}

		// if Flutterwave tx successful
		if txresponse.Status == string(enum.FlutterwaveSuccess) && txresponse.Data.Status == string(enum.FlutterwaveSuccessful) && txresponse.Data.Currency == string(tx.Currency) && txresponse.Data.AmountSettled >= tx.Amount {

			txfailed = function.PBool(false)

			// blache card (if any)
			if txresponse.Data.Card != nil && txresponse.Data.Card.Token != "" {
				checksum := function.StringSha256(function.Stringify(txresponse.Data.Card))

				// check if checksum exists for a card
				err = card.FByMap(typing.SQLMaps{
					WMaps: []typing.SQLMap{
						{
							Map: map[string]interface{}{
								"checksum": checksum,
							},
							JoinOperator:       enum.And,
							ComparisonOperator: enum.Equal,
						},
					},
					WJoinOperator: enum.And,
				}, true)
				if err != nil && err != sql.ErrNoRows {
					return nil, err
				} else if err == sql.ErrNoRows {
					// create card
					card = repository.Card{
						ID:           function.GenerateUUID(),
						UserID:       tx.UserID,
						CustomerID:   tx.CustomerID,
						First6digits: txresponse.Data.Card.First6digits,
						Last4digits:  txresponse.Data.Card.Last4digits,
						Issuer:       txresponse.Data.Card.Issuer,
						Country:      txresponse.Data.Card.Country,
						Type:         txresponse.Data.Card.Type,
						Expiry:       txresponse.Data.Card.Expiry,
						Checksum:     checksum,
						Token:        txresponse.Data.Card.Token,
					}

					// set create card flag
					cardtx = function.PBool(true)
				}
			}

		}
	} else if tx.Gateway != "" && tx.Gateway == enum.Paystack && input.Provider == model.EPaymentGateway(enum.Paystack) {
		// handle Paystack tx

		txresponse := input.Tx.(repository.PaystackTx)

		// if Paystack can't find the tx
		if !txresponse.Status {
			return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found on %s!`, input.Reference, strings.ToTitle(input.Provider.String())), 404)
		}

		// if Paystack tx failed
		if txresponse.Data.Status == string(enum.PaystackFailed) {
			// set failed flag
			txfailed = function.PBool(true)
		}

		// if Paystack tx successful
		// @TODO: is Data.Amount the amount settled or the amount paid?
		if txresponse.Status && txresponse.Data.Status == string(enum.PaystackSuccessful) && txresponse.Data.Currency == string(tx.Currency) && txresponse.Data.Amount >= tx.Amount {

			txfailed = function.PBool(false)

			// blache card (if any)
			if txresponse.Data.Authorization != nil && txresponse.Data.Authorization.AuthorizationCode != "" {
				checksum := function.StringSha256(function.Stringify(txresponse.Data.Authorization))

				// check if checksum exists for a card
				err = card.FByMap(typing.SQLMaps{
					WMaps: []typing.SQLMap{
						{
							Map: map[string]interface{}{
								"checksum": checksum,
							},
							JoinOperator:       enum.And,
							ComparisonOperator: enum.Equal,
						},
					},
					WJoinOperator: enum.And,
				}, true)
				if err != nil && err != sql.ErrNoRows {
					return nil, err
				} else if err == sql.ErrNoRows {
					// create card
					card = repository.Card{
						ID:           function.GenerateUUID(),
						UserID:       tx.UserID,
						CustomerID:   tx.CustomerID,
						First6digits: "",
						Last4digits:  txresponse.Data.Authorization.Last4,
						Issuer:       txresponse.Data.Authorization.Bank,
						Country:      txresponse.Data.Authorization.CountryCode,
						Type:         strings.ToUpper(txresponse.Data.Authorization.Brand),
						Expiry:       fmt.Sprintf("%s/%s", txresponse.Data.Authorization.ExpMonth, txresponse.Data.Authorization.ExpYear),
						Checksum:     checksum,
						Token:        txresponse.Data.Authorization.AuthorizationCode,
					}

					// set create card flag
					cardtx = function.PBool(true)
				}
			}

		}

	} else if tx.Gateway != "" && tx.Gateway == enum.MANUAL && input.Provider == model.EPaymentGateway(enum.MANUAL) {

		// handle manual tx
		txfailed = function.PBool(false)

	} else if tx.Gateway == "" && tx.Method == enum.Wallet && input.Provider == model.EPaymentGateway(enum.Wallet) {
		// handle wallet tx

		txresponse := input.Tx.(repository.WalletTx)

		// if wallet can't find the tx (or some other error)
		if !txresponse.Status {
			return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found on %s!`, input.Reference, strings.ToTitle(input.Provider.String())), 404)
		}

		// if wallet tx failed
		if txresponse.Data.Status == string(enum.WalletFailed) {
			// set failed flag
			txfailed = function.PBool(true)
		}

		// if wallet tx successful
		if txresponse.Status && txresponse.Data.Status == string(enum.WalletSuccess) && txresponse.Data.Currency == string(tx.Currency) && txresponse.Data.Amount >= tx.Amount {

			txfailed = function.PBool(false)

		}

	} else {
		return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" is unknow to configure.Webhook`, input.Reference), 400)
	}

	// update tx, savings account, and card as necessary
	if txfailed != nil && *txfailed {
		// lock the savings account to ensure the post-balance doesn't change until the tx is updated

		s := repository.Transaction{}
		err = s.FUByMap(btx, typing.SQLMaps{
			WMaps: []typing.SQLMap{
				{
					Map: map[string]interface{}{
						"savings.id": tx.SavingsID,
					},
					JoinOperator:       enum.And,
					ComparisonOperator: enum.Equal,
				},
			},
			WJoinOperator: enum.And,
		}, true)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, gql.MakeError(fmt.Sprintf(`Savings account with id "%s" not found!`, tx.SavingsID), 404)
			}
			return nil, err
		}

		// update tx
		query := typing.SQLMaps{
			WMaps: []typing.SQLMap{
				{
					Map: map[string]interface{}{
						"id": tx.ID,
					},
					JoinOperator:       enum.And,
					ComparisonOperator: enum.Equal,
				},
			},
			SMap: typing.SQLMap{
				Map: map[string]interface{}{
					"post_balance": s.SavingsAmount,
					"failed":       true,
					"failed_at":    "now()",
					"final_config": input.Tx,
					"paid":         false,
					"paid_at":      schema.NullTime{},
					"history":      fmt.Sprintf(`history || '{"act": "%s", "by": "%s", "at": "%s"}'::jsonb`, fmt.Sprintf("%s transaction failed", strings.ToTitle(string(tx.Type))), "system", schema.NullTime{Time: time.Now()}),

					"updated_at": "now()", // update updated_at
				},
				JoinOperator:       enum.Comma,
				ComparisonOperator: enum.Equal,
			},
			WJoinOperator: enum.And,
		}
		// update tx
		err = tx.UByMapTx(btx, query)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found!`, input.Reference), 404)
			}
			return nil, err
		}

		if err := btx.Commit(); err != nil {
			return nil, err
		}

		return &tx, nil

	} else if txfailed != nil && !*txfailed {
		// lock the savings account to ensure the post-balance doesn't change until the tx is updated

		s := repository.Transaction{}
		err = s.FUByMap(btx, typing.SQLMaps{
			WMaps: []typing.SQLMap{
				{
					Map: map[string]interface{}{
						"savings.id": tx.SavingsID,
					},
					JoinOperator:       enum.And,
					ComparisonOperator: enum.Equal,
				},
			},
			WJoinOperator: enum.And,
		}, true)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, gql.MakeError(fmt.Sprintf(`Savings account with id "%s" not found!`, tx.SavingsID), 404)
			}
			return nil, err
		}

		// update savings account if tx is a funding tx
		if tx.Type == enum.FUNDING {
			umap := map[string]interface{}{
				"updated_at":     "now()", // update updated_at
				"amount_blached": fmt.Sprintf(`amount_blached + %v`, tx.SavingsAmount),
			}

			err = s.UByMapTx(btx, typing.SQLMaps{
				WMaps: []typing.SQLMap{
					{
						Map: map[string]interface{}{
							"savings.id": s.ID,
						},
						JoinOperator:       enum.And,
						ComparisonOperator: enum.Equal,
					},
				},
				SMap: typing.SQLMap{
					Map:                umap,
					JoinOperator:       enum.Comma,
					ComparisonOperator: enum.Equal,
				},
				RMap: typing.SQLMap{
					Map: map[string]interface{}{"*": nil},
				},
				WJoinOperator: enum.And,
			})
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, gql.MakeError(fmt.Sprintf(`Savings account with id "%s" retrieval failed!`, tx.SavingsID), 500)
				}
				return nil, err
			}
		}

		// update the savings account if tx is withdrawal tx and debitSavings is true
		if tx.Type == enum.WITHDRAWAL && debitSavings {
			umap := map[string]interface{}{
				"updated_at":       "now()", // update updated_at
				"amount_withdrawn": fmt.Sprintf(`amount_withdrawn + %v`, tx.SavingsAmount),
			}

			err = s.UByMapTx(btx, typing.SQLMaps{
				WMaps: []typing.SQLMap{
					{
						Map: map[string]interface{}{
							"savings.id": s.ID,
						},
						JoinOperator:       enum.And,
						ComparisonOperator: enum.Equal,
					},
				},
				SMap: typing.SQLMap{
					Map:                umap,
					JoinOperator:       enum.Comma,
					ComparisonOperator: enum.Equal,
				},
				RMap: typing.SQLMap{
					Map: map[string]interface{}{"*": nil},
				},
				WJoinOperator: enum.And,
			})
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, gql.MakeError(fmt.Sprintf(`Savings account with id "%s" retrieval failed!`, tx.SavingsID), 500)
				}
				return nil, err
			}
		}

		// update tx
		err = tx.UByMapTx(btx, typing.SQLMaps{
			WMaps: []typing.SQLMap{
				{
					Map: map[string]interface{}{
						"id": tx.ID,
					},
					JoinOperator:       enum.And,
					ComparisonOperator: enum.Equal,
				},
			},
			SMap: typing.SQLMap{
				Map: map[string]interface{}{
					"post_balance": s.SavingsAmount,
					"failed":       false,
					"failed_at":    schema.NullTime{},
					"final_config": input.Tx,
					"paid":         true,
					"paid_at":      "now()",
					"history":      fmt.Sprintf(`history || '{"act": "%s", "by": "%s", "at": "%s"}'::jsonb`, fmt.Sprintf("%s transaction successful", strings.ToTitle(string(tx.Type))), "system", schema.NullTime{Time: time.Now()}),

					"updated_at": "now()", // update updated_at
				},
				JoinOperator:       enum.Comma,
				ComparisonOperator: enum.Equal,
			},
			WJoinOperator: enum.And,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found!`, input.Reference), 404)
			}
			return nil, err
		}

		if cardtx != nil && *cardtx {
			err = card.CreateTx(btx)
			if err != nil {
				return nil, err
			}
		}

		if err := btx.Commit(); err != nil {
			return nil, err
		}

		return &tx, nil
	}

	return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" is unknow to configure.Blache`, input.Reference), 400)
}

// Transaction finds and returns a transaction by either of the given input parameters (id, reference)
func Transaction(input model.TransactionFilterInput) (*repository.Transaction, error) {
	t := repository.Transaction{}

	// find transaction
	filter := typing.SQLMaps{
		WMaps: []typing.SQLMap{
			{
				Map:                map[string]interface{}{},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},
		WJoinOperator: enum.And,
	}
	preload := true
	join := false
	if input.TransactionID != nil {
		filter.WMaps[0].Map["transactions.id"] = *input.TransactionID
		err := t.FByMap(filter, preload, join)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, gql.MakeError(fmt.Sprintf(`Transaction with id "%s" not found!`, *input.ProductID), 404)
			}
			return nil, err
		}
		return &t, nil
	}
	if input.Reference != nil {
		filter.WMaps[0].Map["transactions.reference"] = *input.Reference
		err := t.FByMap(filter, preload, join)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, gql.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found!`, *input.Reference), 404)
			}
			return nil, err
		}
		return &t, nil
	}
	return nil, gql.MakeError(`Transaction not found!`, 404)
}

// Transactions finds and returns a list of transactions matching the given input parameters. The transactions returned are restricted to just those belonging to the customer in the request context
func Transactions(input *model.TransactionFilterInput, customer repository.Customer) (repository.Transactions, *repository.Pagination, error) {
	ts := make(repository.Transactions, 0)
	filter := typing.SQLMap{
		Map: map[string]interface{}{
			"customer_id": customer.OrgID,
		},
		JoinOperator:       enum.And,
		ComparisonOperator: enum.Equal,
	}
	search := typing.SQLMap{
		Map:                map[string]interface{}{},
		JoinOperator:       enum.Or,
		ComparisonOperator: enum.ILike,
	}

	// find transaction
	preload := true
	join := false
	offset := constant.PageOffset

	if input == nil {
		// default values
		input = &model.TransactionFilterInput{
			Limit: function.PInt(constant.PageLimit),
			Page:  function.PInt(constant.PageOffset + 1),
			Sort:  &model.AllESort[0],
		}
	}
	if input.Limit == nil {
		input.Limit = function.PInt(constant.PageLimit)
	}
	if input.Page == nil {
		input.Page = function.PInt(constant.PageOffset + 1)
	}
	offset = (*input.Page - 1) * (*input.Limit)
	if input.Sort == nil {
		input.Sort = &model.AllESort[0]
	}

	function.LayerMap(function.StructToMapOfNonNils(input, "json", primitive.Array{"transaction_id", "savings_id", "user_id", "product_id", "customer_id", "type", "reference", "currency", "gateway", "method", "paid", "failed", "cancelled"}, map[string]string{"transaction_id": "id"}), filter.Map)

	if input != nil && input.Search != nil {

		// define search parameters
		search.JoinOperator = enum.Or
		search.ComparisonOperator = enum.ILike
		search.Map["reference"] = fmt.Sprintf("%%%s%%", *input.Search)
		search.Map["remark"] = fmt.Sprintf("%%%s%%", *input.Search)
		search.Map["reference"] = fmt.Sprintf("%%%s%%", *input.Search)
		search.Map["transaction_id"] = fmt.Sprintf("%%%s%%", *input.Search)
	}

	query := typing.SQLMaps{
		WMaps: []typing.SQLMap{
			filter,
			search,
		},
		WJoinOperator: enum.And,
	}

	err := ts.FByMap(query, *input.Limit, offset, input.Sort.String(), preload, join)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	var pagination *repository.Pagination

	if input.Paginate != nil && *input.Paginate {
		// get total count
		total, err := ts.CByMap(query)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, err
		}
		pagination = &repository.Pagination{
			Page:  *input.Page,
			Limit: *input.Limit,
			Pages: int(math.Ceil(float64(total) / float64(*input.Limit))),
		}
	}
	return ts, pagination, nil
}
