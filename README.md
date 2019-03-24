# qiita

Go client library for [qiita API v2](https://qiita.com/api/v2/docs).

[![CircleCI](https://circleci.com/gh/muiscript/qiita/tree/master.svg?style=svg)](https://circleci.com/gh/muiscript/qiita/tree/master)
[![codecov](https://codecov.io/gh/muiscript/qiita/branch/master/graph/badge.svg)](https://codecov.io/gh/muiscript/qiita)

## usage

```go
logger := log.New(os.Stdout, "[LOG]", log.LstdFlags)
qiita := qiita.New("<YOUR_ACCESS_TOKEN>", logger)

ctx := context.Background()

// get user
user, err := qiita.GetUser(ctx, "muiscript")

// get item
item, err := qiita.GetItem(ctx, "b4ca1773580317e7112e")
```

## API list

#### apis available for unauthorized/authorized users

| Done | Endpoint | Method Signature |
| --- | --- | --- |
| :heavy_check_mark: | `GET` - `/users` | `GetUsers(ctx context.Context, page int, perPage int)` |
| :heavy_check_mark: | `GET` - `/users/:user_id` | `GetUser(ctx context.Context, userID string)` |
| :heavy_check_mark: | `GET` - `/users/:user_id/followees` | `GetUserFollowees(ctx context.Context, userID string, page int, perPage int)` |
| :heavy_check_mark: | `GET` - `/users/:user_id/followers` | `GetUserFollowers(ctx context.Context, userID string, page int, perPage int)`|
|  | `GET` - `/users/:user_id/items` | `GetUserItems(ctx context.Context, userID string)` |
|  | `GET` - `/users/:user_id/stocks` | `GetUserStocks(ctx context.Context, userID string)` |
|  | `GET` - `/users/:user_id/following_tags` | `GetUserFollowingTags(ctx context.Context, userID string)` |
|  | `GET` - `/items` | `GetItems(ctx context.Context)` |
| :heavy_check_mark: | `GET` - `/items/:item_id` | `GetItem(ctx context.Context, itemID string)` |
|  | `GET` - `/items/:item_id/stockers` | `GetItemStockers(ctx context.Context, itemID string)` |
|  | `GET` - `/items/:item_id/comments` | `GetItemComments(ctx context.Context, itemID string)` |
|  | `GET` - `/tags` | `GetTags(ctx context.Context)` |
|  | `GET` - `/tags/:tag_id` | `GetTag(ctx context.Context, tagID string)` |
|  | `GET` - `/tags/:tag_id/items` | `GetTagItems(ctx context.Context, tagID string)` |
|  | `GET` - `/comments/:comment_id` | `GetComment(ctx context.Context, commentID string)` |

#### apis only available for authorized users

| Done | Endpoint | Method Signature |
| --- | --- | --- |
| :heavy_check_mark: | `GET` - `/users/:user_id/following` | `IsFollowingUser(ctx context.Context, userID string)` |
|  | `PUT` - `/users/:user_id/following` | `FollowUser(ctx context.Context, userID string)` |
|  | `DELETE` - `/users/:user_id/following` | `UnfollowUser(ctx context.Context, userID string)` |
|  | `GET` - `/authenticated_user` | `GetAuthenticatedUser(ctx context.Context)` |
|  | `GET` - `/authenticated_user/items` | `GetAuthenticatedUserItems(ctx context.Context)` |
|  | `POST` - `/items` | `CreateItem(ctx context.Context, title, body string, private, tweet bool)` |
|  | `PATCH` - `/items/:item_id` | `UpdateItem(ctx context.Context, itemID string, title, body string, private, tweet bool)` |
|  | `DELETE` - `/items/:item_id` | `DeleteItem(ctx context.Context, itemID string)` |
|  | `GET` - `/items/:item_id/stock` | `IsStockedItem(ctx context.Context, itemID string)` |
|  | `PUT` - `/items/:item_id/stock` | `StockItem(ctx context.Context, itemID string)` |
|  | `DELETE` - `/items/:item_id/stock` | `UnstockItem(ctx context.Context, itemID string)` |
|  | `GET` - `/tags/:tag_id/following` | `IsFollowingTag(ctx context.Context, tagID string)` |
|  | `PUT` - `/tags/:tag_id/following` | `FollowTag(ctx context.Context, tagID string)` |
|  | `DELETE` - `/tags/:tag_id/following` | `UnfollowTag(ctx context.Context, tagID string)` |
|  | `POST` - `/items/:item_id/comments` | `CreateItemComment(ctx context.Context, itemID string, body string)` |
|  | `PATCH` - `/comments/:comment_id` | `UpdateComment(ctx context.Context, commentID string, body string)` |
|  | `DELETE` - `/comments/:comment_id` | `DeleteComment(ctx context.Context, commentID string)` |
|  | `PUT` - `/comments/:comment_id/thank` | `ThankComment(ctx context.Context, commentID string)` |
|  | `DELETE` - `/comments/:comment_id/thank` | `UnthankComment(ctx context.Context, commentID string)` |
