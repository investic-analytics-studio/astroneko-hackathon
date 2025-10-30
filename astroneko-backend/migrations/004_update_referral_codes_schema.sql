-- Create user referral codes table if not exists
CREATE TABLE IF NOT EXISTS astroneko_user_referral_codes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    referral_code varchar(255) NOT NULL,
    is_activated bool DEFAULT false NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT astroneko_user_referral_codes_pkey PRIMARY KEY (id),
    CONSTRAINT astroneko_user_referral_codes_user_id_fkey FOREIGN KEY (user_id) REFERENCES astroneko_auth_users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create referral logs table if not exists
CREATE TABLE IF NOT EXISTS astroneko_referral_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    redeemed_by_user_id uuid NOT NULL,
    code_type varchar(255) NOT NULL,
    referral_code varchar(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT astroneko_referral_logs_pkey PRIMARY KEY (id),
    CONSTRAINT astroneko_referral_logs_redeemed_by_user_id_fkey FOREIGN KEY (redeemed_by_user_id) REFERENCES astroneko_auth_users(id) ON DELETE CASCADE ON UPDATE CASCADE
);