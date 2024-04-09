system clear; 
use first; 


insert into categories
(name, description)
values 
(
    "Smartphones",
    "Mobile communication devices with advanced computing capabilities."
); 


















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


