CREATE TABLE evaluations (
    id UUID PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    cv_path TEXT,
    report_path TEXT,
    result JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);