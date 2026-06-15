# Atlas 環境設定(宣言的モードの最小例)
env "local" {
  url     = "postgres://spike:spike@localhost:5432/spike?sslmode=disable&search_path=public"
  # 望む最終状態(desired schema)。ここでは共通スキーマSQLを参照。
  src     = "file://../../schema/schema.sql"
  dev     = "docker://postgres/16/dev?search_path=public"   # 差分計算用の使い捨てDB
  migration {
    dir = "file://migrations"
  }
}
