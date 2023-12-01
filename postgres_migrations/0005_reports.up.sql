CREATE TABLE content_report (
  id SERIAL PRIMARY KEY,
  content_ref TEXT NOT NULL,
  report_state TEXT NOT NULL,
  reported_by TEXT NOT NULL,
  report_reason citext NOT NULL,
  context TEXT NOT NULL,
  reported_user TEXT NOT NULL,
  resolved_by TEXT,
  resolved_at TIMESTAMP,
  action_taken citext,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);