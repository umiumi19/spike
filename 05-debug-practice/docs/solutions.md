# 解答とデバッグ手順

各演習について「実際に出る panic」「ブレークポイントで原因に辿り着く手順」「原因」「修正例」をまとめています。
**まずは自力で粘ってから**読むのがおすすめです。

> 共通のコツ：「実行とデバッグ」ビューの **ブレークポイント** セクションで `Uncaught Exceptions` にチェックを入れておくと、panic の瞬間に自動で止まり、その時点の変数・コールスタックをそのまま観察できます。原因究明が大幅に楽になります。

---

## 演習1 — `ex01_index`

**実行すると出るエラー**
```
panic: runtime error: index out of range [3] with length 3
main.lastValue(...) .../ex01_index/main.go:7
```
メッセージが「長さ3のスライスに添字3でアクセスした」と教えてくれています。

**手順**
1. `exercises/ex01_index/main.go` を開き、7行目（`return nums[len(nums)]`）の左余白をクリックしてブレークポイントを置く。
2. `F5` でデバッグ実行。`lastValue` の中で停止する。
3. **変数パネル**で `nums` を展開 → 要素は `[10 20 30]`、長さは 3。
4. **ウォッチ式**に `len(nums)` を登録 → `3`。**デバッグコンソール**で `nums[3]` と打つと、まさにこれが範囲外だと分かる。
5. 有効な添字は 0〜2 なのに、`len(nums)`（= 3）でアクセスしている、と気づく。

**原因**
スライスの添字は 0 始まりなので、末尾要素は `len-1`。`len(nums)` は常に「最後の次」を指してしまい、必ず範囲外になる。

**修正**
```go
func lastValue(nums []int) int {
	return nums[len(nums)-1]
}
```
（より丁寧にするなら、空スライス `len==0` のときの扱いも決めておく。）

---

## 演習2 — `ex02_nilptr`

**実行すると出るエラー**
```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: ...]
main.main() .../ex02_nilptr/main.go:27
```
落ちているのは 27 行目（`fmt.Printf(... u.Name, u.Age)`）。ここが**症状の出る場所**。

**手順**
1. 27行目にブレークポイントを置いて `F5`。
2. 停止したら変数パネルで `u` を見る → **`<nil>`**。nil のポインタの `u.Name` を読もうとして落ちる、と分かる。
3. では「なぜ `u` が nil なのか」。`u := findUser(users, "Carol")` の行（26行目）にブレークポイントを置き直し、再実行。
4. その行で **ステップイン（`F11`）** して `findUser` の中に入る。
5. ループを `F10` で回しながら `u.Name`（候補側）と探している `name`（`"Carol"`）を見比べる。`"Alice"`,`"Bob"` のどちらとも一致せず、ループを抜けて `return nil` に到達する様子を観察する。
6. **コールスタック**で `main` に戻れば、その nil がそのまま 27 行目に流れていることが確認できる。

**原因**
`findUser` は見つからないと `nil` を返す仕様。呼び出し側が nil チェックをせずに `u.Name` を参照したため、nil ポインタ参照で panic。

**修正**
```go
u := findUser(users, "Carol")
if u == nil {
	fmt.Println("該当ユーザーが見つかりませんでした")
	return
}
fmt.Printf("%s さんは %d 歳です\n", u.Name, u.Age)
```

---

## 演習3 — `ex03_conditional`

**実行すると出るエラー**
```
りんご: 平均単価 = 100
バナナ: 平均単価 = 100
panic: runtime error: integer divide by zero
main.averagePrice(...) .../ex03_conditional/main.go:13
```
2件は表示できてから落ちる → **特定のデータだけ**で壊れている、というのが大きなヒント。

**手順（条件付きブレークポイントの練習）**
1. 13行目（`return s.Revenue / s.Quantity`）にブレークポイントを置く。
2. ただしこのままだと毎ループ止まって面倒。赤丸を**右クリック →「ブレークポイントの編集」**を選び、条件式に `s.Quantity == 0` を入力。
3. `F5` で実行すると、**3件目（さくらんぼ）のときだけ**停止する。
4. 変数パネルで `s` を展開 → `Quantity` が `0`。整数の 0 除算が panic の正体だと確定する。

**原因**
`averagePrice` が割る数（`Quantity`）のゼロチェックをしていない。`Quantity == 0` のレコードで整数の 0 除算が起きる。

**修正（呼び出し側で扱いを分けられるよう ok を返す例）**
```go
func averagePrice(s Sale) (int, bool) {
	if s.Quantity == 0 {
		return 0, false
	}
	return s.Revenue / s.Quantity, true
}

// 呼び出し側
for _, s := range sales {
	if avg, ok := averagePrice(s); ok {
		fmt.Printf("%s: 平均単価 = %d\n", s.Product, avg)
	} else {
		fmt.Printf("%s: 数量が0のため計算をスキップ\n", s.Product)
	}
}
```

---

## 演習4 — `ex04_callstack`

