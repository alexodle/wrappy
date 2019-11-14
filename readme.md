# Wrappy

Interface and wrapper generator for GoLang

## Problem

In Go, if vendor code is written without interfaces, there is no way to generate mocks. `wrappy` fixes that problem by
recursively crawling source code and generating interfaces and struct wrappers for all (or a specific subset) of the public
structs within that code.

## Example

Let's use the example of the following notional "BlobClient", which would be used to download files from a storage service:

**vendor/some/vendor/repo/storage_client.go**
```
...

type StorageClient struct {}

func (s *StorageClient) GetContainerClient(container string) ContainerClient {
	...
}

type ContainerClient struct {}

func (c *ContainerClient) DownloadBlob(path string) (Response, error) {
	...
}

type Response struct {}

func (r *Response) Body() io.ReadCloser {
	...
}
```

Notice how this code uses and references structs directly, making it impossible for us to mock out this code for testing.
In addition, we can't simply write/generate our own interface for `StorageClient`, as the `GetContainerClient` method returns
a concrete `ContainerClient` struct, which we'd also want to mock. What's needed is a set of interfaces + a set of wrappers
that implement slightly altered versions of those methods, which return other generated interfaces instead of structs.

Let's run wrappy on this code and see what happens:

`go run main/main.go vendor/some/vendor/repo/ generated/some/vendor/repo/`

**generated/some/vendor/repo/storage_client.go**
```
type StorageClient interface {
	GetImpl() *orig_client.Client
	GetContainerClient(container string) ContainerClient // Returns a ContainerClient interface
}

type ContainerClient interface {
	GetImpl() *orig_client.ContainerClient
	DownloadBlob(path string) (Response, error) // Returns a Response interface
}

type Response interface {
	GetImpl() *orig_client.Response
	Body() io.ReadCloser
}

func NewClient(impl *orig_client.Client) StorageClient {
	return &clientWrapper{impl: impl}
}

func NewContainerClient(impl *orig_tmp.ContainerClient) ContainerClient {
	return &containerClientWrapper{impl: impl}
}

func NewResponse(impl *orig_tmp.Response) Response {
	return &responseWrapper{impl: impl}
}

... // Implementation details omitted

```

Now we can implement our code using these generated interfaces, and write tests using generated mocks. Note that mock
generation is outside the scope of this project. There are plenty of good mock generators out there. Mockery
(https://github.com/vektra/mockery), for example, could generate mocks here by running:

`mockery -all -dir generated/some/vendor/repo/`
