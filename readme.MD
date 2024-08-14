# Assessment Tsel

## Personal Notes
- Untuk kemudahan pengetesan oleh user, saya menggunakan free `PostgreSQL` database via [Railway](__https://railway.app/__), jadi user tidak perlu install database dan bisa langsung hit endpoint
- Jika anda menginginkan menggunakan database local anda sendiri, pastikan anda sudah menginstall postgreSQL pada komputer anda, lalu silahkan menjalankan DDL di bagian paling bawah readme ini untuk membuat table dan beberapa data contoh

### How to run this project?

We can run this Assessment project with or without Docker. Here, I am providing both ways to run this project.

- Clone this project

```bash
# Move to your workspace
cd your-workspace

# Clone this project into your workspace
git clone https://github.com/armiariyan/assessment-tsel.git

# Move to the project root directory
cd assessment-tsel
```

#### Run with Docker

- Create a file `.env` similar to `.env.example` at the root directory with your configuration.
- Install Docker and Docker Compose.
- Run `docker-compose up -d`.
- Access API using `http://localhost:9999`

#### Run without Docker

- Create a file `.env` similar to `.env.example` at the root directory with your configuration.
- Install `go` if not installed on your machine.
- If you plan to use your local database, make sure you install `PostgreSQL` if not installed on your machine
- Change your `.env` section `env` into `"local"` so when app running, its used local database
- Important: Change your `.env` section `postgresql.products.*` into your local database properties
- Run `go run main.go`.
- Access API using `http://localhost:9999`

### Postman Collection
I'm attached my these APIs collection inside this folder, please check it out the file named `assessment_tsel.postman_collection`

### How to run the test?

```bash
# Run all tests for usecase
go test ./...
```

### How to generate the mock code?

In this project, to test, we need to generate mock code for repository, and who knows later we need too for usecase or database

```bash
# Generate mock code for the usecase and repository
make mocks
```

Whenever you make changes in the interface repository, you need to run the corresponding command to regenerate the mock code for testing.

### Diagram Data Language (DDL) for Local Database
Sekali lagi pastikan anda sudah punya database postgresql, lalu lakukan query ini sebelum memulai pengetesan endpoint.

Anda tidak perlu melakukan query ini jika memutuskan menggunakan database onlines

```sql
create table public.products
(
    id          serial
        primary key,
    name        varchar(255)                not null,
    description text,
    price       numeric(10, 2)              not null,
    variety     jsonb,
    rating      numeric(2, 1) default 0
        constraint products_rating_check
            check ((rating >= 0.0) AND (rating <= 5.0)),
    stock       integer                     not null,
    created_at  timestamp     default now() not null,
    updated_at  timestamp     default now() not null,
    deleted_at  timestamp
);

INSERT INTO public.products (id, name, description, price, variety, rating, stock, created_at, updated_at, deleted_at) VALUES (1, 'Kaos Polos Pria', 'Kaos polos pria cuttingan oversize', 125000.00, '{"sizes": ["S", "M", "L", "XXL"], "colors": ["red", "blue", "green"]}', 3.9, 150, '2024-08-14 10:29:14.856867', '2024-08-14 22:12:45.023702', '2024-08-14 22:52:38.219952');
INSERT INTO public.products (id, name, description, price, variety, rating, stock, created_at, updated_at, deleted_at) VALUES (2, 'UPDATE Kaos slimitfit wanita', 'Kaos slimfit wanita', 100000.00, '{"sizes": ["S", "M", "L", "XL"], "colors": ["black", "white"]}', 4.5, 150, '2024-08-14 10:48:49.945826', '2024-08-14 22:46:54.426361', null);
INSERT INTO public.products (id, name, description, price, variety, rating, stock, created_at, updated_at, deleted_at) VALUES (3, 'Updated Jaket Denim', 'Jaket denim updated nih', 125000.00, '{"sizes": ["S", "M", "L", "XL"], "colors": ["red", "blue", "green"]}', 4.5, 150, '2024-08-14 11:20:01.624094', '2024-08-14 11:29:06.354793', null);
INSERT INTO public.products (id, name, description, price, variety, rating, stock, created_at, updated_at, deleted_at) VALUES (4, 'Sarung jamur', 'Sarung jamur deskripsi', 99009.00, '{"sizes": ["S", "M", "L", "XL"], "colors": ["red", "blue", "green"]}', null, 150, '2024-08-14 22:25:56.743047', '2024-08-14 22:25:56.743047', null);
INSERT INTO public.products (id, name, description, price, variety, rating, stock, created_at, updated_at, deleted_at) VALUES (5, 'Test kurang kurang', 'Sarung jamur deskripsi', 99009.00, '{"sizes": ["S", "M", "L", "XL"], "colors": ["red", "blue", "green"]}', 4.4, 150, '2024-08-14 22:28:01.449332', '2024-08-14 22:29:06.210720', null);
INSERT INTO public.products (id, name, description, price, variety, rating, stock, created_at, updated_at, deleted_at) VALUES (6, 'Sarung', 'Sarung untuk mandi pagi deskripsi', 99009.00, '{"sizes": ["S", "M", "L", "XL"], "colors": ["red", "blue", "green"]}', null, 150, '2024-08-14 22:30:55.358444', '2024-08-14 22:30:55.358444', '2024-08-15 00:46:31.913650');


```