**実行すると出るエラー**
```
panic: runtime error: index out of range [-1]
main.priceOf(...) .../ex04_callstack/main.go:18
```
添字が **`-1`** という点が決定的なヒント。「見つからなかった」を表す番兵値がそのまま添字に使われている匂いがします。

**手順（ステップイン／アウトとコールスタックの練習）**
1. 18行目（`return catalog[idx]`）にブレークポイントを置いて `F5`。
2. 停止したら `idx` を確認。問題のループ回では **`idx == -1`**。
3. その `idx` は `indexOf(items, name)` の戻り値。`priceOf` の `idx := indexOf(...)` 行で **ステップイン（`F11`）** して `indexOf` に入る。
4. `indexOf` 内を `F10` で進めると、`name`（このとき `"マグカップ"`）が `items`（`ペン/ノート/ランプ`）のどれとも一致せず、最後に `return -1` に到達するのが見える。
5. **ステップアウト（`Shift+F11`）** で `priceOf` に戻り、`catalog[-1]` で落ちることを確認。
6. さらに **コールスタック**で `main` の階層をクリックすると、その時の `name` 変数が見え、「`order` の3番目 `"マグカップ"` が引き金」だと特定できる。

**原因**
`indexOf` は見つからないと番兵値 `-1` を返す。`priceOf` がその `-1` をチェックせず添字に使ったため、`catalog[-1]` で範囲外。要は演習2と同じ「番兵値の未チェック」が、関数をまたいで起きている。

**修正**
```go
func priceOf(catalog []int, items []string, name string) (int, bool) {
	idx := indexOf(items, name)
	if idx < 0 {
		return 0, false
	}
	return catalog[idx], true
}

// 呼び出し側
for _, name := range order {
	if p, ok := priceOf(catalog, items, name); ok {
		total += p
	} else {
		fmt.Printf("カタログにない商品: %s\n", name)
	}
}
```

---

## 演習5 — `ex05_assertion`

**実行すると出るエラー**
```
panic: interface conversion: interface {} is string, not int
main.sumNumbers(...) .../ex05_assertion/main.go:9
```
「string が入っていたのに int として取り出そうとした」と明言されています。

**手順**
1. 9行目（`total += v.(int)`）にブレークポイントを置いて `F5`。
2. ループのたびに止まるので、`F5`（続行）で進めながら**変数パネルで `v` を展開**し、各回の動的な型を見る。`1`,`2` は int、3回目で `v` が文字列 `"3"` になっているのが分かる。
3. （応用）条件付きブレークポイントは型では書きづらいので、ここはステップ実行か、**デバッグコンソール**で各回 `v` を評価して観察するのが手軽。

**原因**
データに数値でない値（文字列 `"3"`）が混ざっているのに、`v.(int)` という**チェックなしの型アサーション**を使っている。型が違うと即 panic する。

**修正（カンマ ok 形式の型アサーションで安全に）**
```go
func sumNumbers(values []interface{}) int {
	total := 0
	for _, v := range values {
		if n, ok := v.(int); ok {
			total += n
		} else {
			fmt.Printf("数値でない値をスキップ: %v (%T)\n", v, v)
		}
	}
	return total
}
```
`x.(T)` は失敗すると panic、`n, ok := x.(T)` は失敗しても `ok == false` になるだけで安全、という違いがポイント。

---

## 演習6 — `ex06_logpoint`

**実行すると出る結果（panic ではない）**
```
最終残高: 1300
```
手計算では `1000 -200 +500 -100 -50 +300 = 1450` のはず。**落ちないが答えが違う**タイプのバグで、ここがログポイントの出番です。

**手順（ログポイントの練習）**
1. `exercises/ex06_logpoint/main.go` の `balance = start + d` の行（9行目）の赤丸を**右クリック →「ログポイントの追加」**。
2. メッセージに `d={d}, balance={balance}` と入力（`{}` の中は式として評価される）。
3. `F5` で実行。**一度も止まらず**、デバッグコンソールに各ループの値が流れる：
   ```
   d=-200, balance=800
   d=500, balance=1500
   d=-100, balance=900
   d=-50, balance=950
   d=300, balance=1300
   ```
4. `balance` が累積しておらず、毎回 `start(1000) + その回の d` になっているのが一目で分かる。最後の `d=300` の結果 `1300` がそのまま残って返っている。

**原因**
`balance = start + d` は、毎回 `start` から計算し直して**上書き**している。前回までの結果を引き継がないので、実質「`start` + 最後の取引」しか反映されない。

**修正**
```go
for _, d := range deltas {
	balance += d
}
```
`balance += d`（= `balance = balance + d`）にすれば、前回の残高に積み上がる。修正後は `1450`、`balance_test.go` も通る。

**ついでに試したテクニックの答え合わせ**
- **debug test**：`balance_test.go` の `TestFinalBalance` の上に出る「debug test」を押すと、修正前は `finalBalance() = 1300, want 1450` で停止・失敗。`finalBalance` 内にブレークポイントを置けば、テスト経由でそのまま中を追える。
- **Set Value**：停止中に `balance` を `9999` に書き換えて続行すると、出力もそれに応じて変わる。仮説の「もしこの値だったら?」をコードを直さず試せる。

