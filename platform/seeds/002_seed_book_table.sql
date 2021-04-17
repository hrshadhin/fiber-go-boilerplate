-- Insert demo data into book table
INSERT INTO public.book(user_id, title, author, status, meta) VALUES
(2, 'Wolla Pizza', 'H.R. Shadhin', 1, '{ "rating": 5, "picture": "http://pic.me/wollap.png", "description": "test book 1"}'),
(2, 'Wolla Pizza 2', 'H.R. Shadhin', 1, '{ "rating": 2, "picture": "http://pic.me/wollap2.png", "description": "test book 2"}'),
(2, 'Wolla Pizza 3', 'H.R. Shadhin', 1, '{ "rating": 8, "picture": "http://pic.me/wolla3.png", "description": "test book 3"}'),
(2, 'Wolla Pizza 4', 'H.R. Shadhin', 1, '{ "rating": 1, "picture": "http://pic.me/wollap4.png", "description": "test book 4"}'),
(2, 'Wolla Pizza 5', 'H.R. Shadhin', 0, '{ "rating": 0, "picture": "http://pic.me/wollap5.png", "description": "test book 5"}');
