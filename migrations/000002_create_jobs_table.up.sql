CREATE TYPE job_status AS ENUM (
    'applied',
    'interview',
    'offer',
    'rejected'
);

CREATE TABLE jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company         VARCHAR(255) NOT NULL,
    role            VARCHAR(255) NOT NULL,
    status          job_status NOT NULL DEFAULT 'applied',
    notes           TEXT,
    applied_date    DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_status ON jobs(status);