\c postgres
CREATE EXTENSION IF NOT EXISTS dblink;
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'users') THEN
      PERFORM dblink_exec('dbname=postgres user=' || current_user, 'CREATE DATABASE bossdb');
   END IF;
END
$$;
\c bossdb
DO
$$
    BEGIN
        -- Создаем таблицу ролей пользователей
        CREATE TABLE IF NOT EXISTS user_roles
        (
            id   serial PRIMARY KEY,
            name varchar(50) NOT NULL UNIQUE -- Имя роли (например, администратор, пользователь), уникальное поле
        );

-- Создаем таблицу категорий продуктов
        CREATE TABLE IF NOT EXISTS categories_product
        (
            id   serial PRIMARY KEY,
            name varchar(50) NOT NULL UNIQUE -- Название категории продукта
        );

-- Создаем таблицу типов продвижения
        CREATE TABLE IF NOT EXISTS types_promotion
        (
            id        serial PRIMARY KEY,
            name      varchar(50)    NOT NULL UNIQUE,             -- Название типа продвижения
            price     numeric(10, 2) NOT NULL CHECK (price >= 0), -- Цена продвижения
            time_live interval       NOT NULL                     -- Время действия продвижения
        );

-- Создаем таблицу пользователей
        CREATE TABLE IF NOT EXISTS users
        (
            id                  int PRIMARY KEY,
            path_ava            text,
            username            varchar(50) NOT NULL UNIQUE,                       -- Имя пользователя, уникальное поле
            firstname           varchar(50) NOT NULL,                              -- Имя
            lastname            varchar(50),                                       -- Фамилия
            number_phone        varchar(12),                                       -- Номер телефона
            rating              numeric(3, 2)        DEFAULT 0.00,                 -- Рейтинг пользователя
            verification_status varchar(20) NOT NULL DEFAULT 'unverified',         -- Статус верификации
            role_id             int         NOT NULL,                              -- Идентификатор роли пользователя (внешний ключ)
            CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES user_roles (id) -- Связь с таблицей user_roles
        );

-- Создаем таблицу объявлений
        CREATE TABLE IF NOT EXISTS advertisements
        (
            id                    serial PRIMARY KEY,
            user_id               int            NOT NULL,                                          -- Внешний ключ на пользователя
            name                  varchar(50)    NOT NULL,                                          -- Название объявления
            description           varchar(255),                                                     -- Описание объявления
            price                 numeric(10, 2) NOT NULL CHECK (price >= 0),                       -- Цена, не может быть отрицательной
            date_placement        timestamp DEFAULT CURRENT_TIMESTAMP,                              -- Дата размещения
            location              varchar(50),                                                      -- Местоположение
            type_id               int,                                                              -- Внешний ключ на тип продвижения
            views_count           int       DEFAULT 0,                                              -- Количество просмотров
            date_expire_promotion timestamp,                                                        -- Дата окончания продвижения
            category_id           int            NOT NULL,                                          -- Внешний ключ на категорию продукта
            CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES categories_product (id), -- Связь с категорией продукта
            CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),                      -- Связь с пользователем
            CONSTRAINT fk_type_id FOREIGN KEY (type_id) REFERENCES types_promotion (id)             -- Связь с типом продвижения
        );

-- Создаем таблицу отзывов
        CREATE TABLE IF NOT EXISTS reviews
        (
            id               serial PRIMARY KEY,
            text             text,                                                                       -- Текст отзыва
            mark             int CHECK (mark BETWEEN 1 AND 5),                                           -- Оценка от 1 до 5
            reviewer_id      int NOT NULL,                                                               -- Внешний ключ на пользователя, который оставил отзыв
            advertisement_id int NOT NULL,                                                               -- Внешний ключ на объявление, к которому относится отзыв
            CONSTRAINT fk_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES users (id),                   -- Связь с пользователем
            CONSTRAINT fk_advertisement_id FOREIGN KEY (advertisement_id) REFERENCES advertisements (id) -- Связь с объявлением
        );

-- Создаем таблицу фотографий объявлений
        CREATE TABLE IF NOT EXISTS ad_photos
        (
            id               serial PRIMARY KEY,
            path             text NOT NULL,                                                              -- Путь к изображению
            advertisement_id int  NOT NULL,                                                              -- Внешний ключ на объявление
            CONSTRAINT fk_advertisement_id FOREIGN KEY (advertisement_id) REFERENCES advertisements (id) -- Связь с объявлением
        );

-- Создаем таблицу сделок
        CREATE TABLE IF NOT EXISTS deals
        (
            id               serial PRIMARY KEY,
            advertisement_id int NOT NULL,                                                                                  -- Внешний ключ на объявление
            buyer_id         int NOT NULL,                                                                                  -- Внешний ключ на покупателя
            date_deal        timestamp DEFAULT CURRENT_TIMESTAMP,                                                           -- Дата и время сделки
            CONSTRAINT fk_advertisement_id FOREIGN KEY (advertisement_id) REFERENCES advertisements (id) ON DELETE CASCADE, -- Связь с таблицей объявлений
            CONSTRAINT fk_buyer_id FOREIGN KEY (buyer_id) REFERENCES users (id) ON DELETE CASCADE                           -- Связь с таблицей пользователей (покупателей)
        );
        RAISE NOTICE 'Таблицы успешно созданы.';
    EXCEPTION
        WHEN OTHERS THEN
            RAISE EXCEPTION 'Ошибка при создании таблиц: %', SQLERRM;
    END
$$;
