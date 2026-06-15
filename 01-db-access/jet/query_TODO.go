package jetspike

// TODO: お題のクエリを書く。
// jet は型付きSQLを Go で組み立て、結果を入れ子オブジェクトへ自動マッピングできる。
//
// イメージ(生成された .gen の table/model を import して使う):
//   stmt := SELECT(
//       Users.AllColumns, Posts.AllColumns, Comments.AllColumns,
//   ).FROM(
//       Users.
//           LEFT_JOIN(Posts, Posts.UserID.EQ(Users.ID).AND(Posts.DeletedAt.IS_NULL())).
//           LEFT_JOIN(Comments, Comments.PostID.EQ(Posts.ID).AND(Comments.DeletedAt.IS_NULL())),
//   ).WHERE(
//       Users.ID.EQ(UUID(userID)).AND(Users.DeletedAt.IS_NULL()),
//   )
//
//   var dest struct {
//       model.Users
//       Posts []struct {
//           model.Posts
//           Comments []model.Comments
//       }
//   }
//   err := stmt.Query(db, &dest)   // ← 入れ子に自動マッピングされる点に注目
//
// TODO: 実際に生成物を import して上を完成させ、書き心地を sqlc/ent と比べる。
