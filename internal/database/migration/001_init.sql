-- Create enum types with check conditions to prevent errors
CREATE TYPE status AS ENUM ('active', 'inactive');

CREATE TYPE user_status AS ENUM ('active', 'inactive', 'deleted');

CREATE TYPE project_role AS ENUM ('owner', 'editor', 'viewer');

CREATE TYPE feedback_state AS ENUM ('new', 'in_review', 'planned', 'in_progress', 'completed', 'declined');

CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'expired', 'declined');

CREATE TYPE campaign_type AS ENUM ('feedback', 'nps', 'csat');

CREATE TYPE notification_type AS ENUM ('new_feedback', 'system');

-- =============================================
-- User Authentication and Management
-- =============================================

-- Create users table (plural form to avoid keyword conflict)
CREATE TABLE users (
                       id TEXT PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       name VARCHAR(255) DEFAULT '' NOT NULL,
                       status user_status NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                       updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- Create verification codes table
CREATE TABLE verification_codes (
                                    id TEXT PRIMARY KEY,
                                    email TEXT NOT NULL,
                                    ip TEXT NOT NULL,
                                    code VARCHAR(10) NOT NULL,
                                    expires_at TIMESTAMPTZ NOT NULL,
                                    is_used BOOLEAN NOT NULL,
                                    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_verification_codes_ip_created_at ON verification_codes(ip, created_at);
CREATE INDEX idx_verification_codes_ip_email_code ON verification_codes(ip, email, code);
CREATE INDEX idx_verification_codes_expires_at ON verification_codes(expires_at);

-- Create user sessions table
CREATE TABLE user_sessions (
                               id TEXT PRIMARY KEY,
                               user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               device_id VARCHAR(255) NOT NULL,
                               refresh_token VARCHAR(255) NOT NULL UNIQUE,
                               expires_at TIMESTAMPTZ NOT NULL,
                               device_info TEXT DEFAULT '{}' NOT NULL,
                               created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                               updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_device_id_refresh_token ON user_sessions(device_id, refresh_token);

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
                          stats_total_feedbacks INTEGER DEFAULT 0 NOT NULL,
                          created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                          updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_slug ON projects(slug);

-- Create project settings table
CREATE TABLE project_settings (
                                  id TEXT PRIMARY KEY,
                                  project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                  domain VARCHAR(255) DEFAULT '' NOT NULL,
                                  primary_color VARCHAR(7) DEFAULT '' NOT NULL,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  UNIQUE(project_id)
);

CREATE INDEX idx_project_settings_project_id ON project_settings(project_id);
CREATE INDEX idx_project_settings_project_domain ON project_settings(project_id, domain);

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

-- Create project campaigns table
CREATE TABLE project_campaigns (
                                   id TEXT PRIMARY KEY,
                                   project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                   name VARCHAR(255) NOT NULL,
                                   campaign_type campaign_type NOT NULL,
                                   status status NOT NULL,
                                   settings JSONB DEFAULT '{}'::JSONB NOT NULL,
                                   stats_total_feedbacks INTEGER DEFAULT 0 NOT NULL,
                                   created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                   updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_project_campaigns_project_campaign_type ON project_campaigns(project_id, campaign_type);

-- Create project campaign categories table
CREATE TABLE project_campaign_categories (
                                  id TEXT PRIMARY KEY,
                                  campaign_id TEXT NOT NULL REFERENCES project_campaigns(id) ON DELETE CASCADE,
                                  category_id TEXT NOT NULL REFERENCES project_categories(id) ON DELETE CASCADE,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  UNIQUE(campaign_id, category_id)
);

CREATE INDEX idx_project_campaign_categories_campaign_id ON project_campaign_categories(campaign_id);

-- Create feedbacks table
CREATE TABLE feedbacks (
                           id TEXT PRIMARY KEY,
                           project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                           campaign_id TEXT NOT NULL REFERENCES project_campaigns(id) ON DELETE CASCADE,
                           app_user_id VARCHAR(255),
                           email VARCHAR(255),
                           category_id TEXT NOT NULL,
                           content TEXT NOT NULL,
                           rating INTEGER NOT NULL CHECK (rating IS NULL OR (rating >= 1 AND rating <= 10)),
                           is_anonymous BOOLEAN NOT NULL,
                           state feedback_state NOT NULL,
                           campaign_type campaign_type NOT NULL,
                           ip VARCHAR(45) NOT NULL,
                           country_code VARCHAR(2) NOT NULL,
                           search_vector tsvector NOT NULL,
                           stats_total_replies INTEGER DEFAULT 0 NOT NULL,
                           created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                           updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_feedbacks_project_id ON feedbacks(project_id);
CREATE INDEX idx_feedbacks_campaign_id ON feedbacks(campaign_id);

CREATE INDEX idx_feedbacks_query_patterns ON feedbacks(project_id, campaign_id, category_id, rating, created_at DESC);
CREATE INDEX idx_feedbacks_search_vector ON feedbacks USING GIN(search_vector);

-- Create a custom text search configuration that preserves stop words
CREATE TEXT SEARCH CONFIGURATION english_nostop (COPY = english);
ALTER TEXT SEARCH CONFIGURATION english_nostop ALTER MAPPING FOR asciihword, asciiword, hword, word WITH simple;

-- Updated function to generate search vector with improved email handling
CREATE OR REPLACE FUNCTION feedbacks_search_vector_update() RETURNS trigger AS $$
BEGIN
  NEW.search_vector =
    setweight(to_tsvector('english_nostop', COALESCE(NEW.content, '')), 'A') ||
    setweight(to_tsvector('english_nostop', COALESCE(NEW.email, '')), 'B') ||
    setweight(to_tsvector('english_nostop',
      COALESCE(regexp_replace(NEW.email, '[@.]', ' ', 'g'), '')
    ), 'B');
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger remains the same
CREATE TRIGGER feedbacks_search_vector_update_trigger
    BEFORE INSERT OR UPDATE ON feedbacks
                         FOR EACH ROW
                         EXECUTE FUNCTION feedbacks_search_vector_update();

-- Create feedback replies table
CREATE TABLE feedback_replies (
                                  id TEXT PRIMARY KEY,
                                  feedback_id TEXT NOT NULL REFERENCES feedbacks(id) ON DELETE CASCADE,
                                  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                  content TEXT NOT NULL,
                                  is_edited BOOLEAN NOT NULL,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_feedback_replies_feedback_id ON feedback_replies(feedback_id);
CREATE INDEX idx_feedback_replies_user_id ON feedback_replies(user_id);

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
                                  status invitation_status NOT NULL,
                                  expires_at TIMESTAMPTZ NOT NULL,
                                  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Composite indexes for common query patterns
CREATE INDEX idx_user_invitations_email_status ON user_invitations(email, status);
CREATE INDEX idx_user_invitations_project_status ON user_invitations(project_id, status);
CREATE INDEX idx_user_invitations_project_email_created ON user_invitations(project_id, email, created_at DESC);

-- Status + expiration for background jobs
CREATE INDEX idx_user_invitations_status_expires ON user_invitations(status, expires_at);

-- Partial index for pending invitations
CREATE INDEX idx_user_invitations_pending ON user_invitations(email, project_id) WHERE status = 'pending';

-- =============================================
-- Notification
-- =============================================

-- Create notification reminders table
CREATE TABLE notification_reminders (
                                id TEXT PRIMARY KEY,
                                project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_notification_reminders_created_at ON notification_reminders(created_at);

-- Create notifications table
CREATE TABLE notifications (
                               id TEXT PRIMARY KEY,
                               user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               type notification_type NOT NULL,
                               is_read BOOLEAN DEFAULT FALSE NOT NULL,
                               metadata JSONB DEFAULT '{}'::JSONB NOT NULL,
                               created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                               updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_notifications_user_id_created_at ON notifications(user_id, created_at);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Create user project notification settings
CREATE TABLE user_project_notification_settings (
                            id TEXT PRIMARY KEY,
                            user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                            receive_new_feedback BOOLEAN DEFAULT TRUE NOT NULL,
                            receive_daily_summary BOOLEAN DEFAULT FALSE NOT NULL,
                            receive_weekly_summary BOOLEAN DEFAULT TRUE NOT NULL,
                            created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                            updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                            UNIQUE(user_id, project_id)
);

CREATE INDEX idx_upns_user_id_project_id ON user_project_notification_settings(user_id, project_id);
