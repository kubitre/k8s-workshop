# Step 1: Docker

## Checking docker installation
1. Run `docker run hello-world`
    - Connect can fail in case of "sudo" permissions;
    - Expected to receive: 
    ```
    Hello from Docker!
    This message shows that your installation appears to be working correctly.
    ...
    ```

## Case A: docker for cli utilities
1. Lets imagine that we have an SQL instance: 
    - Register new account at [remotemysql.com](remotemysql.com);
    - Create free database, i.e:
    ```
    You have successfully created a new database. The details are below.
    Username: c3STeeV6TKK
    Database name: c3STeeV6TKK
    Password: QlHdcrxVlY
    Server: remotemysql.com
    Port: 3306
    ```

2. Now we need to manage our database. We may want to i.e. install adminer. But we need to install php, download some code, run a web server.. Other way we can use docker to avoid extra installations:
    - Run `docker run --name adminer --rm -p 8080:8080 -d adminer:4.7.2-standalone`
    - Navigate to `127.0.0.1`, enter credentials from step 1;
    - Congrats, now you can manage your remote sql without extra installations;
    - Run `docker stop adminer` to stop & cleanup used adminer.

## Case B: Our docker network
1. Lets imagine that we want to build a stateful application, some php code and some database. To simplify support we may want to use docker:
    - Lets prepare a project directory, like `mkdir -p ~/projects/k8s-workshop/step-1 && cd ~/projects/k8s-workshop/step-1`;
    - Then lets prepare some project structure:
        - `mkdir -r shared/db` - create a directory for mounted data to keep some state from docker;
    - Firstly, we need a database; 
        - Run `docker run --name db --rm -e MYSQL_DATABASE=app -e MYSQL_ROOT_PASSWORD=root -e MYSQL_USER=user -e MYSQL_PASSWORD=pass -v $PWD/shared/db/mysql:/var/lib/mysql -u $(id -u):$(id -g) -d mariadb:10.1`;
    - Then, we need an application, adminer is still ok:
        - Run `docker run --name adminer --rm -p 8080:8080 -d adminer:4.7.2-standalone`;
    - Then we need to make a "bridge" between them. We will use docker network:
        - Run `docker network create my-app` to create a network;
        - Run `docker network connect my-app db` to connect our "db" container to the network;
        - Run `docker network connect my-app adminer` to connect our "adminer" container to the network;
    - Now you can use your app (adminer) at `127.0.0.1:8080` to connect the db: `user:pass@db:3306/app`;
    - Run `docker stop adminer db && docker network remove my-app` to stop your application;
        - Anyway you still have your db data at `./shared/db`;
    - You can run you application again and make sure that your is in safe. 

2. But its too verbose way. To ease your work you can use `docker-compose`;
    - Lets declare our services from point 1 using docker compose:
    ```yaml
    version: "3"
    services:
      db:
        image: mariadb:10.1
        user: ${UID}:${GID}
        environment:
          - MYSQL_DATABASE=mydb
          - MYSQL_ROOT_PASSWORD=root
          - MYSQL_USER=user
          - MYSQL_PASSWORD=pass
        volumes:
          - ./shared/db/mysql:/var/lib/mysql
      adminer:
        user: ${UID}:${GID}
        image: adminer
        ports:
          - 8080:8080
    ```
    - To automate `${UID}:${GID}` bindings we need to create `.env`:
        - Run `echo "$(id -u):$(id -g)" > .env`;
    - Now we can raise the same applications set: `docker-compose up -d`;
        - (you can check that its working, `127.0.0.1:8080`, etc.);
    - You can call `docker-compose down` to shutdown everything;
    - You can use i.e. `restart: always` option to make sure that your service will be restarted after crash/system reboot/etc.
 