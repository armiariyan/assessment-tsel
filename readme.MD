# Assessment Tsel

## Personal Notes
- For ease of testing by users, I’m using a trial `PostgreSQL` database via [Railway](__https://railway.app/__), so users don’t need to install a database and can directly hit the endpoints.
- If you prefer to use your own local database, make sure you have PostgreSQL installed on your computer, then run the `DDL` at the bottom of this `README` to create the tables and some sample data.

### How to run this project?

We can run this Assessment project with or without Docker. Here, I am providing both ways to run this project.



#### Run with Docker

- Install docker in your machine
- Run this command
```bash
# pull image from the docker hub
docker pull armiariyan/assessment-tsel

# after complete pull, run this
docker run -p 9999:9999 armiariyan/assessment-tsel
```

#### Run without Docker

- Clone this project

```bash
# Move to your workspace
cd your-workspace

# Clone this project into your workspace
git clone https://github.com/armiariyan/assessment-tsel.git

# Move to the project root directory
cd assessment-tsel
```
- Create a file `.env` similar to `.env.example` at the root directory with your configuration.
- Install `go` if not installed on your machine.
- If you plan to use your local database, make sure you install `PostgreSQL` if not installed on your machine
- If you want to use local database (i recommend just use the online one, so skip this step) Change your `.env` section `env` into `"local"` so when app running, its used local database
- Important: Change your `.env` section `postgresql.products.*` into your local database properties
- Run `go run main.go`.
- Access API using `http://localhost:9999`

### API documentation of this project

#### Online docs
<a href="https://documenter.getpostman.com/view/26445251/2sA3s6GViK" target="_blank">
    <img alt="View API Doc Button" src="https://github.com/amitshekhariitbhu/go-backend-clean-architecture/blob/main/assets/button-view-api-docs.png?raw=true" width="200" height="60"/>
</a>

#### Offline docs
If you prefer direct import to postman, I'm attached my these APIs collection inside this folder, please check it out the file named `assessment_tsel.postman_collection`

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
Once again, make sure you have a PostgreSQL database, then run this query before starting to test the endpoints.  

You don’t need to run this query if you decide to use the online database.

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

### Example API Request and Response
This cover all positive case, for negative case refer to postman online documentation above at section API documentation of this project
- Healthcheck

  - Request

  ```
  curl --location 'http://localhost:9999'
  ```

  - Response

  ```json
  {
    "message": "Assessment Tsel - Server up and running",
    "server_time": "Thu, 15 Aug 2024 01:01:37 WIB",
    "version": "v1.0.0"
  }
  ```

- Create Product

  - Request

  ```
    curl --location 'http://localhost:9999/v1/products' \
    --header 'Content-Type: application/json' \
    --data '{
        "name": "Test Product",
        "description": "Sarung untuk mandi pagi deskripsi",
        "price": 99009,
        "variety": {
            "colors": ["red", "blue", "green"],
            "sizes": ["S", "M", "L", "XL"]
        },
        "stock": 150
    }'
  ```

  - Response

  ```json
    {
        "status": "200",
        "message": "success",
        "data": null
    }
  ```

- Update Product

  - Request

  ```
    curl --location --request PATCH 'http://localhost:9999/v1/products' \
    --header 'Content-Type: application/json' \
    --data '{
        "id": 1,
        "name": "UPDATE Kaos slimitfit wanita",
        "description": "Kaos slimfit wanita",
        "price": 100000,
        "stock": 150,
        "rating": 4,
        "variety": {
            "sizes": [
                "S",
                "M",
                "L",
                "XL"
            ],
            "colors": [
                "black",
                "white"
            ]
        },
        "createdAt": "2024-08-14T10:48:49.945826Z",
        "updatedAt": "2024-08-14T10:48:49.945826Z"
    }'
  ```

  - Response

  ```json
    {
        "status": "200",
        "message": "success update data",
        "data": {
            "id": 2,
            "name": "UPDATE Kaos slimitfit wanita",
            "description": "Kaos slimfit wanita",
            "price": 100000,
            "stock": 150,
            "rating": 4.5,
            "variety": {
                "sizes": [
                    "S",
                    "M",
                    "L",
                    "XL"
                ],
                "colors": [
                    "black",
                    "white"
                ]
            },
            "createdAt": "2024-08-14T10:48:49.945826Z",
            "updatedAt": "2024-08-14T22:46:54.426361Z"
        }
    }
  ```

- Delete Product

  - Request

  ```
    curl --location --request DELETE 'http://localhost:9999/v1/products/6'
  ```

  - Response

  ```json
    {
        "status": "200",
        "message": "success delete data",
        "data": null
    }
  ```

- Get List Products

  - Request

  ```
    curl --location 'http://localhost:9999/v1/products?page=1&limit=10'
  ```

  - Response

  ```json
    {
        "status": "200",
        "message": "success",
        "data": {
            "results": [
                {
                    "id": 3,
                    "name": "Updated Jaket Denim",
                    "description": "Jaket denim updated nih",
                    "price": 125000,
                    "stock": 150,
                    "rating": 4.5,
                    "variety": {
                        "sizes": [
                            "S",
                            "M",
                            "L",
                            "XL"
                        ],
                        "colors": [
                            "red",
                            "blue",
                            "green"
                        ]
                    },
                    "createdAt": "2024-08-14T11:20:01.624094Z",
                    "updatedAt": "2024-08-14T11:29:06.354793Z"
                },
                {
                    "id": 4,
                    "name": "Sarung jamur",
                    "description": "Sarung jamur deskripsi",
                    "price": 99009,
                    "stock": 150,
                    "rating": null,
                    "variety": {
                        "sizes": [
                            "S",
                            "M",
                            "L",
                            "XL"
                        ],
                        "colors": [
                            "red",
                            "blue",
                            "green"
                        ]
                    },
                    "createdAt": "2024-08-14T22:25:56.743047Z",
                    "updatedAt": "2024-08-14T22:25:56.743047Z"
                },
                {
                    "id": 5,
                    "name": "Test kurang kurang",
                    "description": "Sarung jamur deskripsi",
                    "price": 99009,
                    "stock": 150,
                    "rating": 4.4,
                    "variety": {
                        "sizes": [
                            "S",
                            "M",
                            "L",
                            "XL"
                        ],
                        "colors": [
                            "red",
                            "blue",
                            "green"
                        ]
                    },
                    "createdAt": "2024-08-14T22:28:01.449332Z",
                    "updatedAt": "2024-08-14T22:29:06.21072Z"
                },
                {
                    "id": 2,
                    "name": "UPDATE Kaos slimitfit wanita",
                    "description": "Kaos slimfit wanita",
                    "price": 100000,
                    "stock": 150,
                    "rating": 4.5,
                    "variety": {
                        "sizes": [
                            "S",
                            "M",
                            "L",
                            "XL"
                        ],
                        "colors": [
                            "black",
                            "white"
                        ]
                    },
                    "createdAt": "2024-08-14T10:48:49.945826Z",
                    "updatedAt": "2024-08-14T22:46:54.426361Z"
                }
            ],
            "pagination": {
                "page": 1,
                "totalPages": 1,
                "totalItems": 4,
                "limit": 10,
                "hasNext": false,
                "hasPrevious": false
            }
        }
    }
  ```

- Get Detail Product by ID

  - Request

  ```
    curl --location 'http://localhost:9999/v1/products/6'
  ```

  - Response

  ```json
    {
        "status": "200",
        "message": "success",
        "data": {
            "id": 2,
            "name": "Kaos Polos Wanita",
            "description": "Kaos slimfit wanita",
            "price": 100000,
            "stock": 150,
            "rating": null,
            "variety": {
                "sizes": [
                    "S",
                    "M",
                    "L",
                    "XL"
                ],
                "colors": [
                    "black",
                    "white"
                ]
            },
            "createdAt": "2024-08-14T10:48:49.945826Z",
            "updatedAt": "2024-08-14T10:48:49.945826Z"
        }
    }
  ```