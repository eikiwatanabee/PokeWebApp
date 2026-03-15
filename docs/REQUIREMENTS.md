# PokeBookManager - 要件定義書

## 1. プロジェクト概要

ポケモンをテーマにした読書管理Webアプリケーション。
本を読了するとランダムでポケモンをゲットでき、図鑑を埋めていく楽しさを提供する。

## 2. コアコンセプト

**「本を読む = ポケモンをゲットする冒険」**

- 読了 → ランダムでポケモン1匹ゲット（PokeAPIから取得）
- 進化なし。進化後のポケモンも普通にランダムで出現
- レベルなし。シンプルにゲット＆図鑑コンプリートにフォーカス

## 3. ユーザー・認証

### 3.1 認証
- **Google OAuth2** によるSSO
- JWTトークンによるAPI認証

### 3.2 テナント
- ユーザーはテナント（組織/グループ）に所属
- テナント内でタグをグローバル共有

### 3.3 ユーザーロール
- **管理者**: テナント管理、ユーザー招待
- **一般ユーザー**: 本の登録・読書・ポケモン管理

## 4. 機能要件

### 4.1 本の管理（MUST）

| 機能 | 詳細 |
|---|---|
| 本の登録 | タイトル、著者、ジャンル（タグ）等を入力して登録 |
| ステータス管理 | 未読 → 読書中 → 読了 の3段階 |
| 本の一覧 | ステータスやタグでフィルタリング・検索 |
| 本の編集・削除 | 登録情報の修正、削除 |

### 4.2 メモ機能（MUST）

| 機能 | 詳細 |
|---|---|
| メモ追加 | 本ごとにシンプルテキストメモを複数追加可能 |
| タイムスタンプ | 各メモに作成日時を自動記録 |
| メモ一覧 | 本ごとのメモ一覧表示（時系列順） |
| メモ編集・削除 | 既存メモの修正、削除 |

### 4.3 タグ機能（MUST）

| 機能 | 詳細 |
|---|---|
| タグ作成 | 自由入力でタグを作成 |
| テナント内共有 | 作成されたタグはテナント内の全ユーザーが利用可能 |
| タグ付け | 本に複数タグを付与可能（N:M） |
| タグ検索 | タグで本を絞り込み |

### 4.4 ポケモンシステム（MUST）

| 機能 | 詳細 |
|---|---|
| ポケモンゲット | 本を読了した時にランダムでポケモン1匹をゲット |
| PokeAPI連携 | PokeAPI v2からポケモン情報（名前、画像、タイプ等）を取得 |
| 図鑑 | ゲットしたポケモンの一覧（図鑑形式）|
| ポケモン詳細 | 各ポケモンの情報 + どの本を読んでゲットしたかの紐付け |

### 4.5 ユーザーダッシュボード（NICE TO HAVE）

- 読書統計（読了数、読書中数）
- ゲットしたポケモン数 / 全ポケモン数
- 最近の読了履歴

## 5. 非機能要件

### 5.1 パフォーマンス
- API応答時間: 200ms以内（DB操作含む）
- PokeAPI呼び出し: キャッシュを活用（レスポンスをDBに保存）

### 5.2 セキュリティ
- OAuth2 + JWT認証
- テナント間のデータ分離
- SQLインジェクション対策（GORMパラメータバインディング）
- CORS適切な設定

### 5.3 可用性
- Docker Compose でローカル環境構築可能
- マイグレーション機能による DB スキーマ管理

## 6. 技術スタック

### 6.1 バックエンド
| 技術 | バージョン | 用途 |
|---|---|---|
| Go | 1.22+ | 言語 |
| Gin | v1.10+ | Webフレームワーク |
| GORM | v2 | ORM + マイグレーション |
| PostgreSQL | 16 | データベース |
| JWT | - | トークン認証 |

### 6.2 フロントエンド
| 技術 | バージョン | 用途 |
|---|---|---|
| Next.js | 15 (App Router) | Reactフレームワーク |
| React | 19 | UI |
| TypeScript | 5.x | 型安全 |
| Tailwind CSS | v4 | スタイリング |

### 6.3 インフラ
| 技術 | 用途 |
|---|---|
| Docker Compose v2 | コンテナオーケストレーション |
| Go 1.22-alpine | APIコンテナ |
| node:22-alpine | フロントエンドコンテナ |

## 7. アーキテクチャ

### 7.1 バックエンドアーキテクチャ
- **Clean Architecture** + **DDD** (Domain-Driven Design)
- **CQRS** (Command Query Responsibility Segregation)
- **Unit of Work** パターン

### 7.2 レイヤー構成

```
Presentation Layer (Gin Handlers)
    ↓
Application Layer (Commands / Queries) ← CQRS
    ↓
Domain Layer (Entities / Value Objects / Domain Services)
    ↓
Infrastructure Layer (Repository Impl / UoW / PokeAPI Client)
```

### 7.3 ドメインモデル

```
[Tenant] ──1:N──▶ [User]
                     │
                     ├──1:N──▶ [Book] ──N:M──▶ [Tag]
                     │           │
                     │           └──1:N──▶ [Memo]
                     │
                     └──1:N──▶ [UserPokemon]

[Tag] ── テナントスコープ（テナント内グローバル）
```

