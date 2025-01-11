package logic

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"cendit.io/gate/http/graphql/exception"
	"cendit.io/gate/http/graphql/model"

	"cendit.io/auth/repository"
	"cendit.io/auth/schema"
	"cendit.io/garage/function"
	"cendit.io/garage/primer/constant"
	"cendit.io/garage/primitive"
	garage "cendit.io/garage/schema"
	"cendit.io/garage/xiao"
)

// Webhook performs post-fuding & post-withdrawal processing on the input and makes a repository call to update balances
func Webhook(input model.WebhookInput, debitSavings bool) (*schema.Transaction, error) {

	// start a database transaction
	ptx, err := repository.Transaction().Tx(context.Background())
	if err != nil {
		return nil, err
	}
	defer ptx.Rollback()

	// find saving
	filter := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"reference": input.Reference,
				},
				JoinOperator:       xiao.And,
				ComparisonOperator: xiao.Equal,
			},
		},
		WJoinOperator: xiao.And,
	}

	// find tx and lock
	_, err = repository.Transaction().FindAndLockByMap(context.Background(), ptx, filter, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found!`, input.Reference), 404)
		}
		return nil, err
	}

	return nil, exception.MakeError(fmt.Sprintf(`Transaction with reference "%s" is unknow to configure.Cendit`, input.Reference), 400)
}

// Transaction finds and returns a transaction by either of the given input parameters (id, reference)
func Transaction(input model.TransactionFilterInput) (*schema.Transaction, error) {

	// find transaction
	filter := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map:                map[string]interface{}{},
				JoinOperator:       xiao.And,
				ComparisonOperator: xiao.Equal,
			},
		},
		WJoinOperator: xiao.And,
	}

	if input.TransactionID != nil {
		filter.WMaps[0].Map["transactions.id"] = *input.TransactionID
		t, err := repository.Transaction().FindByMap(context.Background(), filter, true)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, exception.MakeError(fmt.Sprintf(`Transaction with id "%s" not found!`, *input.ProductID), 404)
			}
			return nil, err
		}
		return t, nil
	}
	if input.Reference != nil {
		filter.WMaps[0].Map["transactions.reference"] = *input.Reference
		t, err := repository.Transaction().FindByMap(context.Background(), filter, true)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, exception.MakeError(fmt.Sprintf(`Transaction with reference "%s" not found!`, *input.Reference), 404)
			}
			return nil, err
		}
		return t, nil
	}
	return nil, exception.MakeError(`Transaction not found!`, 404)
}

// Transactions finds and returns a list of transactions matching the given input parameters. The transactions returned are restricted to just those belonging to the customer in the request context
func Transactions(input *model.TransactionFilterInput, customer schema.User) (schema.Transactions, *garage.Pagination, error) {

	filter := xiao.SQLMap{
		Map:                map[string]interface{}{},
		JoinOperator:       xiao.And,
		ComparisonOperator: xiao.Equal,
	}
	search := xiao.SQLMap{
		Map:                map[string]interface{}{},
		JoinOperator:       xiao.Or,
		ComparisonOperator: xiao.ILike,
	}

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
		search.JoinOperator = xiao.Or
		search.ComparisonOperator = xiao.ILike
		search.Map["reference"] = fmt.Sprintf("%%%s%%", *input.Search)
		search.Map["remark"] = fmt.Sprintf("%%%s%%", *input.Search)
		search.Map["reference"] = fmt.Sprintf("%%%s%%", *input.Search)
		search.Map["transaction_id"] = fmt.Sprintf("%%%s%%", *input.Search)
	}

	query := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			filter,
			search,
		},
		WJoinOperator: xiao.And,
		OMap: xiao.SQLMap{
			Map: map[string]interface{}{
				"transactions.created_at": input.Sort.String(),
			},
		},
		Pagination: xiao.Pagination{
			Limit:  *input.Limit,
			Offset: offset,
		},
	}

	ts, err := repository.Transaction().FindAllByMap(context.Background(), query, true)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	var pagination *garage.Pagination

	if input.Paginate != nil && *input.Paginate {
		// get total count
		total, err := repository.Transaction().CountByMap(context.Background(), query)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, err
		}
		pagination = &garage.Pagination{
			Page:  *input.Page,
			Limit: *input.Limit,
			Pages: int(math.Ceil(float64(total) / float64(*input.Limit))),
		}
	}

	return ts, pagination, nil
}
