-- Migration: Create users table
-- Created: 2025-01-10
-- Description: Creates the astroneko_auth_users table for user authentication

CREATE TABLE IF NOT EXISTS astroneko_auth_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    is_activated_referral BOOLEAN NOT NULL DEFAULT FALSE,
    latest_login_at TIMESTAMP WITH TIME ZONE,
    firebase_uid VARCHAR(255) NOT NULL,
    profile_image_url TEXT,
    display_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_astroneko_auth_users_email ON astroneko_auth_users(email);

-- Create index on firebase_uid for faster lookups
CREATE INDEX IF NOT EXISTS idx_astroneko_auth_users_firebase_uid ON astroneko_auth_users(firebase_uid);
