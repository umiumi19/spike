# ent

## セットアップ
```sh
go get entgo.io/ent
# エンティティの雛形を生成(User / Post / Comment)
go run -mod=mod entgo.io/ent/cmd/ent new User Post Comment
# → ent/schema/{user.go, post.go, comment.go} ができる
```
その後 `ent/schema/*.go` に **フィールドと edge(関係)** を定義する(TODO)。
- User: name, timestamps, deleted_at / edge: posts(1対多), comments
- Post: title, body, ... / edge: author(User), comments
- Comment: body, ... / edge: post, author

```sh
go generate ./ent      # 定義から型付きクライアントを生成
```

## お題のクエリ(TODO)
ent は SQL を書かず、関係を `.With...()` でたどる:
```go
// u, err := client.User.Query().
//     Where(user.ID(userID), user.DeletedAtIsNil()).
//     WithPosts(func(pq *ent.PostQuery) {
//         pq.Where(post.DeletedAtIsNil()).
//            WithComments(func(cq *ent.CommentQuery) {
//                cq.Where(comment.DeletedAtIsNil())
//            })
//     }).
//     Only(ctx)
// u.Edges.Posts[0].Edges.Comments[0] のように入れ子で型付きアクセスできる
```

## マイグレーション(ent統合)
ent はスキーマからマイグレーションも生成できる:
```sh
# 例: 自動マイグレーションで開発DBに反映(本番は差分マイグレーション推奨)
# client.Schema.Create(ctx) を main で呼ぶ、もしくは ent の migrate diff を使う
```
