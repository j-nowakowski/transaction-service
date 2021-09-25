# @source: https://medium.com/@pybrarian/mysql-databases-that-dont-retain-data-293bc2ed7f02
FROM mysql:5
ENV MYSQL_ROOT_PASSWORD=root
RUN mkdir /var/lib/mysql-no-volume
# RUN rm /tmp/output/errors.txt
# RUN rm /tmp/output/email.csv
CMD ["--datadir", "/var/lib/mysql-no-volume"]
