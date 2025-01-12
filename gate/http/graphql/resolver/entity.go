package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"cendit.io/auth/repository"
	"cendit.io/garage/function"
	"cendit.io/garage/xiao"
	graphql1 "cendit.io/gate/http/graphql"
	"cendit.io/gate/http/graphql/model"
)

// FindAddressByUserID is the resolver for the findAddressByUserID field.
func (r *entityResolver) FindAddressByUserID(ctx context.Context, userID string) (*model.Address, error) {
	if userID != "" {

		address, err := repository.Address().FindByMap(context.Background(), xiao.SQLMaps{
			WMaps: []xiao.SQLMap{
				{
					Map: map[string]interface{}{
						"user_id": userID,
					},
					JoinOperator:       xiao.And,
					ComparisonOperator: xiao.Equal,
				},
			},
		}, true)
		if err == nil {
			a := model.Address{}

			function.Parse(address, &a)

			return &a, nil
		}
	}
	return nil, fmt.Errorf("address not found")
}

// FindSecuritySettingByUserID is the resolver for the findSecuritySettingByUserID field.
func (r *entityResolver) FindSecuritySettingByUserID(ctx context.Context, userID string) (*model.SecuritySetting, error) {
	if userID != "" {

		setting, err := repository.SecuritySetting().FindByMap(context.Background(), xiao.SQLMaps{
			WMaps: []xiao.SQLMap{
				{
					Map: map[string]interface{}{
						"user_id": userID,
					},
					JoinOperator:       xiao.And,
					ComparisonOperator: xiao.Equal,
				},
			},
		}, true)
		if err == nil {
			s := model.SecuritySetting{}

			function.Parse(setting, &s)

			return &s, nil
		}
	}
	return nil, fmt.Errorf("security setting not found")
}

// FindUserByID is the resolver for the findUserByID field.
func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	if id != "" {

		user, err := repository.User().FindByMap(context.Background(), xiao.SQLMaps{
			WMaps: []xiao.SQLMap{
				{
					Map: map[string]interface{}{
						"id": id,
					},
					JoinOperator:       xiao.And,
					ComparisonOperator: xiao.Equal,
				},
			},
		}, true)
		if err == nil {
			u := model.User{}

			function.Parse(user, &u)

			return &u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// FindWalletByUserID is the resolver for the findWalletByUserID field.
func (r *entityResolver) FindWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	if userID != "" {

		wallet, err := repository.Wallet().FindByMap(context.Background(), xiao.SQLMaps{
			WMaps: []xiao.SQLMap{
				{
					Map: map[string]interface{}{
						"user_id": userID,
					},
					JoinOperator:       xiao.And,
					ComparisonOperator: xiao.Equal,
				},
			},
		}, true)
		if err == nil {
			w := model.Wallet{}

			function.Parse(wallet, &w)

			return &w, nil
		}
	}
	return nil, fmt.Errorf("wallet not found")
}

// Entity returns graphql1.EntityResolver implementation.
func (r *Resolver) Entity() graphql1.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
