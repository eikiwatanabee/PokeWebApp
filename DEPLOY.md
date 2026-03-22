# デプロイ手順（Vercel + Supabase + Render）

全て無料枠で運用可能な構成です。

## 1. Supabase（データベース）

1. [supabase.com](https://supabase.com) でプロジェクトを作成
2. **Settings → Database** から接続情報を取得:
   - `Connection string (URI)` をコピー（`[YOUR-PASSWORD]` を実際のパスワードに置換）
   - 例: `postgresql://postgres.xxxx:password@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres`
3. この値を `DATABASE_URL` として後述の Render 環境変数に設定

> **注意**: Transaction mode (port 6543) を使用してください。Session mode (port 5432) でも動作します。

## 2. Render（バックエンド API）

1. [render.com](https://render.com) でアカウント作成
2. **New → Web Service** を選択
3. GitHubリポジトリを接続
4. 設定:
   - **Name**: `pokegithub-api`
   - **Region**: Singapore（日本に近い）
   - **Runtime**: Docker
   - **Dockerfile Path**: `./backend/Dockerfile`
   - **Docker Context**: `./backend`
   - **Plan**: Free
5. 環境変数を設定:

| 変数名 | 値 |
|-------|---|
| `PORT` | `8080` |
| `DATABASE_URL` | Supabaseの接続文字列 |
| `JWT_SECRET` | ランダムな文字列（`openssl rand -hex 32`） |
| `GITHUB_CLIENT_ID` | GitHub OAuth App の Client ID |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth App の Client Secret |
| `GITHUB_REDIRECT_URL` | `https://pokegithub-api.onrender.com/api/auth/github/callback` |
| `GITHUB_WEBHOOK_SECRET` | Webhookの署名検証用シークレット |
| `GOOGLE_CLIENT_ID` | （任意）Google OAuth |
| `GOOGLE_CLIENT_SECRET` | （任意）Google OAuth |
| `GOOGLE_REDIRECT_URL` | `https://pokegithub-api.onrender.com/api/auth/google/callback` |

6. **Create Web Service** をクリック → 自動ビルド＆デプロイ

## 3. Vercel（フロントエンド）

1. [vercel.com](https://vercel.com) でアカウント作成
2. **Add New → Project** → GitHubリポジトリをインポート
3. 設定:
   - **Framework Preset**: Next.js
   - **Root Directory**: `frontend`
4. 環境変数を設定:

| 変数名 | 値 |
|-------|---|
| `NEXT_PUBLIC_API_URL` | `https://pokegithub-api.onrender.com`（RenderのURL） |

5. **Deploy** をクリック

## 4. GitHub OAuth App の設定

1. [GitHub Developer Settings](https://github.com/settings/developers) → **New OAuth App**
2. 設定:
   - **Application name**: PokeGitHub
   - **Homepage URL**: Vercelのデプロイ先URL
   - **Authorization callback URL**: `https://pokegithub-api.onrender.com/api/auth/github/callback`
3. Client ID / Client Secret を Render の環境変数に設定

## 5. GitHub Webhook の設定

対象リポジトリの **Settings → Webhooks → Add webhook**:
- **Payload URL**: `https://pokegithub-api.onrender.com/api/webhook/github`
- **Content type**: `application/json`
- **Secret**: Render に設定した `GITHUB_WEBHOOK_SECRET` と同じ値
- **Events**: Push, Pull requests, Pull request reviews, Issues, Deployments

## ローカル開発

```bash
# Docker Compose で全部起動
docker compose up

# または個別起動
cd backend && go run ./cmd/server  # API (localhost:8080)
cd frontend && npm run dev         # Frontend (localhost:3000)
```

## 注意事項

- Render 無料枠はスリープあり（15分無通信でスリープ → 次のリクエストで起動に30秒ほど）
- Supabase 無料枠は500MB、1週間無活動でpause（手動で再開可能）
- Vercel は制限なしで快適に使える
