# \TransactionApi

All URIs are relative to *https://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthTransactionGet**](TransactionApi.md#AuthTransactionGet) | **Get** /auth/transaction | Get transactions
[**AuthTransactionPost**](TransactionApi.md#AuthTransactionPost) | **Post** /auth/transaction | Create an transaction


# **AuthTransactionGet**
> TransactionTransaction AuthTransactionGet(ctx, name, optional)
Get transactions

Search transactions based on query parameters

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Account user name | 
 **optional** | ***AuthTransactionGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthTransactionGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **type_** | **optional.Int32**| Outgoing(0)/Incoming(1)/Both(2) transactions for the account | 
 **start** | **optional.String**| The start of the query time range, format like this: 2006-01-02 15:04:05 | 
 **end** | **optional.String**| The end of the query time range, format like this: 2006-01-02 15:04:05 | 

### Return type

[**TransactionTransaction**](transaction.Transaction.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AuthTransactionPost**
> TransactionTransaction AuthTransactionPost(ctx, transaction)
Create an transaction

Transfer money from one account to another

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **transaction** | [**TransactionTransaction**](TransactionTransaction.md)| JSON structure | 

### Return type

[**TransactionTransaction**](transaction.Transaction.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

