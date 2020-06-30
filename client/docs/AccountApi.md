# \AccountApi

All URIs are relative to *https://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthAccountGet**](AccountApi.md#AuthAccountGet) | **Get** /auth/account | List all accounts
[**AuthAccountIdDelete**](AccountApi.md#AuthAccountIdDelete) | **Delete** /auth/account/{id} | Delete an account
[**AuthAccountIdGet**](AccountApi.md#AuthAccountIdGet) | **Get** /auth/account/{id} | Show an account
[**AuthAccountIdPut**](AccountApi.md#AuthAccountIdPut) | **Put** /auth/account/{id} | Update an account
[**AuthPost**](AccountApi.md#AuthPost) | **Post** /auth | Login an account
[**RegPost**](AccountApi.md#RegPost) | **Post** /reg | Create an account


# **AuthAccountGet**
> AccountAccount AuthAccountGet(ctx, optional)
List all accounts

List all the accounts

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AuthAccountGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthAccountGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **line** | **optional.Int32**| List result number no more than the value of &#39;line&#39; | 

### Return type

[**AccountAccount**](account.Account.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AuthAccountIdDelete**
> AccountAccount AuthAccountIdDelete(ctx, id)
Delete an account

Delete by account ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **int32**| Account ID | 

### Return type

[**AccountAccount**](account.Account.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AuthAccountIdGet**
> AccountAccount AuthAccountIdGet(ctx, id)
Show an account

Get by account ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **int32**| Account ID | 

### Return type

[**AccountAccount**](account.Account.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AuthAccountIdPut**
> AccountAccount AuthAccountIdPut(ctx, id, fields)
Update an account

Update by account ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **int32**| Account ID | 
  **fields** | [**AccountAccount**](AccountAccount.md)| Update fields/values as JSON | 

### Return type

[**AccountAccount**](account.Account.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AuthPost**
> AccountLoginRsp AuthPost(ctx, req)
Login an account

Login using username and password

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **req** | [**AccountLoginReq**](AccountLoginReq.md)| account.LoginReq struct JSON | 

### Return type

[**AccountLoginRsp**](account.LoginRsp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RegPost**
> AccountAccount RegPost(ctx, account)
Create an account

Create an account by JSON format parameters

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **account** | [**AccountAccount**](AccountAccount.md)| account.Account struct JSON | 

### Return type

[**AccountAccount**](account.Account.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

