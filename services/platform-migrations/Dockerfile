FROM migrate/migrate:v4.17.0

WORKDIR /migrations

COPY . .

RUN chmod a+x ./execute.sh

ENTRYPOINT ["./execute.sh"]