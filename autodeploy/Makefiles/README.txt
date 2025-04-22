Запуск сценариев make-файла осуществляется по одному. Обращаться к сценарию необходимо по алиасу. Н-р:

make -f makefile.txt postgres

make -f autodeploy/Makefiles/makefile.txt migrateup

Путь к make-файлу указывается относительно нахождения в текущем каталоге.