---

## 演習7 — `ex07_deadlock`

**実行すると出るエラー**
```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
	.../ex07_deadlock/main.go:17 +0x...
```
ポイントは2つ。`all goroutines are asleep`（全 goroutine が待ち状態 ＝ 誰も先に進めない）と、`goroutine 1 [chan receive]`（main がチャネル受信で止まっている）。

**手順（ダンプと goroutine の読み方）**
1. まず普通に実行し、上記の **goroutine ダンプ**を読む。`main.main()` が `main.go:17`（`sum += <-results`）で `chan receive` のまま固まっている、と分かる。
2. コードを数える：goroutine は `jobs`（`2,3,4` の **3個**）ぶん起動し、それぞれ 1 回だけ `results` に送信する → 合計 **3 回送信**。
3. 一方 main の受信ループは `for i := 0; i < 4; i++` で **4 回受信**しようとしている。
4. 3 個の送信を受け取り終えると、送信側 goroutine は全て終了。残った main は 4 回目の受信で送り手のいないチャネルを待ち続ける。動ける goroutine が 1 つもなくなり、ランタイムがデッドロックを検出して落とす。

**原因**
送信回数（goroutine の個数 = 3）と受信回数（4）が食い違っている。受信が 1 回多く、その 1 回が永遠に満たされない。

**修正**
```go
for i := 0; i < len(jobs); i++ {
	sum += <-results
}
```
受信回数を送信回数（`len(jobs)`）に合わせる。修正後は `二乗の合計: 29`（4+9+16）。

**デバッガで goroutine を見る練習**
- 今回はランタイムが自動でデッドロックを検出して終了するため、main の停止位置はダンプで分かる。
- **検出されないハング**（一部の goroutine が動き続けているなどで `all goroutines asleep` にならず、ただ固まる）の場合は、デバッグ実行中に**一時停止 (Pause)** ボタンを押すと、その瞬間で止まる。**コールスタックパネルの上部で goroutine を切り替える**と、各 goroutine が今どの行で何待ちかを一つずつ確認できる。「お互いがお互いの結果を待っている」円環も、こうして複数 goroutine の停止位置を並べて見ると見つけられる。
- 補足：チャネルの送受信数の不一致のほか、`sync.WaitGroup` の `Done()` 忘れ（`Wait()` が永久に返らない）も定番。どちらも「誰が・どこで・何を待っているか」を goroutine ごとに見るのが突破口です。

**詰まる過程をデバッガで観察するステップ**

1. 17行目（`sum += <-results`）にブレークポイントを置いて `F5`。最初の受信の直前で止まる。
2. **コールスタックパネル**を見る。`Goroutine 1 - main.main` に加えて、`Goroutine 6/7/8 - main.main.func1` が並び、いずれも送信行（11行目）で止まっているのが見える。これが goroutine 切り替えの観察ポイント。「送り手は受け手待ちで全員 park している」状況がそのまま見える。
3. ワーカー goroutine をクリックすると変数パネルがそのワーカーの文脈に切り替わる。各自の `n`（2 / 3 / 4）を確認できる。
4. `F5`（続行）を押すたびに次の受信へ進む。変数パネルの `i` と `sum` が進み、コールスタックのワーカー goroutine が 1 つずつ消えていく（受信が成立した送り手は仕事を終えて終了するため）。
5. `i` が `2` まで進んだ時点（3 回受信完了）で、コールスタックには `main` しか残っていないことを確認する。送り手はもういない。
6. 4 回目の受信（`i == 3`）で **ステップオーバー（`F10`）** を押すと、コマンドが返ってこない（ツールバーが「実行中」のまま何も進まない）。「ステップが返ってこない」こと自体が、その行でブロックしたという症状。
7. **一時停止（Pause、⏸）** を押す。実行はランタイムのチャネル受信処理の中で止まり、コールスタックは `main` が 17 行目の受信で待っていることを示す。送り手ゼロの受信で永久に待つ、と確定できる。

---

## まとめ：print デバッグからの卒業ポイント

- **症状の場所 ≠ 原因の場所**。panic の行は「結果」で、原因は手前の関数にあることが多い（演習2・4）。コールスタックとステップインで遡る。
- **止めて覗く**ので、コードを汚さない・消し忘れない・再ビルド不要。条件付きブレークポイントで「問題のケースだけ」を狙い撃ちできる（演習3）。
- **`Uncaught Exceptions` で自動停止**を有効にしておくと、panic の瞬間の状態をそのまま調べられる。
- Go 特有の頻出 panic：範囲外添字 / nil ポインタ参照 / nil マップ代入 / 0 除算 / 型アサーション失敗。いずれも「想定外の値がどこから来たか」をデバッガで遡れば原因に行き着きます。
