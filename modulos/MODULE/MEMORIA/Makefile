obj-m +=  ram_grupo5.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean

run:
	sudo insmod ram_grupo5.ko

remove: 
	sudo rmmod ram_grupo5.ko