### 7.4 Aggregate境界

| Aggregate Root | 含むEntity/VO |
|---|---|
| User | UserPokemon |
| Book | Memo, BookStatus(VO) |
| Tenant | (単独) |
| Tag | (単独, テナントスコープ) |

### 7.5 ロック順序ルール

デッドロック防止のため、以下の順序でロックを取得する。

```
Tenant(1) → User(2) → Book(3) → Memo(4) → Tag(5) → UserPokemon(6)
```

- 常にこの順序でSELECT ... FOR UPDATEを発行
- 逆順のロック取得は禁止
- コードレビューで順序違反をチェック

## 8. DBスキーマ

### tenants
| カラム | 型 | 説明 |
|---|---|---|
| id | UUID | PK |
| name | VARCHAR(255) | テナント名 |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### users
| カラム | 型 | 説明 |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK → tenants |
| google_id | VARCHAR(255) | Google OAuth ID |
| email | VARCHAR(255) | メールアドレス |
| name | VARCHAR(255) | 表示名 |
| role | VARCHAR(50) | admin / member |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### books
| カラム | 型 | 説明 |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → users |
| title | VARCHAR(255) | 書籍タイトル |
| author | VARCHAR(255) | 著者 |
| status | VARCHAR(50) | unread / reading / finished |
| finished_at | TIMESTAMP | 読了日時（NULL可） |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### memos
| カラム | 型 | 説明 |
|---|---|---|
| id | UUID | PK |
| book_id | UUID | FK → books |
| content | TEXT | メモ内容 |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### tags
| カラム | 型 | 説明 |
|---|---|---|
| id | UUID | PK |
| tenant_id | UUID | FK → tenants |
| name | VARCHAR(100) | タグ名 |
| created_at | TIMESTAMP | 作成日時 |

### book_tags
| カラム | 型 | 説明 |
|---|---|---|
| book_id | UUID | FK → books |
| tag_id | UUID | FK → tags |

### user_pokemon
| カラム | 型 | 説明 |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → users |
| book_id | UUID | FK → books（どの本でゲットしたか） |
| pokemon_id | INT | PokeAPI上のポケモンID |
| pokemon_name | VARCHAR(100) | ポケモン名 |
| pokemon_sprite | VARCHAR(500) | スプライト画像URL |
| pokemon_types | JSONB | タイプ情報 |
| caught_at | TIMESTAMP | ゲット日時 |

## 9. API エンドポイント

### 認証
| Method | Path | 説明 |
|---|---|---|
| GET | /api/auth/google | Google OAuth2 認証開始 |
| GET | /api/auth/google/callback | Google OAuth2 コールバック |
| POST | /api/auth/refresh | トークンリフレッシュ |

### 本
| Method | Path | 説明 | CQRS |
|---|---|---|---|
| POST | /api/books | 本を登録 | Command |
| GET | /api/books | 本一覧取得 | Query |
| GET | /api/books/:id | 本詳細取得 | Query |
| PUT | /api/books/:id | 本情報更新 | Command |
| DELETE | /api/books/:id | 本削除 | Command |
| POST | /api/books/:id/start | 読書開始（→読書中） | Command |
| POST | /api/books/:id/finish | 読了（→ポケモンゲット） | Command |

### メモ
| Method | Path | 説明 | CQRS |
|---|---|---|---|
| POST | /api/books/:id/memos | メモ追加 | Command |
| GET | /api/books/:id/memos | メモ一覧取得 | Query |
| PUT | /api/memos/:id | メモ更新 | Command |
| DELETE | /api/memos/:id | メモ削除 | Command |

### タグ
| Method | Path | 説明 | CQRS |
|---|---|---|---|
| POST | /api/tags | タグ作成 | Command |
| GET | /api/tags | タグ一覧取得（テナント内） | Query |
| DELETE | /api/tags/:id | タグ削除 | Command |

### ポケモン図鑑
| Method | Path | 説明 | CQRS |
|---|---|---|---|
| GET | /api/pokedex | 図鑑一覧 | Query |
| GET | /api/pokedex/:id | ポケモン詳細 | Query |

### テナント
| Method | Path | 説明 | CQRS |
|---|---|---|---|
| GET | /api/tenant | テナント情報取得 | Query |
| POST | /api/tenant/invite | ユーザー招待 | Command |

## 10. 画面構成

| 画面 | パス | 説明 |
|---|---|---|
| ログイン | /login | Google SSO |
| ダッシュボード | / | 読書統計・最近の活動 |
| 本一覧 | /books | 本の一覧・フィルタ・検索 |
| 本詳細 | /books/:id | 本情報・メモ一覧・ステータス変更 |
| 本登録 | /books/new | 新規本登録フォーム |
| ポケモン図鑑 | /pokedex | ゲットしたポケモン一覧 |
| ポケモン詳細 | /pokedex/:id | ポケモン情報・紐付いた本 |
| 設定 | /settings | ユーザー設定・テナント管理 |
