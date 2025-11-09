-- User preferences table

CREATE TABLE user_preferences (
    login VARCHAR(255) PRIMARY KEY,
    noon TIME NOT NULL DEFAULT '00:00:00',
    lang VARCHAR(255) NOT NULL DEFAULT 'ru'
);

CREATE INDEX idx_user_preferences_login ON user_preferences(login);

