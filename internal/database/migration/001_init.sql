-- Create enum types
CREATE TYPE status AS ENUM ('active', 'inactive', 'deleted');
CREATE TYPE project_role AS ENUM ('owner', 'editor', 'viewer');
CREATE TYPE feedback_state AS ENUM ('new', 'in_review', 'planned', 'in_progress', 'completed', 'declined');
CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'expired', 'declined');
CREATE TYPE vote_type AS ENUM ('up', 'down');

-- =============================================
-- User Authentication and Management
-- =============================================

-- Create users table
CREATE TABLE user (
                      id TEXT PRIMARY KEY,
                      email VARCHAR(255) UNIQUE NOT NULL,
                      name VARCHAR(255) DEFAULT NULL,
                      status status NOT NULL,
                      invited_by TEXT REFERENCES user(id) ON DELETE SET NULL DEFAULT NULL,
                      created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                      updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_user_email ON user(email);
CREATE INDEX idx_user_status ON user(status);
CREATE INDEX idx_user_invited_by ON user(invited_by);

-- Create verification codes table
CREATE TABLE sign_in_verification_code (
                                   id TEXT PRIMARY KEY,
                                   user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
                                   code VARCHAR(10) NOT NULL,
                                   expires_at TIMESTAMPTZ NOT NULL,
                                   is_used BOOLEAN NOT NULL,
                                   created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_sign_in_verification_code_user_id ON sign_in_verification_code(user_id);
CREATE INDEX idx_sign_in_verification_code_code ON sign_in_verification_code(code);
CREATE INDEX idx_sign_in_verification_code_expires_at ON sign_in_verification_code(expires_at);

-- Create user session table
CREATE TABLE user_session (
                              id TEXT PRIMARY KEY,
                              user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
                              device_id VARCHAR(255) NOT NULL UNIQUE,
                              refresh_token VARCHAR(255) NOT NULL UNIQUE,
                              expiry_time TIMESTAMPTZ NOT NULL,
                              device_info TEXT DEFAULT NULL,
                              ip VARCHAR(255) DEFAULT NULL,
                              created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                              updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_user_session_user_id ON user_session(user_id);
CREATE INDEX idx_user_session_refresh_token ON user_session(refresh_token);
CREATE INDEX idx_user_session_expiry_time ON user_session(expiry_time);
CREATE INDEX idx_user_session_is_active ON user_session(is_active);

-- Create user invitation table
CREATE TABLE user_invitation (
                                 id TEXT PRIMARY KEY,
                                 email VARCHAR(255) NOT NULL,
                                 inviter_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
                                 project_id TEXT REFERENCES project(id) ON DELETE CASCADE DEFAULT NULL,
                                 role project_role NOT NULL,
                                 token VARCHAR(100) UNIQUE NOT NULL,
                                 code VARCHAR(10) NOT NULL,
                                 status invitation_status NOT NULL,
                                 expires_at TIMESTAMPTZ NOT NULL,
                                 created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                 updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Essential lookup indexes
CREATE UNIQUE INDEX idx_user_invitation_token ON user_invitation(token);
CREATE UNIQUE INDEX idx_user_invitation_active_code ON user_invitation(code) WHERE status = 'pending';

-- Composite indexes for common query patterns
CREATE INDEX idx_user_invitation_email_status ON user_invitation(email, status);
CREATE INDEX idx_user_invitation_project_status ON user_invitation(project_id, status);
CREATE INDEX idx_user_invitation_project_email_created ON user_invitation(project_id, email, created_at DESC);

-- Status + expiration for background jobs
CREATE INDEX idx_user_invitation_status_expires ON user_invitation(status, expires_at);

-- Partial index for pending invitations
CREATE INDEX idx_user_invitation_pending ON user_invitation(email, project_id) WHERE status = 'pending';

-- =============================================
-- Project Management
-- =============================================

-- Create projects table
CREATE TABLE project (
                         id TEXT PRIMARY KEY,
                         user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
                         title VARCHAR(255) NOT NULL,
                         description TEXT DEFAULT NULL,
                         slug VARCHAR(255) UNIQUE NOT NULL,
                         status status NOT NULL,
                         created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                         updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_project_user_id ON project(user_id);
CREATE INDEX idx_project_slug ON project(slug);
CREATE INDEX idx_project_status ON project(status);

-- Create project settings table
CREATE TABLE project_setting (
                                 id TEXT PRIMARY KEY,
                                 project_id TEXT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
                                 is_feedback_public BOOLEAN NOT NULL,
                                 allow_anonymous_feedback BOOLEAN NOT NULL,
                                 require_approval BOOLEAN NOT NULL,
                                 enable_voting BOOLEAN NOT NULL,
                                 created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                 updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                 UNIQUE(project_id)
);

CREATE INDEX idx_project_setting_project_id ON project_setting(project_id);
CREATE INDEX idx_project_setting_is_feedback_public ON project_setting(is_feedback_public);

-- Create project collaborator table
CREATE TABLE project_collaborator (
                                      id TEXT PRIMARY KEY,
                                      project_id TEXT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
                                      user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
                                      role project_role NOT NULL,
                                      is_active BOOLEAN NOT NULL,
                                      created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                      updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                      UNIQUE(project_id, user_id)
);

CREATE INDEX idx_project_collaborator_project_id ON project_collaborator(project_id);
CREATE INDEX idx_project_collaborator_user_id ON project_collaborator(user_id);
CREATE INDEX idx_project_collaborator_is_active ON project_collaborator(is_active);

-- =============================================
-- Feedback Management
-- =============================================

-- Create project category table
CREATE TABLE project_category (
                                  id TEXT PRIMARY KEY,
                                  project_id TEXT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
                                  name VARCHAR(100) NOT NULL,
                                  slug VARCHAR(100) NOT NULL,
                                  status status NOT NULL,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  UNIQUE(project_id, name),
                                  UNIQUE(project_id, slug)
);

CREATE INDEX idx_project_category_project_id ON project_category(project_id);
CREATE INDEX idx_project_category_name ON project_category(name);
CREATE INDEX idx_project_category_slug ON project_category(slug);
CREATE INDEX idx_project_category_status ON project_category(status);

-- Create feedback tag table (project-specific tags)
CREATE TABLE feedback_tag (
                              id TEXT PRIMARY KEY,
                              project_id TEXT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
                              name VARCHAR(100) NOT NULL,
                              slug VARCHAR(100) NOT NULL,
                              status status NOT NULL,
                              created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                              updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                              UNIQUE(project_id, name),
                              UNIQUE(project_id, slug)
);

CREATE INDEX idx_feedback_tag_project_id ON feedback_tag(project_id);
CREATE INDEX idx_feedback_tag_name ON feedback_tag(name);
CREATE INDEX idx_feedback_tag_slug ON feedback_tag(slug);
CREATE INDEX idx_feedback_tag_status ON feedback_tag(status);

-- Create feedback table
CREATE TABLE feedback (
                          id TEXT PRIMARY KEY,
                          project_id TEXT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
                          user_id TEXT REFERENCES user(id) ON DELETE SET NULL DEFAULT NULL,
                          category_id TEXT REFERENCES project_category(id) ON DELETE SET NULL DEFAULT NULL,
                          content TEXT NOT NULL,
                          rating INTEGER DEFAULT NULL CHECK (rating IS NULL OR (rating >= 1 AND rating <= 5)),
                          is_anonymous BOOLEAN NOT NULL,
                          tags TEXT[] DEFAULT '{}' NOT NULL,
                          state feedback_state NOT NULL,
                          status status NOT NULL,
                          created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                          updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_feedback_project_id ON feedback(project_id);
CREATE INDEX idx_feedback_category_id ON feedback(category_id);
CREATE INDEX idx_feedback_user_id ON feedback(user_id);
CREATE INDEX idx_feedback_status ON feedback(status);
CREATE INDEX idx_feedback_state ON feedback(state);
CREATE INDEX idx_feedback_tags ON feedback USING GIN (tags);

-- Create feedback reply table
CREATE TABLE feedback_reply (
                                id TEXT PRIMARY KEY,
                                feedback_id TEXT NOT NULL REFERENCES feedback(id) ON DELETE CASCADE,
                                user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
                                content TEXT NOT NULL,
                                is_edited BOOLEAN NOT NULL,
                                status status NOT NULL,
                                created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_feedback_reply_feedback_id ON feedback_reply(feedback_id);
CREATE INDEX idx_feedback_reply_user_id ON feedback_reply(user_id);
CREATE INDEX idx_feedback_reply_status ON feedback_reply(status);

-- Create feedback vote table
CREATE TABLE feedback_vote (
                               id TEXT PRIMARY KEY,
                               feedback_id TEXT NOT NULL REFERENCES feedback(id) ON DELETE CASCADE,
                               user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
                               vote_type vote_type NOT NULL,
                               created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                               UNIQUE(feedback_id, user_id)
);

CREATE INDEX idx_feedback_vote_feedback_id ON feedback_vote(feedback_id);
CREATE INDEX idx_feedback_vote_user_id ON feedback_vote(user_id);
CREATE INDEX idx_feedback_vote_vote_type ON feedback_vote(vote_type);