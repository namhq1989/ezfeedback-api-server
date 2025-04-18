-- Create enum types with check conditions to prevent errors
CREATE TYPE status AS ENUM ('active', 'inactive');

CREATE TYPE user_status AS ENUM ('active', 'inactive', 'deleted');

CREATE TYPE project_role AS ENUM ('owner', 'editor', 'viewer');

CREATE TYPE feedback_state AS ENUM ('new', 'in_review', 'planned', 'in_progress', 'completed', 'declined');

CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'expired', 'declined');

-- =============================================
-- User Authentication and Management
-- =============================================

-- Create users table (plural form to avoid keyword conflict)
CREATE TABLE users (
                       id TEXT PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       name VARCHAR(255) DEFAULT '' NOT NULL,
                       status user_status NOT NULL,
                       invited_by TEXT REFERENCES users(id) ON DELETE SET NULL DEFAULT NULL,
                       created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                       updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_invited_by ON users(invited_by);

-- Create verification codes table
CREATE TABLE verification_codes (
                                    id TEXT PRIMARY KEY,
                                    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    code VARCHAR(10) NOT NULL,
                                    expires_at TIMESTAMPTZ NOT NULL,
                                    is_used BOOLEAN NOT NULL,
                                    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_verification_codes_user_id ON verification_codes(user_id);
CREATE INDEX idx_verification_codes_code ON verification_codes(code);
CREATE INDEX idx_verification_codes_expires_at ON verification_codes(expires_at);

-- Create user sessions table
CREATE TABLE user_sessions (
                               id TEXT PRIMARY KEY,
                               user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               device_id VARCHAR(255) NOT NULL,
                               refresh_token VARCHAR(255) NOT NULL UNIQUE,
                               expiry_time TIMESTAMPTZ NOT NULL,
                               device_info TEXT DEFAULT '{}' NOT NULL,
                               created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                               updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_device_id ON user_sessions(device_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX idx_user_sessions_expiry_time ON user_sessions(expiry_time);

-- =============================================
-- Project Management
-- =============================================

-- Create projects table
CREATE TABLE projects (
                          id TEXT PRIMARY KEY,
                          user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          title VARCHAR(255) NOT NULL,
                          description TEXT DEFAULT '' NOT NULL,
                          slug VARCHAR(255) UNIQUE NOT NULL,
                          status status NOT NULL,
                          created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                          updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_slug ON projects(slug);
CREATE INDEX idx_projects_status ON projects(status);

-- Create project settings table
CREATE TABLE project_settings (
                                  id TEXT PRIMARY KEY,
                                  project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                  is_feedback_public BOOLEAN NOT NULL,
                                  allow_anonymous_feedback BOOLEAN NOT NULL,
                                  require_approval BOOLEAN NOT NULL,
                                  enable_voting BOOLEAN NOT NULL,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  UNIQUE(project_id)
);

CREATE INDEX idx_project_settings_project_id ON project_settings(project_id);
CREATE INDEX idx_project_settings_is_feedback_public ON project_settings(is_feedback_public);

-- Create project collaborators table
CREATE TABLE project_collaborators (
                                       id TEXT PRIMARY KEY,
                                       project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                       user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                       role project_role NOT NULL,
                                       created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                       updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                       UNIQUE(project_id, user_id)
);

CREATE INDEX idx_project_collaborators_project_id ON project_collaborators(project_id);
CREATE INDEX idx_project_collaborators_user_id ON project_collaborators(user_id);

-- =============================================
-- Feedback Management
-- =============================================

-- Create project categories table
CREATE TABLE project_categories (
                                    id TEXT PRIMARY KEY,
                                    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                    name VARCHAR(100) NOT NULL,
                                    slug VARCHAR(100) NOT NULL,
                                    status status NOT NULL,
                                    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                    UNIQUE(project_id, name),
                                    UNIQUE(project_id, slug)
);

CREATE INDEX idx_project_categories_project_id ON project_categories(project_id);
CREATE INDEX idx_project_categories_name ON project_categories(name);
CREATE INDEX idx_project_categories_slug ON project_categories(slug);
CREATE INDEX idx_project_categories_status ON project_categories(status);

-- Create feedback tags table (project-specific tags)
CREATE TABLE feedback_tags (
                               id TEXT PRIMARY KEY,
                               project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                               name VARCHAR(100) NOT NULL,
                               slug VARCHAR(100) NOT NULL,
                               status status NOT NULL,
                               created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                               updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                               UNIQUE(project_id, name),
                               UNIQUE(project_id, slug)
);

CREATE INDEX idx_feedback_tags_project_id ON feedback_tags(project_id);
CREATE INDEX idx_feedback_tags_name ON feedback_tags(name);
CREATE INDEX idx_feedback_tags_slug ON feedback_tags(slug);
CREATE INDEX idx_feedback_tags_status ON feedback_tags(status);

-- Create feedbacks table
CREATE TABLE feedbacks (
                           id TEXT PRIMARY KEY,
                           project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                           user_id TEXT REFERENCES users(id) ON DELETE SET NULL DEFAULT NULL,
                           category_id TEXT REFERENCES project_categories(id) ON DELETE SET NULL DEFAULT NULL,
                           content TEXT NOT NULL,
                           rating INTEGER DEFAULT NULL CHECK (rating IS NULL OR (rating >= 1 AND rating <= 5)),
                           is_anonymous BOOLEAN NOT NULL,
                           tags TEXT[] DEFAULT '{}' NOT NULL,
                           state feedback_state NOT NULL,
                           status status NOT NULL,
                           created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                           updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_feedbacks_project_id ON feedbacks(project_id);
CREATE INDEX idx_feedbacks_category_id ON feedbacks(category_id);
CREATE INDEX idx_feedbacks_user_id ON feedbacks(user_id);
CREATE INDEX idx_feedbacks_status ON feedbacks(status);
CREATE INDEX idx_feedbacks_state ON feedbacks(state);
CREATE INDEX idx_feedbacks_tags ON feedbacks USING GIN (tags);

-- Create feedback replies table
CREATE TABLE feedback_replies (
                                  id TEXT PRIMARY KEY,
                                  feedback_id TEXT NOT NULL REFERENCES feedbacks(id) ON DELETE CASCADE,
                                  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                  content TEXT NOT NULL,
                                  is_edited BOOLEAN NOT NULL,
                                  status status NOT NULL,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_feedback_replies_feedback_id ON feedback_replies(feedback_id);
CREATE INDEX idx_feedback_replies_user_id ON feedback_replies(user_id);
CREATE INDEX idx_feedback_replies_status ON feedback_replies(status);

-- Create feedback votes table
CREATE TABLE feedback_votes (
                                id TEXT PRIMARY KEY,
                                feedback_id TEXT NOT NULL REFERENCES feedbacks(id) ON DELETE CASCADE,
                                user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                UNIQUE(feedback_id, user_id)
);

CREATE INDEX idx_feedback_votes_feedback_id ON feedback_votes(feedback_id);
CREATE INDEX idx_feedback_votes_user_id ON feedback_votes(user_id);

-- Create user invitations table
CREATE TABLE user_invitations (
                                  id TEXT PRIMARY KEY,
                                  email VARCHAR(255) NOT NULL,
                                  inviter_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                  project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                  role project_role NOT NULL,
                                  token VARCHAR(100) UNIQUE NOT NULL,
                                  code VARCHAR(10) NOT NULL,
                                  status invitation_status NOT NULL,
                                  expires_at TIMESTAMPTZ NOT NULL,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Essential lookup indexes
CREATE UNIQUE INDEX idx_user_invitations_token ON user_invitations(token);
CREATE UNIQUE INDEX idx_user_invitations_active_code ON user_invitations(code) WHERE status = 'pending';

-- Composite indexes for common query patterns
CREATE INDEX idx_user_invitations_email_status ON user_invitations(email, status);
CREATE INDEX idx_user_invitations_project_status ON user_invitations(project_id, status);
CREATE INDEX idx_user_invitations_project_email_created ON user_invitations(project_id, email, created_at DESC);

-- Status + expiration for background jobs
CREATE INDEX idx_user_invitations_status_expires ON user_invitations(status, expires_at);

-- Partial index for pending invitations
CREATE INDEX idx_user_invitations_pending ON user_invitations(email, project_id) WHERE status = 'pending';