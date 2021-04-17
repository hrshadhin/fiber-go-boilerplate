-- Insert demo data into user table
-- Password: demo#123
INSERT INTO "user" (is_admin,username,email,password,first_name,last_name) values
(TRUE, 'hrshadhin', 'dev@hrshadhin.me', '$2a$10$AWuR/dHlOdwYY0Vwtez28.thr67ir8LoB964QQr8QS2tX/eYKh8yS', 'H.R.', 'Shadhin'),
(FALSE, 'demo1', 'demo@hrshadhin.me', '$2a$10$AWuR/dHlOdwYY0Vwtez28.thr67ir8LoB964QQr8QS2tX/eYKh8yS', 'Mr.', 'Demo 1'),
(FALSE, 'demo2', 'demo2@hrshadhin.me', '$2a$10$AWuR/dHlOdwYY0Vwtez28.thr67ir8LoB964QQr8QS2tX/eYKh8yS', 'Mr.', 'Demo 2'),
(FALSE, 'demo3', 'demo3@hrshadhin.me', '$2a$10$AWuR/dHlOdwYY0Vwtez28.thr67ir8LoB964QQr8QS2tX/eYKh8yS', 'Mr.', 'Demo 3'),
(FALSE, 'demo4', 'demo4@hrshadhin.me', '$2a$10$AWuR/dHlOdwYY0Vwtez28.thr67ir8LoB964QQr8QS2tX/eYKh8yS', 'Mr.', 'Demo 4');

