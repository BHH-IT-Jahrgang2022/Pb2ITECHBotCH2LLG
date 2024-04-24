FROM ubuntu
RUN apt-get update
RUN apt-get install -y apache2
RUN apt-get install -y apache2-utils
RUN apt-get clean

COPY ./web /var/www/html

EXPOSE 80

CMD ["apache2ctl", "-D", "FOREGROUND"]
