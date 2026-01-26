
# PostFlow

PostFlow is a command‑line tool and Go‑based HTTP API for managing post‑production projects and media assets. It is designed to be fast, usable at scale, and deployable in environments where you control the database, object storage, and server URL.

At its core, PostFlow provides:

- User registration and authentication using access/refresh tokens.
- Project creation, membership, and role management (`admin`, `staff`, `viewer`).
- Asset upload, listing, download, and deletion, backed by PostgreSQL and S3.
- Bulk project import/export between a local filesystem and remote storage.

The same binary runs both the server (`postflow serve`) and the CLI commands (e.g. `postflow projects ls`, `postflow assets upload ...`).

---

## Components and Data Flow

- **Server**
  - Go HTTP server.
  - Uses PostgreSQL for users, projects, membership, and asset metadata.
  - Uses Amazon S3 for asset files.
  - Issues presigned S3 URLs so clients can upload/download directly.

- **CLI**
  - Single binary (`postflow`) installed on the user’s machine.
  - Talks only to the server over HTTP/JSON.
  - Provides commands to:
    - Register/login users.
    - Create and manage projects and membership.
    - Upload and download assets.
    - Push local directory trees into a project, and clone projects back to disk.

This separation allows you to run the server on any reachable host while developers and operators use the CLI from their own machines.

---

## Prerequisites

- Go 1.22+ (for building the binary).
- PostgreSQL 13+ (local or remote).
- An S3 bucket with appropriate permissions.
- Docker (optional) if you prefer to run the server in a container.

---

## Configuration

The server is configured through environment variables.

Typical server environment:

```env
# PostgreSQL
DATABASE_URL="postgres://postgres:postgres@localhost:5432/postflow?sslmode=disable"

# HTTP listen address for the server
PORT=":8080"

# S3 / AWS
AWS_REGION="us-east-1"
S3_BUCKET="your-postflow-bucket-name"
AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY_ID"
AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_ACCESS_KEY"


The CLI needs to know where the server is running:

```env
# Base URL used by the CLI to reach the server
BASE_URL="http://localhost:8080"
```

You can keep these in a local `.env` file (not committed) and export them via your shell or a tool like `direnv`.

---

## Running the Server (local)

From the repository root:

1. Ensure your `DATABASE_URL`, `AWS_REGION`, `S3_BUCKET`, and AWS credentials are exported.
2. Build and start the server:

```bash
go build -o postflow .
./postflow serve
```

The server will listen on `PORT` (default `:8080`), for example at `http://localhost:8080`.

---

## Running the Server (Docker)

The repo includes a `Dockerfile` that builds a minimal server image.

Example:

```bash
# Build image
docker build -t postflow-server .

# Run container
docker run --rm \
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/postflow?sslmode=disable" \
  -e AWS_REGION="us-east-1" \
  -e S3_BUCKET="your-postflow-bucket-name" \
  -e AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY_ID" \
  -e AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_ACCESS_KEY" \
  -e PORT=":8080" \
  -p 8080:8080 \
  postflow-server
```

Adjust `DATABASE_URL` and S3 settings to match your environment. Once running, the server is reachable at `http://localhost:8080` (or any host/port you configure).

---

## Installing the CLI

To use PostFlow as a global CLI on your machine:

```bash
# From repo root
go build -o postflow .

# Install into a directory on PATH (macOS/Linux)
sudo mv postflow /usr/local/bin/postflow
sudo chmod +x /usr/local/bin/postflow
```

Verify:

```bash
which postflow
postflow help
```

Point the CLI at your server:

```bash
export BASE_URL="http://localhost:8080"
```

You can now run commands from any directory, for example:

```bash
postflow register <username> <email> <password>
postflow login <email> <password>
postflow projects create <project_name> [description]
```

---

## Core CLI Functionality

The CLI surface is intentionally compact and focused on common workflows.

High‑level capabilities include:

- **User and session management**
  - Register and login users.
  - Manage access tokens automatically on the client.

- **Projects**
  - Create projects with optional descriptions.
  - List projects for the current user.
  - Add and remove members with role control (`admin`, `staff`, `viewer`).
  - Inspect project membership.
  - Delete projects and cascade related metadata and assets.

- **Assets**
  - Upload assets to a project:
    - Metadata stored in PostgreSQL.
    - File content stored in S3 under a structured key.
  - List assets for a project.
  - Obtain presigned download URLs for assets.
  - Delete assets from both the database and S3.

- **Bulk operations**
  - `projects push`:
    - Create a project.
    - Walk a local directory tree.
    - Use the innermost folder name as each file’s tag.
    - Upload all files to S3 and create asset records concurrently.
  - `projects clone`:
    - Query all assets in a project.
    - Reconstruct a directory tree locally:
      - `<destination>/<project_name>/<tag>/<asset_name>`
    - Download all assets from S3 in parallel using presigned URLs.

For detailed command syntax and options, use the built‑in help:

```bash
postflow help
postflow help projects
postflow help assets
```

---

## Typical Usage Flow

A typical usage sequence looks like this:

1. Start the server (locally or via Docker).
2. Configure the CLI with `BASE_URL`.
3. Register and log in:

   ```bash
   postflow register <username> <email> <password>
   postflow login <email> <password>
   ```

4. Create a project and add collaborators:

   ```bash
   postflow projects create my_project "Rough cut and final assets"
   postflow projects addmem my_project collaborator@example.com staff
   ```

5. Upload assets:

   ```bash
   postflow assets upload my_project /path/to/file.mov "video"
   postflow assets ls my_project
   ```

6. Bulk import or export:

   - Import a local folder tree:

     ```bash
     postflow projects push my_project /path/to/local/folder
     ```

   - Clone a remote project locally:

     ```bash
     postflow projects clone my_project /path/to/destination
     ```

7. Clean up when needed:

   ```bash
   postflow assets delete my_project file.mov
   postflow projects delete my_project
   ```


