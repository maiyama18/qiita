# qiita

Go client library for [qiita API v2](https://qiita.com/api/v2/docs).

## usage

```go
logger := log.New(os.Stdout, "[LOG]", log.LstdFlags)
qiita := qiita.New(logger)

ctx := context.Background()

// get user
user, err := qiita.GetUser(ctx, "muiscript")

// get item
item, err := qiita.GetItem(ctx, "b4ca1773580317e7112e")
```

## API list

#### apis for authorization

|  | Endpoint | Method Signature |
| --- | --- | --- |
| [ ] | `POST` - `/api/v2/access_tokens` | |
| [ ] | `DELETE` - `/api/v2/access_tokens/:access_token` | |

#### apis available for unauthorized/authorized users

|  | Endpoint | Method Signature |
| --- | --- | --- |
| [ ] | `GET` - `/users` | |
| [x] | `GET` - `/users/:user_id` | `GetUser(ctx context.Context, userID string) (*User, error)` |
| [ ] | `GET` - `/users/:user_id/followees` | |
| [ ] | `GET` - `/users/:user_id/followers` | |
| [ ] | `GET` - `/users/:user_id/items` | |
| [ ] | `GET` - `/users/:user_id/stocks` | |
| [ ] | `GET` - `/users/:user_id/following_tags` | |
| [ ] | `GET` - `/items` | |
| [x] | `GET` - `/items/:item_id` | `GetItem(ctx context.Context, itemID string) (*Post, error)` |
| [ ] | `GET` - `/items/:item_id/stockers` | |
| [ ] | `GET` - `/items/:item_id/comments` | |
| [ ] | `GET` - `/tags` | |
| [ ] | `GET` - `/tags/:tag_id` | |
| [ ] | `GET` - `/tags/:tag_id/items` | |
| [ ] | `GET` - `/comments/:comment_id` | |

#### apis available for unauthorized/authorized users

|  | Endpoint | Method Signature |
| --- | --- | --- |
| [ ] | `GET` - `/users/:user_id/following` | |
| [ ] | `DELETE` - `/users/:user_id/following` | |
| [ ] | `PUT` - `/users/:user_id/following` | |
| [ ] | `GET` - `/authenticated_user` | |
| [ ] | `GET` - `/authenticated_user/items` | |
| [ ] | `POST` - `/items` | |
| [ ] | `DELETE` - `/items/:item_id` | |
| [ ] | `PATCH` - `/items/:item_id` | |
| [ ] | `DELETE` - `/items/:item_id` | |
| [ ] | `GET` - `/items/:item_id/stock` | |
| [ ] | `PUT` - `/items/:item_id/stock` | |
| [ ] | `DELETE` - `/items/:item_id/stock` | |
| [ ] | `GET` - `/tags/:tag_id/following` | |
| [ ] | `PUT` - `/tags/:tag_id/following` | |
| [ ] | `DELETE` - `/tags/:tag_id/following` | |
| [ ] | `POST` - `/items/:item_id/comments` | |
| [ ] | `PATCH` - `/comments/:comment_id` | |
| [ ] | `DELETE` - `/comments/:comment_id` | |
| [ ] | `PUT` - `/comments/:comment_id/thank` | |
| [ ] | `DELETE` - `/comments/:comment_id/thank` | |
