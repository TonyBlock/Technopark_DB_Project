FROM golang:latest AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/api/main.go

FROM ubuntu:latest
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER Forum_user WITH SUPERUSER PASSWORD 'db_password';" &&\
    createdb -O Forum_user db_forum &&\
    psql -f ./db/db.sql -d db_forum &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main .

EXPOSE 5000
USER root
CMD service postgresql start && ./main