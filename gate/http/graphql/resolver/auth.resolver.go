package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"fmt"

	"cendit.io/auth/logic"
	"cendit.io/auth/schema"
	"cendit.io/garage/function"
	"cendit.io/garage/logger"
	"cendit.io/garage/primer/constant"
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/primer/typing"
	"cendit.io/gate/http/graphql/exception"
	"cendit.io/gate/http/graphql/model"
	"cendit.io/gate/http/rest/interceptor"
)

// SignupStepOne is the resolver for the signupStepOne field.
func (r *mutationResolver) SignupStepOne(ctx context.Context, input model.StepOneInput) (model.RespondWithUser, error) {
	logger.GetLogger().Debug(fmt.Sprintf(`START :: [%v] :: mutationResolver.SignupStepOne with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	act := schema.Activity{
		ID:          function.GenerateUUID(),
		Resolver:    "mutationResolver.SignupStepOne",
		Payload:     function.Stringify(input),
		Description: "made an action for step one of account creation",
		Error:       "",
		Status:      "",
	}
	act.Date()

	user, err := logic.SignupStepOne(input)
	if err != nil {
		if errors.As(err, &exception.Error{}) {
			intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
			intrusion.Status = err.(exception.Error).Status
			return exception.MakeSubgraphError(err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})), function.Stringify(act)), nil
		}

		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = constant.CodeISE
		return exception.MakeSubgraphError(`Something went wrong while creating your account! Please try again.`, constant.CodeISE, fmt.Sprintf(`[%v] :: {message} :: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error()), function.Stringify(act)), nil
	}

	// log activity
	act.Role = user.Role
	act.By = user.ID
	act.Status = enum.Success
	if err := act.Insert(); err != nil {
		logger.GetLogger().Debug(`ACTIVITY LOG :: SAVE ERROR :: ` + err.Error())
	}

	logger.GetLogger().Debug(fmt.Sprintf(`END :: [%v] :: mutationResolver.SignupStepOne with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
	intrusion.Session = typing.Session{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	var u model.User
	function.Parse(user, &u)
	return &model.ResponseWithUser{
		Message: "Account created successfully. Please proceed to the next step.",
		Data:    &u,
	}, nil
}

// SignupStepTwo is the resolver for the signupStepTwo field.
func (r *mutationResolver) SignupStepTwo(ctx context.Context, input model.StepTwoInput) (model.RespondWithUser, error) {
	logger.GetLogger().Debug(fmt.Sprintf(`START :: [%v] :: mutationResolver.SignupStepTwo with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	act := schema.Activity{
		ID:          function.GenerateUUID(),
		Resolver:    "mutationResolver.SignupStepTwo",
		Payload:     function.Stringify(input),
		Description: "made an action for step two of account creation",
		Error:       "",
		Status:      "",
	}
	act.Date()

	// access control
	_, authObj, err := interceptor.Authorize(ctx)
	// we only expect a exception.Error here
	if err != nil && errors.As(err, &exception.Error{}) {
		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = err.(exception.Error).Status
		return nil, exception.MakeGraphQLError(ctx, err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})))
	}

	user, err := logic.SignupStepTwo(input, authObj.(schema.User))
	if err != nil {
		if errors.As(err, &exception.Error{}) {
			intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
			intrusion.Status = err.(exception.Error).Status
			return exception.MakeSubgraphError(err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})), function.Stringify(act)), nil
		}

		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = constant.CodeISE

		return exception.MakeSubgraphError(`Something went wrong while saving your personal information! Please try again.`, constant.CodeISE, fmt.Sprintf(`[%v] :: {message} :: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error()), function.Stringify(act)), nil
	}

	// log activity
	act.Role = user.Role
	act.By = user.ID
	act.Status = enum.Success
	if err := act.Insert(); err != nil {
		logger.GetLogger().Debug(`ACTIVITY LOG :: SAVE ERROR :: ` + err.Error())
	}

	logger.GetLogger().Debug(fmt.Sprintf(`END :: [%v] :: mutationResolver.SignupStepTwo with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	var u model.User
	function.Parse(user, &u)
	return &model.ResponseWithUser{
		Message: "Personal information saved successfully. Please proceed to the next step.",
		Data:    &u,
	}, nil
}

// SignupStepThree is the resolver for the signupStepThree field.
func (r *mutationResolver) SignupStepThree(ctx context.Context, input model.StepThreeInput) (model.RespondWithAddress, error) {
	logger.GetLogger().Debug(fmt.Sprintf(`START :: [%v] :: mutationResolver.signupStepThree with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	act := schema.Activity{
		ID:          function.GenerateUUID(),
		Resolver:    "mutationResolver.SignupStepThree",
		Payload:     function.Stringify(input),
		Description: "made an action for step three of account creation",
		Error:       "",
		Status:      "",
	}
	act.Date()

	// access control
	auth, authObj, err := interceptor.Authorize(ctx)
	// we only expect a exception.Error here
	if err != nil && errors.As(err, &exception.Error{}) {
		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = err.(exception.Error).Status
		return nil, exception.MakeGraphQLError(ctx, err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})))
	}

	address, err := logic.SignupStepThree(input, authObj.(schema.User))
	if err != nil {
		if errors.As(err, &exception.Error{}) {
			intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
			intrusion.Status = err.(exception.Error).Status
			return exception.MakeSubgraphError(err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})), function.Stringify(act)), nil
		}

		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = constant.CodeISE

		return exception.MakeSubgraphError(`Something went wrong while saving your address! Please try again.`, constant.CodeISE, fmt.Sprintf(`[%v] :: {message} :: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error()), function.Stringify(act)), nil
	}

	// log activity
	act.Role = auth.Role
	act.By = auth.ID
	act.Status = enum.Success
	if err := act.Insert(); err != nil {
		logger.GetLogger().Debug(`ACTIVITY LOG :: SAVE ERROR :: ` + err.Error())
	}

	logger.GetLogger().Debug(fmt.Sprintf(`END :: [%v] :: mSutationResolver.signupStepThree with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	var add model.Address
	function.Parse(address, &add)

	return &model.ResponseWithAddress{
		Message: "Address saved successfully. Please proceed to the next step.",
		Data:    &add,
	}, nil
}

// SetPasscode is the resolver for the setPasscode field.
func (r *mutationResolver) SetPasscode(ctx context.Context, input model.SetPasscodeInput) (model.RespondWithUser, error) {
	logger.GetLogger().Debug(fmt.Sprintf(`START :: [%v] :: mutationResolver.SetPasscode with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	act := schema.Activity{
		ID:          function.GenerateUUID(),
		Resolver:    "mutationResolver.SetPasscode",
		Payload:     function.Stringify(input),
		Description: "made an action to set passcode",
		Error:       "",
		Status:      "",
	}
	act.Date()

	// access control
	_, authObj, err := interceptor.Authorize(ctx)
	// we only expect a exception.Error here
	if err != nil && errors.As(err, &exception.Error{}) {
		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = err.(exception.Error).Status
		return nil, exception.MakeGraphQLError(ctx, err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})))
	}

	user, err := logic.SetPasscode(input, authObj.(schema.User))
	if err != nil {
		if errors.As(err, &exception.Error{}) {
			intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
			intrusion.Status = err.(exception.Error).Status
			return exception.MakeSubgraphError(err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})), function.Stringify(act)), nil
		}

		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = constant.CodeISE

		return exception.MakeSubgraphError(`Something went wrong while setting your passcode! Please try again.`, constant.CodeISE, fmt.Sprintf(`[%v] :: {message} :: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error()), function.Stringify(act)), nil
	}

	// log activity
	act.Role = user.Role
	act.By = user.ID
	act.Status = enum.Success
	if err := act.Insert(); err != nil {
		logger.GetLogger().Debug(`ACTIVITY LOG :: SAVE ERROR :: ` + err.Error())
	}

	logger.GetLogger().Debug(fmt.Sprintf(`END :: [%v] :: mutationResolver.SetPasscode with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	var u model.User
	function.Parse(user, &u)
	return &model.ResponseWithUser{
		Message: "Passcode set successfully.",
		Data:    &u,
	}, nil
}

// StartSession is the resolver for the startSession field.
func (r *mutationResolver) StartSession(ctx context.Context, input model.StartSessionInput) (model.RespondWithUser, error) {
	logger.GetLogger().Debug(fmt.Sprintf(`START :: [%v] :: mutationResolver.StartSession with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	act := schema.Activity{
		ID:          function.GenerateUUID(),
		Resolver:    "mutationResolver.StartSession",
		Payload:     function.Stringify(input),
		Description: "made an action to start session",
		Error:       "",
		Status:      "",
	}
	act.Date()

	user, err := logic.StartSession(input)
	if err != nil {
		if errors.As(err, &exception.Error{}) {
			intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
			intrusion.Status = err.(exception.Error).Status
			return exception.MakeSubgraphError(err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{})), function.Stringify(act)), nil
		}

		intrusion := ctx.Value(typing.CtxIntrusionKey{}).(*typing.Intrusion)
		intrusion.Status = constant.CodeISE

		return exception.MakeSubgraphError(`Something went wrong while creating your session! Please try again.`, constant.CodeISE, fmt.Sprintf(`[%v] :: {message} :: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error()), function.Stringify(act)), nil
	}

	// log activity
	act.Role = user.Role
	act.By = user.ID
	act.Status = enum.Success
	if err := act.Insert(); err != nil {
		logger.GetLogger().Debug(`ACTIVITY LOG :: SAVE ERROR :: ` + err.Error())
	}

	logger.GetLogger().Debug(fmt.Sprintf(`END :: [%v] :: mutationResolver.StartSession with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	var u model.User
	function.Parse(user, &u)
	return &model.ResponseWithUser{
		Message: "Welcome back.",
		Data:    &u,
	}, nil
}

// Signin is the resolver for the signin field.
func (r *queryResolver) Signin(ctx context.Context, input model.SigninInput) (model.RespondWithUser, error) {
	logger.GetLogger().Debug(fmt.Sprintf(`START :: [%v] :: queryResolver.Siginin with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	act := schema.Activity{
		ID:          function.GenerateUUID(),
		Resolver:    "queryResolver.Signin",
		Payload:     function.Stringify(input),
		Description: "made an action to sign in",
		Error:       "",
		Status:      "",
	}
	act.Date()

	user, err := logic.Signin(input)
	if err != nil {
		if errors.As(err, &exception.Error{}) {
			return exception.MakeSubgraphError(err.(exception.Error).Message, err.(exception.Error).Status, fmt.Sprintf(`[%v] :: {message}`, ctx.Value(typing.CtxTraceKey{}))), nil
		}

		return exception.MakeSubgraphError(`Something went wrong while signing in! Please try again.`, constant.CodeISE, fmt.Sprintf(`[%v] :: {message} :: %s`, ctx.Value(typing.CtxTraceKey{}), err.Error())), nil
	}

	// log activity
	act.Role = user.Role
	act.By = user.ID
	act.Status = enum.Success
	if err := act.Insert(); err != nil {
		logger.GetLogger().Debug(`ACTIVITY LOG :: SAVE ERROR :: ` + err.Error())
	}

	logger.GetLogger().Debug(fmt.Sprintf(`END :: [%v] :: queryResolver.Siginin with input: %+v`, ctx.Value(typing.CtxTraceKey{}), function.Stringify(input)))

	var u model.User
	function.Parse(user, &u)
	return &model.ResponseWithUser{
		Message: "Welcome back.",
		Data:    &u,
	}, nil
}
