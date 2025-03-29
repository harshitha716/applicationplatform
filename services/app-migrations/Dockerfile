FROM migrate/migrate:v4.17.1

WORKDIR /migrations

COPY . .

RUN chmod a+x ./execute.sh

ENTRYPOINT ["./execute.sh"]