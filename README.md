# PostFlow

PostFlow is a **Command-Line Post Production Project Management tool** designed to operate at scale for professional media workflows. Built in **Go (Golang)** for performance and reliability, PostFlow focuses on fast bulk uploads/downloads, strict access control, and flexible infrastructure configuration.

Think of PostFlow as a **CLI-first alternative to tools like Frame.io**, optimized for engineers, studios, and technical teams who want full control over their storage, permissions, and deployment environment.

---

## ✨ Key Features

- **CLI-first workflow** – Designed for automation, scripting, and power users  
- **High-performance Go backend** – Optimized for large file transfers  
- **Amazon S3–backed storage** – Bring your own bucket, keys, and IAM policies  
- **Role-based access control** – Admin, Staff, and Viewer permissions  
- **Scalable architecture** – Designed to run locally, on-prem, or on any hosted server  
- **Project-level isolation** – Each project manages its own members and assets  
- **Presigned asset access** – Secure, time-limited S3 download links  

---

## 🧠 Architecture Overview

PostFlow is split into two primary components:

1. **CLI Client**
   - Handles user commands, local file traversal, and API communication
   - Designed to be scriptable and automation-friendly

2. **Server**
   - Exposes HTTP endpoints for authentication, projects, and assets
   - Can be hosted at any URL of the user’s choosing
   - Persists metadata to a database
   - Uploads and downloads files directly to/from Amazon S3

### Storage Model

PostFlow does **not** manage storage for you. Instead:

- You create and own your **Amazon S3 bucket**
- You define **IAM permissions** appropriate for your workflow
- You provide PostFlow with your S3 credentials
- All assets are stored in your bucket, under your control

This design allows for:
- Full compliance with studio security policies
- Easy migration or backup strategies
- Cost control and transparency

---

## 🔐 Authentication & Authorization

PostFlow uses a **refresh token + access token** authentication model:

- **Refresh tokens**
  - Long-lived
  - Stored server-side
  - Used to re-issue access keys when expired

- **Access keys**
  - Short-lived
  - Required for all authenticated requests
  - Automatically rotated when needed

Permissions are enforced at both the **project** and **asset** level.

---

## 🧑‍🤝‍🧑 User Roles

Each project supports three user roles:

- **Admin**
  - Full project control
  - Add/remove members
  - Delete projects
  - Manage all assets

- **Staff**
  - Upload and delete assets
  - View and download assets

- **Viewer**
  - View and download assets only
  - No upload or deletion permissions

---

## 📦 Commands

Below is a complete list of PostFlow CLI commands and their behavior.

---

### Authentication

#### `Register <username> <user_email> <password>`

Create a new user account.

- Creates a new user in the database
- Issues a **refresh token** and **access key**
- Automatically logs the user in

---

#### `Login <email> <password>`

Authenticate an existing user.

- Validates credentials
- Checks refresh token expiration
- Issues a new refresh token if expired
- Issues a new access key

---

#### `Help`

Displays an overview of all available PostFlow commands.

---

## 📁 Projects

#### `Projects create <project_name> <optional_description>`

Create a new project.

- Registers the project in the database
- The creator is automatically assigned **Admin** status

---

#### `Projects addmem <project_name> <user_email> <user_status>`

Add a user to a project.

- User status must be one of:
  - `admin`
  - `staff`
  - `viewer`

---

#### `Projects ls`

List all projects the currently logged-in user is a member of.

---

#### `Projects userlist <project_name>`

List all users associated with a project.

- Displays username and email for each member

---

#### `Projects delete <project_name>`

Delete a project and all related data.

- Admin-only
- Uses **ON DELETE CASCADE** to remove:
  - Project members
  - Assets
  - Metadata
- Does **not** delete raw files from S3 unless explicitly configured to do so

---

#### `Project delmem <project_name> <user_email>`

Remove a user from a project.

- Admin-only
- If removing another admin:
  - Only the **original project creator** may perform the action

---

#### `Projects clone <project_name> <destination_directory>`

Clone a project to the local filesystem.

- Downloads all assets
- Organizes files into folders based on their tag  
  (e.g. `Audio/`, `Rough_Cuts/`)

---

#### `Projects push <source_directory> <project_name> <optional_description>`

Push a local directory into PostFlow.

- Creates asset records in the database
- Uploads all files to S3
- Uses local folder names as asset tags

---

## 🎞 Assets

#### `Assets ls <project_name>`

List all assets in a project.

- Verifies project membership before displaying results

---

#### `Assets upload <project_name> <asset_filepath> <asset_tag>`

Upload an asset to a project.

- Copies the file to S3
- Stores metadata and tags in the database
- Requires **Staff** or **Admin** permissions

---

#### `Assets view <project_name> <asset_name>`

View/download an asset.

- Verifies project membership
- Returns a **presigned S3 URL**
- URL expires automatically for security

---

#### `Assets delete <project_name> <asset_name>`

Delete an asset from a project.

- Requires **Staff** or **Admin** permissions
- Removes metadata and deletes the file from S3

---

## 🚀 Use Cases

- Post-production studios managing large audio/video assets
- Technical teams needing automation-friendly project management
- Secure client review workflows without web UIs
- Internal tooling for media-heavy pipelines
- Backend-focused portfolio projects demonstrating real-world architecture

---

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Storage:** Amazon S3
- **Authentication:** Refresh tokens + access keys
- **Interface:** Command Line Interface (CLI)
- **Architecture:** Client/Server, API-driven

---

## 📄 License

This project is provided as-is for educational and portfolio purposes.  
Refer to the license file in this repository for usage details.

---


