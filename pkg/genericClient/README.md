# customclient

customclient is a collection of clients used by access to different API, write in GO 

The packages are kept separate to keep it as small and concise as possible.

* Note: We have used package-oriented-design to create this library, you can know more about this
  [link](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html)


## Packages


| PKG           | expected kubernetes compatibility                |
|---------------|--------------------------------------------------|
| genericClient | a client library in Go to access any API         |


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
