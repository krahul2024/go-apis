system clear; 
use first; 

-- describe users; 
-- describe orders ; 
-- describe brands ; 
-- describe categories; 
-- describe products; 


-- INSERT INTO products (name, summary, description, price, quantity, brandId, categoryId, imageUrl, productCode)
-- VALUES
-- ("Samsung Galaxy S24", "Smartphone", "Advanced camera and AI", 800, 10, 1, 1, "samsung-galaxy-s22.jpg", "prod044"),
-- ("Apple MacBook Air M3", "Laptop", "Thin and lightweight", 1000, 10, 25, 3, "apple-macbook-air.jpg", "prod045"),
-- ("Dell XPS 15", "Laptop", "High-performance processing", 1200, 10, 36, 3, "dell-xps-15.jpg", "prod046"),
-- ("Google Pixel 6", "Smartphone", "Advanced camera and AI", 700, 15, 30, 1, "google-pixel-6.jpg", "prod047"),
-- ("Microsoft Surface Book", "Laptop", "2-in-1 design", 1200, 10, 27, 3, "microsoft-surface-book.jpg", "prod050"),
-- ("Sony Xperia 1", "Smartphone", "Advanced camera and AI", 800, 10, 2, 1, "sony-xperia-1.jpg", "prod053"),
-- ("Microsoft Surface Laptop 3", "Laptop", "Thin and lightweight", 1000, 10, 27, 3, "microsoft-surface-laptop-3.jpg", "prod059"),
-- ("Samsung Galaxy Tab S8", "Tablet", "High-performance processing", 500, 10, 1, 5, "samsung-galaxy-tab-s8.jpg", "prod062"),
-- ("Sony Xperia 5", "Smartphone", "Advanced camera and AI", 600, 10, 2, 1, "sony-xperia-5.jpg", "prod063"),
-- ("Microsoft Surface Go", "Tablet", "High-performance processing", 400, 10, 27, 5, "microsoft-surface-go.jpg", "prod069"); 


-- create the brands table 
-- create table if not exists products(
--     id int not null auto_increment, 
--     name varchar(200), 
--     summary varchar(500), 
--     description varchar(2000), 
--     price int not null, 
--     quantity int default 0 , 
--     brand_id int not null, 
--     category_id int not null, 
--     image_url varchar(1000), 
--     createdAt datetime default current_timestamp, 
--     updatedAt datetime default current_timestamp on update current_timestamp,
--     primary key (id), 
--     foreign key (brand_id) references brands(id), 
--     foreign key (category_id) references categories(id)
-- );


-- -- create the orders table 
-- create table if not exists orders (
--     id int not null auto_increment, 
--     order_date datetime, 
--     delivery_date datetime, 
--     product_id int not null, 
--     user_id int not null, 
--     amount int default 0, 
--     shipping_address varchar(1000), 
--     createdAt datetime default current_timestamp, 
--     updatedAt datetime default current_timestamp on update current_timestamp,
--     primary key (id), 
--     foreign key (product_id) references products(id), 
--     foreign key (user_id) references users(id) 
-- );


