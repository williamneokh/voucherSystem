FROM arm64v8/mysql:oracle

COPY voucher_system.sql /tmp


CMD [ "mysqld", "--init-file=/tmp/voucher_system.sql" ]