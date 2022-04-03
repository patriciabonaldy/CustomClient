# Form3PKG

Form3 is a collection of clients used by access to different API, write in GO 

The packages are kept separate to keep it as small and concise as possible.

## Packages


| PKG           | expected kubernetes compatibility                |
|---------------|--------------------------------------------------|
| accountClient | a client library in Go to access our account API |
| genericClient | a client library in Go to access any API         |


### accountClient pkg

  it's a client library in Go to access our fake account API, it's a suitable for use in 
  another software project, it has different methods available:
~~~bash
* Delete(ctx context.Context, accountID string, version int) error
* GetByAccountID(ctx context.Context, accountID string) (RequestAccount, error)
* Save(ctx context.Context, accountData AccountData) (RequestAccount, error)
  ~~~

### How to use:

1.- Import the library in your project.

2.- create a new account client 

~~~go
ac := New(baseURL string)
~~~

Example:

 * creating a new account
~~~go
 accountData:=  AccountData{
              ID:             "b8fc26d7-ca23-4b32-a5ad-3a5e39b048be",
              OrganisationID: "b8fc26d7-ca23-4b32-a5ad-3a5e39b058be",
              Type:           "accounts",
              Attributes: &Attributes{
                  Country:    getString("GB"),
                  Name:       []string{"John"},
                  BankID:     "400302",
                  Bic:        "NWBKGB42",
                  BankIDCode: "GBDSC",
              },
          }
 
 ac.Save(context.Background(), tt.accountData)
~~~

* getting an account

~~~go
resp = ac.GetByAccountID(context.Background(), "b8fc26d7-ca23-4b32-a5ad-3a5e39b048be")

fmt.Println(resp)  
// output
RequestAccount{
  Account: AccountData{
    Attributes: &Attributes{
        AccountClassification:   nil,
        AccountMatchingOptOut:   nil,
        AccountNumber:           "41426819",
        AlternativeNames:        nil,
        BankID:                  "400300",
        BankIDCode:              "GBDSC",
        BaseCurrency:            "GBP",
        Bic:                     "NWBKGB22",
        Country:                 getString("GB"),
        Iban:                    "GB11NWBK40030041426819",
        JointAccount:            nil,
        Name:                    nil,
        SecondaryIdentification: "",
        Status:                  getString("confirmed"),
        Switched:                nil,
    },
    ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
    OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
    Type:           "accounts",
},

~~~


* deleting an account

~~~go
ac.Delete(context.Background(), "b8fc26d7-ca23-4b32-a5ad-3a5e39b048be", 0)
~~~
----------
### genericClient pkg


it's a client library in Go to access any API, it's a suitable library for use in
another software project, it has different methods available:
~~~bash
* Delete(ctx context.Context, url string, headers ...Header) error
* Get(ctx context.Context, url string) (resp *http.Response, err error)
* Post(ctx context.Context, url string, data []byte, headers ...Header) (resp *http.Response, err error)
  ~~~
### How to use:

1.- Import the library in your project.

2.- create a new client

~~~go
client:= genericClient.New()
~~~

You can add values to the header
~~~go
// Header represents Header in the request.
type Header struct {
  Key   string
  Value string
}
~~~

* Post method
~~~go
headers := []genericClient.Header{{Key: "Content-Type", Value: "application/json"}}
resp, err :=client.Post(ctx, a.baseURL, body, headers...)
~~~



--------

### Testing

~~~bash
make tests
~~~